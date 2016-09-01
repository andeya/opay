package opay

import (
	"time"
)

type Context struct {
	account     Accounter //instance of Account Interface
	withAccount Accounter //the second party's instance of Account Interface
	Request
}

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
	return ctx.Request.IOrder.ToPend(ctx.Request.Tx, ctx.Request.response)
}

// 标记订单为正在处理状态，或有相关异步回调操作
func (ctx *Context) ToDo() error {
	return ctx.Request.IOrder.ToDo(ctx.Request.Tx, ctx.Request.response)
}

// 处理账户并标记订单为成功状态
func (ctx *Context) ToSucceed() error {
	return ctx.Request.IOrder.ToSucceed(ctx.Request.Tx, ctx.Request.response)
}

// 标记订单为撤销状态
func (ctx *Context) ToCancel() error {
	return ctx.Request.IOrder.ToCancel(ctx.Request.Tx, ctx.Request.response)
}

// 标记订单为失败状态
func (ctx *Context) ToFail() error {
	return ctx.Request.IOrder.ToFail(ctx.Request.Tx, ctx.Request.response)
}

// 同步处理订单，并标记为成功状态
func (ctx *Context) SyncDeal() error {
	return ctx.Request.IOrder.SyncDeal(ctx.Request.Tx, ctx.Request.response)
}

// 针对 Uid-Aid 账户，修改账户余额。
func (ctx *Context) UpdateBalance() error {
	return ctx.account.UpdateBalance(ctx.Request.IOrder.GetUid(), ctx.Request.IOrder.GetAmount(), ctx.Request.Tx, ctx.Request.response)
}

// 针对 Uid-Aid 账户，回滚账户余额。
func (ctx *Context) RollbackBalance() error {
	return ctx.account.UpdateBalance(ctx.Request.IOrder.GetUid(), -ctx.Request.IOrder.GetAmount(), ctx.Request.Tx, ctx.Request.response)
}

// 针对 Uid2-Aid 账户，修改账户余额。
func (ctx *Context) UpdateUid2Balance() error {
	return ctx.account.UpdateBalance(ctx.Request.IOrder.GetUid2(), -ctx.Request.IOrder.GetAmount(), ctx.Request.Tx, ctx.Request.response)
}

// 针对 Uid2-Aid 账户，回滚账户余额。
func (ctx *Context) RollbackUid2Balance() error {
	return ctx.account.UpdateBalance(ctx.Request.IOrder.GetUid2(), ctx.Request.IOrder.GetAmount(), ctx.Request.Tx, ctx.Request.response)
}

// 针对 Uid-Aid2 账户，修改账户余额。
func (ctx *Context) UpdateAid2Balance() error {
	return ctx.withAccount.UpdateBalance(ctx.Request.IOrder.GetUid(), ctx.Request.IOrder.GetAmount2(), ctx.Request.Tx, ctx.Request.response)
}

// 针对 Uid-Aid2 账户，回滚账户余额。
func (ctx *Context) RollbackAid2Balance() error {
	return ctx.withAccount.UpdateBalance(ctx.Request.IOrder.GetUid(), -ctx.Request.IOrder.GetAmount2(), ctx.Request.Tx, ctx.Request.response)
}

// Write response values.
func (ctx *Context) Write(k string, v interface{}) {
	ctx.Request.write(k, v)
}
