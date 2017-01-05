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
		err := ctx.Request.Stakeholder.Pend(ctx.Request.Tx, ctx)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.Pend(ctx.Request.Tx, ctx)
}

// Do marks the order as being in progress, and maybe have an associated asynchronous callback operation.
func (ctx *Context) Do() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.Do(ctx.Request.Tx, ctx)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.Do(ctx.Request.Tx, ctx)
}

// Succeed processes the account and marks the order as successful.
func (ctx *Context) Succeed() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.Succeed(ctx.Request.Tx, ctx)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.Succeed(ctx.Request.Tx, ctx)
}

// Cancel marks the order as Canceled.
func (ctx *Context) Cancel() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.Cancel(ctx.Request.Tx, ctx)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.Cancel(ctx.Request.Tx, ctx)
}

// Fail marks the order as failed.
func (ctx *Context) Fail() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.Fail(ctx.Request.Tx, ctx)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.Fail(ctx.Request.Tx, ctx)
}

// SyncDeal The order is processed synchronously and marked as a successful status.
func (ctx *Context) SyncDeal() error {
	if ctx.Request.Stakeholder != nil {
		err := ctx.Request.Stakeholder.SyncDeal(ctx.Request.Tx, ctx)
		if err != nil {
			return err
		}
	}
	return ctx.Request.Initiator.SyncDeal(ctx.Request.Tx, ctx)
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

// KV key-value
type KV interface {
	Get(k string) interface{}
	Set(k string, v interface{})
}

// Get gets a temporary variable.
func (ctx *Context) Get(k string) interface{} {
	return ctx.Request.get(k)
}

// Set sets a temporary variable.
func (ctx *Context) Set(k string, v interface{}) {
	ctx.Request.set(k, v)
}
