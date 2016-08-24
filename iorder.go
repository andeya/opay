package opay

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	// 订单接口
	IOrder interface {
		//上下文执行控制
		context.Context

		//获取订单类型
		GetType() string

		//获取用户ID
		GetUid() string

		//获取相对应的用户ID
		GetWithUid() string

		//获取资产ID
		GetAid() string

		//获取相对应的资产ID（如用于资产间兑换业务）
		GetWithAid() string

		//更新订单状态
		// @tx 当在一个事务中时，作为数据库的操作句柄
		Update(status int, notes string, tx *sqlx.Tx) error

		//完成处理，写回结果
		Writeback(err error)
	}

	// 订单处理接口
	Handler interface {
		ServeOpay(Context) error
	}

	// 上下文
	Context struct {
		IOrder                //订单接口
		Account     Accounter //账户操作接口实例
		WithAccount Accounter //相对应的账户操作接口实例
		*sqlx.DB              //数据库操作对象
	}
)

// 订单处理接口函数
type HandlerFunc func(Context) error

func (hf HandlerFunc) ServeOpay(ctx Context) error {
	return hf(ctx)
}

// 订单操作接口路由
type ServeMux struct {
	mu sync.RWMutex
	m  map[string]Handler
}

var (
	ErrNotFoundHandler = errors.New("Not Found Handler")
)

// 通过路由执行订单处理
func (mux *ServeMux) Exec(
	iOrd IOrder,
	accounter Accounter,
	withAccounter Accounter,
	db *sqlx.DB,
) {
	mux.mu.RLock()
	h, ok := mux.m[iOrd.GetType()]
	mux.mu.RUnlock()

	if !ok {
		iOrd.Writeback(ErrNotFoundHandler)
		return
	}

	iOrd.Writeback(h.ServeOpay(Context{
		IOrder:      iOrd,
		Account:     accounter,
		WithAccount: withAccounter,
		DB:          db,
	}))
}

// 订单操作接口的全局路由
var globalServeMux = &ServeMux{
	m: make(map[string]Handler),
}

// 注册订单处理接口
func Handle(typ string, handler Handler) error {
	globalServeMux.mu.Lock()
	defer globalServeMux.mu.Unlock()
	_, ok := globalServeMux.m[typ]
	if ok {
		return errors.New("Handler \"" + typ + "\" has been registered.")
	}
	globalServeMux.m[typ] = handler
	return nil
}

// 注册订单处理接口
func HandleFunc(typ string, handler func(Context) error) {
	Handle(typ, HandlerFunc(handler))
}

// 处理超时的订单
func dealTimeout(iOrd IOrder) (time.Duration, error) {
	deadline, ok := iOrd.Deadline()

	// 无超时限制
	if !ok {
		return 0, nil
	}

	timeout := deadline.Sub(time.Now())

	// 已超时，取消订单处理
	if timeout <= 0 {
		iOrd.Writeback(ErrTimeout)
		return timeout, ErrTimeout
	}

	// 未超时
	return timeout, nil
}
