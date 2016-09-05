package opay

import (
	"time"
)

type Context struct {
	account     IAccount //instance of Account Interface
	withAccount IAccount //the second party's instance of Account Interface
	Request
	Accuracy
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
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.ToPend(ctx.Request.Tx, ctx.Request.response)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.ToPend(ctx.Request.Tx, ctx.Request.response)
}

// 标记订单为正在处理状态，或有相关异步回调操作
func (ctx *Context) ToDo() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.ToDo(ctx.Request.Tx, ctx.Request.response)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.ToDo(ctx.Request.Tx, ctx.Request.response)
}

// 处理账户并标记订单为成功状态
func (ctx *Context) ToSucceed() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.ToSucceed(ctx.Request.Tx, ctx.Request.response)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.ToSucceed(ctx.Request.Tx, ctx.Request.response)
}

// 标记订单为撤销状态
func (ctx *Context) ToCancel() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.ToCancel(ctx.Request.Tx, ctx.Request.response)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.ToCancel(ctx.Request.Tx, ctx.Request.response)
}

// 标记订单为失败状态
func (ctx *Context) ToFail() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.ToFail(ctx.Request.Tx, ctx.Request.response)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.ToFail(ctx.Request.Tx, ctx.Request.response)
}

// 同步处理订单，并标记为成功状态
func (ctx *Context) SyncDeal() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.SyncDeal(ctx.Request.Tx, ctx.Request.response)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.SyncDeal(ctx.Request.Tx, ctx.Request.response)
}

func (ctx *Context) HasStakeholder() bool {
	return ctx.Request.Stakeholder != nil
}

// 修改账户余额。
func (ctx *Context) UpdateBalance() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.account.UpdateBalance(
			ctx.Request.Stakeholder.GetUid(),
			ctx.Request.Stakeholder.GetAmount(),
			ctx.Request.Tx,
			ctx.Request.response,
		)
		if err != nil {
			return err
		}
	}
	return ctx.account.UpdateBalance(
		ctx.Request.Initiator.GetUid(),
		ctx.Request.Initiator.GetAmount(),
		ctx.Request.Tx,
		ctx.Request.response,
	)
}

// 回滚账户余额。
func (ctx *Context) RollbackBalance() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.account.UpdateBalance(
			ctx.Request.Stakeholder.GetUid(),
			-ctx.Request.Stakeholder.GetAmount(),
			ctx.Request.Tx,
			ctx.Request.response,
		)
		if err != nil {
			return err
		}
	}

	return ctx.account.UpdateBalance(
		ctx.Request.Initiator.GetUid(),
		-ctx.Request.Initiator.GetAmount(),
		ctx.Request.Tx,
		ctx.Request.response,
	)
}

// Write response values.
func (ctx *Context) Write(k string, v interface{}) {
	ctx.Request.write(k, v)
}
