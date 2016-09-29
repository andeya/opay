package opay

import (
	"time"
)

type Context struct {
	initiatorSettle   SettleFunc
	stakeholderSettle SettleFunc
	Request
	*Response
	Accuracy
}

// 获取处理超时，不填则不限时
func (ctx *Context) Deadline() time.Time {
	return ctx.Request.Deadline
}

// 新建订单，并标记为等待处理状态
func (ctx *Context) Pend() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.Pend(ctx.Request.Tx)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.Pend(ctx.Request.Tx)
}

// 标记订单为正在处理状态，或有相关异步回调操作
func (ctx *Context) Do() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.Do(ctx.Request.Tx)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.Do(ctx.Request.Tx)
}

// 处理账户并标记订单为成功状态
func (ctx *Context) Succeed() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.Succeed(ctx.Request.Tx)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.Succeed(ctx.Request.Tx)
}

// 标记订单为撤销状态
func (ctx *Context) Cancel() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.Cancel(ctx.Request.Tx)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.Cancel(ctx.Request.Tx)
}

// 标记订单为失败状态
func (ctx *Context) Fail() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.Fail(ctx.Request.Tx)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.Fail(ctx.Request.Tx)
}

// 同步处理订单，并标记为成功状态
func (ctx *Context) SyncDeal() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.SyncDeal(ctx.Request.Tx)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.SyncDeal(ctx.Request.Tx)
}

func (ctx *Context) HasStakeholder() bool {
	return ctx.Request.Stakeholder != nil
}

// 修改账户余额。
func (ctx *Context) UpdateBalance() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.stakeholderSettle(
			ctx.Request.Stakeholder.GetUid(),
			ctx.Request.Stakeholder.GetAmount(),
			ctx.Request.Tx,
		)
		if err != nil {
			return err
		}
	}
	return ctx.initiatorSettle(
		ctx.Request.Initiator.GetUid(),
		ctx.Request.Initiator.GetAmount(),
		ctx.Request.Tx,
	)
}

// 回滚账户余额。
func (ctx *Context) RollbackBalance() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.stakeholderSettle(
			ctx.Request.Stakeholder.GetUid(),
			-ctx.Request.Stakeholder.GetAmount(),
			ctx.Request.Tx,
		)
		if err != nil {
			return err
		}
	}

	return ctx.initiatorSettle(
		ctx.Request.Initiator.GetUid(),
		-ctx.Request.Initiator.GetAmount(),
		ctx.Request.Tx,
	)
}

func (ctx *Context) Param(k string) interface{} {
	return ctx.Request.param(k)
}

// Write response values.
func (ctx *Context) Write(k string, v interface{}) {
	ctx.Response.write(k, v)
}
