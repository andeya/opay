package opay

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type (
	// Operation interface of order.
	IOrder interface {
		// Get the most recent Action, the default value is UNSET==0.
		LastAction() Action

		// Get user's id.
		GetUid() string

		// Get asset id.
		GetAid() string

		// Get the amount of change for the Uid-Aid account,
		// balance of positive and negative representation.
		GetAmount() float64

		// Async execution, and mark pending.
		ToPend(tx *sqlx.Tx, values Values) error

		// Async execution, and mark the doing.
		ToDo(tx *sqlx.Tx, values Values) error

		// Async execution, and mark the successful.
		ToSucceed(tx *sqlx.Tx, values Values) error

		// Async execution, and mark canceled.
		ToCancel(tx *sqlx.Tx, values Values) error

		// Async execution, and mark failure.
		ToFail(tx *sqlx.Tx, values Values) error

		// Sync execution, and mark the successful.
		SyncDeal(tx *sqlx.Tx, values Values) error
	}

	// handling order's action
	Action int
)

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

	ErrInvalidAction   = errors.New("Invalid Action.")
	ErrReprocess       = errors.New("Repeat process order.")
	ErrDifferentAction = errors.New("Initiator's Action and Stakeholder must be same.")
)
