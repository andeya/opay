package opay

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	// 上下文
	Context struct {
		account     Accounter //账户操作接口实例
		withAccount Accounter //相对应的账户操作接口实例
		Request               //请求
	}

	// 请求
	Request struct {
		Key      string          //指定处理类型
		Action   Action          //指定订单处理行为
		Deadline time.Time       //处理超时，不填则不限时
		IOrder                   //订单接口实例
		Addition interface{}     //附加参数
		*sqlx.Tx                 //可选，数据库事务操作
		done     chan<- struct{} //处理结束的信号
	}

	// 订单接口
	IOrder interface {
		// 获取订单当前状态，默认为UNSET==0
		CurAction() Action

		// 获取用户ID
		GetUid() string

		// 获取相对应的用户ID
		GetWithUid() string

		// 获取资产ID
		GetAid() string

		// 获取相对应的资产ID（如用于资产间兑换业务）
		GetWithAid() string

		// 获取针对 Uid-Aid 账户的变化量，正负表示收支。
		GetAmount() float64

		// 获取针对 Uid-WithAid 账户的变化量，正负表示收支。
		GetWithAidAmount() float64

		// 异步处理类方法

		// 新建订单，并标记为等待处理状态
		ToPend(tx *sqlx.Tx, addition interface{}) error

		// 标记订单为正在处理状态，或有相关异步回调操作
		ToDo(tx *sqlx.Tx, addition interface{}) error

		// 处理账户并标记订单为成功状态
		ToSucceed(tx *sqlx.Tx, addition interface{}) error

		// 标记订单为撤销状态
		ToCancel(tx *sqlx.Tx, addition interface{}) error

		// 标记订单为失败状态
		ToFail(tx *sqlx.Tx, addition interface{}) error

		// 同步处理方法

		// 同步处理至成功
		SyncDeal(tx *sqlx.Tx, addition interface{}) error

		// 处理结果

		// 读取处理错误
		Err() error

		// 回写错误
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

// 新建订单，并标记为等待处理状态
func (ctx *Context) ToPend() error {
	return ctx.Request.IOrder.ToPend(ctx.Request.Tx, ctx.Request.Addition)
}

// 标记订单为正在处理状态，或有相关异步回调操作
func (ctx *Context) ToDo() error {
	return ctx.Request.IOrder.ToDo(ctx.Request.Tx, ctx.Request.Addition)
}

// 处理账户并标记订单为成功状态
func (ctx *Context) ToSucceed() error {
	return ctx.Request.IOrder.ToSucceed(ctx.Request.Tx, ctx.Request.Addition)
}

// 标记订单为撤销状态
func (ctx *Context) ToCancel() error {
	return ctx.Request.IOrder.ToCancel(ctx.Request.Tx, ctx.Request.Addition)
}

// 标记订单为失败状态
func (ctx *Context) ToFail() error {
	return ctx.Request.IOrder.ToFail(ctx.Request.Tx, ctx.Request.Addition)
}

// 同步处理订单，并标记为成功状态
func (ctx *Context) SyncDeal() error {
	return ctx.Request.IOrder.SyncDeal(ctx.Request.Tx, ctx.Request.Addition)
}

// 针对 Uid-Aid 账户，修改账户余额。
func (ctx *Context) UpdateBalance() error {
	return ctx.account.UpdateBalance(ctx.Request.IOrder.GetUid(), ctx.Request.IOrder.GetAmount(), ctx.Request.Tx, ctx.Request.Addition)
}

// 针对 Uid-Aid 账户，回滚账户余额。
func (ctx *Context) RollbackBalance() error {
	return ctx.account.UpdateBalance(ctx.Request.IOrder.GetUid(), -ctx.Request.IOrder.GetAmount(), ctx.Request.Tx, ctx.Request.Addition)
}

// 针对 WithUid-Aid 账户，修改账户余额。
func (ctx *Context) UpdateWithUidBalance() error {
	return ctx.account.UpdateBalance(ctx.Request.IOrder.GetWithUid(), -ctx.Request.IOrder.GetAmount(), ctx.Request.Tx, ctx.Request.Addition)
}

// 针对 WithUid-Aid 账户，回滚账户余额。
func (ctx *Context) RollbackWithUidBalance() error {
	return ctx.account.UpdateBalance(ctx.Request.IOrder.GetWithUid(), ctx.Request.IOrder.GetAmount(), ctx.Request.Tx, ctx.Request.Addition)
}

// 针对 Uid-WithAid 账户，修改账户余额。
func (ctx *Context) UpdateWithAidBalance() error {
	return ctx.withAccount.UpdateBalance(ctx.Request.IOrder.GetUid(), ctx.Request.IOrder.GetWithAidAmount(), ctx.Request.Tx, ctx.Request.Addition)
}

// 针对 Uid-WithAid 账户，回滚账户余额。
func (ctx *Context) RollbackWithAidBalance() error {
	return ctx.withAccount.UpdateBalance(ctx.Request.IOrder.GetUid(), -ctx.Request.IOrder.GetWithAidAmount(), ctx.Request.Tx, ctx.Request.Addition)
}

// 六种订单处理行为状态
const (
	FAIL      Action = UNSET - 2 //处理失败
	CANCEL    Action = UNSET - 1 //取消订单
	UNSET     Action = 0         //未设置
	PEND      Action = UNSET + 1 //等待处理
	DO        Action = UNSET + 2 //正在处理
	SUCCEED   Action = UNSET + 3 //处理成功
	SYNC_DEAL Action = UNSET + 4 //同步处理至成功
)

var (
	actions = map[Action]bool{
		FAIL:      true,
		CANCEL:    true,
		UNSET:     true,
		PEND:      true,
		DO:        true,
		SUCCEED:   true,
		SYNC_DEAL: true,
	}

	ErrInvalidAction = errors.New("Invalid Action.")
	ErrReprocess     = errors.New("Repeat process order.")
)

// 检查处理行为Action是否合法
func (req *Request) ValidateAction() error {
	// 检查是否超出Action范围
	if !actions[req.Action] {
		return ErrInvalidAction
	}

	// 检查是否为重复处理
	if req.IOrder.CurAction() == req.Action {
		return ErrReprocess
	}
	return nil
}

// 标记请求处理结束，并回写错误
func (req *Request) writeback(err error) {
	close(req.done)
	req.IOrder.SetErr(err)
}

// 空请求
var emptyRequest Request
