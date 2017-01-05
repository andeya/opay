package opay

import (
	"github.com/jmoiron/sqlx"
)

type (
	// Operation interface of order.
	IOrder interface {
		// for the handler of dealing.
		GetMeta() *Meta

		// Get the previous status.
		PreStatus() int64

		// Get the target status.
		TargetStatus() int64

		// Get user's id.
		GetUid() string

		// Get asset id.
		GetAid() string

		// Get the amount of change for the Uid-Aid account,
		// balance of positive and negative representation.
		GetAmount() float64

		// Async execution, and mark pending.
		Pend(*sqlx.Tx, KV) error

		// Async execution, and mark the doing.
		Do(*sqlx.Tx, KV) error

		// Async execution, and mark the successful.
		Succeed(*sqlx.Tx, KV) error

		// Async execution, and mark canceled.
		Cancel(*sqlx.Tx, KV) error

		// Async execution, and mark failure.
		Fail(*sqlx.Tx, KV) error

		// Sync execution, and mark the successful.
		SyncDeal(*sqlx.Tx, KV) error
	}
)
