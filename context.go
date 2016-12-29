package opay

import (
	"time"
)

// Context is used to process order information.
type Context struct {
	initiatorSettle   SettleFunc
	stakeholderSettle SettleFunc
	Request
	*Response
	*Floater
}

// Deadline gets processing deadline, not limited if not fill.
func (ctx *Context) Deadline() time.Time {
	return ctx.Request.Deadline
}

// Pend creates an order, and marks it as pending.
func (ctx *Context) Pend() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.Pend(ctx.Request.Tx)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.Pend(ctx.Request.Tx)
}

// Do marks the order as being in progress, and maybe have an associated asynchronous callback operation.
func (ctx *Context) Do() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.Do(ctx.Request.Tx)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.Do(ctx.Request.Tx)
}

// Succeed processes the account and marks the order as successful.
func (ctx *Context) Succeed() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.Succeed(ctx.Request.Tx)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.Succeed(ctx.Request.Tx)
}

// Cancel marks the order as Canceled.
func (ctx *Context) Cancel() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.Cancel(ctx.Request.Tx)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.Cancel(ctx.Request.Tx)
}

// Fail marks the order as failed.
func (ctx *Context) Fail() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.Fail(ctx.Request.Tx)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.Fail(ctx.Request.Tx)
}

// SyncDeal The order is processed synchronously and marked as a successful status.
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

// Modify the account balance.
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

// Roll back the account balance.
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

// Param returns response value.
func (ctx *Context) Param(k string) interface{} {
	return ctx.Request.param(k)
}

// Write response values.
func (ctx *Context) Write(k string, v interface{}) {
	ctx.Response.write(k, v)
}
