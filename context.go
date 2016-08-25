package opay

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	// 上下文
	Context struct {
		Account     Accounter //账户操作接口实例
		WithAccount Accounter //相对应的账户操作接口实例
		Request               //请求
	}

	// 请求
	Request struct {
		Key      string          //指定处理类型
		Action   Action          //指定订单处理行为
		Deadline time.Time       //处理超时，不填则不限时
		IOrder                   //订单接口实例
		*sqlx.Tx                 //可选，数据库事务操作
		done     chan<- struct{} //处理结束的信号
	}

	// 订单接口
	IOrder interface {
		// 获取用户ID
		GetUid() string

		// 获取相对应的用户ID
		GetWithUid() string

		// 获取资产ID
		GetAid() string

		// 获取相对应的资产ID（如用于资产间兑换业务）
		GetWithAid() string

		// 新建订单，并标记为等待处理状态
		ToPend(tx *sqlx.Tx) error

		// 标记订单为正在处理状态，或有相关异步回调操作
		ToDo(tx *sqlx.Tx) error

		// 处理账户并标记订单为成功状态
		ToSucceed(tx *sqlx.Tx) error

		// 标记订单为撤销状态
		ToCancel(tx *sqlx.Tx) error

		// 标记订单为失败状态
		ToFail(tx *sqlx.Tx) error

		//回写错误
		SetErr(err error)
	}

	// 订单处理行为
	Action int
)

// 获取指定处理类型
func (ctx *Context) Key() string {
	return ctx.Request.Key
}

// 获取指定订单处理行为
func (ctx *Context) Action() Action {
	return ctx.Request.Action
}

// 获取处理超时，不填则不限时
func (ctx *Context) Deadline() time.Time {
	return ctx.Request.Deadline
}

// 五种订单处理行为
const (
	FAIL    Action = PEND - 2 //处理失败
	CANCEL  Action = PEND - 1 //取消订单
	PEND    Action = 0        //等待处理
	DO      Action = PEND + 1 //正在处理
	SUCCEED Action = PEND + 2 //处理成功
)

// 空请求
var emptyRequest Request

// 标记请求处理结束，并回写错误
func (req *Request) writeback(err error) {
	close(req.done)
	req.IOrder.SetErr(err)
}
