package opay

import (
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	Context struct {
		account     Accounter //instance of Account Interface
		withAccount Accounter //the second party's instance of Account Interface
		Request
	}

	Request struct {
		Key      string      //the specified handler
		Action   Action      //the specified handler's action
		Deadline time.Time   //handle timeouts, if do not fill, no limit
		IOrder               //instance of Order Interface
		Addition interface{} //addition params
		response Response
		respChan chan<- Response //result signal
		*sqlx.Tx                 //the optional, database transaction
		done     bool
	}

	// The result of dealing request.
	Response struct {
		Body map[string]interface{}
		Err  error
	}

	// Operation interface of order.
	IOrder interface {
		// Get the most recent Action, the default value is UNSET==0.
		LastAction() Action

		// Get user's id.
		GetUid() string

		// Get the second party's user id.
		GetUid2() string

		// Get asset id.
		GetAid() string

		// Get the second party's asset id. (for example, the currency exchange business)
		GetAid2() string

		// Get the amount of change for the Uid-Aid account,
		// balance of positive and negative representation.
		GetAmount() float64

		// Get the amount of change for the Uid-Aid2 account,
		// balance of positive and negative representation.
		GetAmount2() float64

		// Async execution, and mark pending.
		ToPend(tx *sqlx.Tx, addition interface{}) error

		// Async execution, and mark the doing.
		ToDo(tx *sqlx.Tx, addition interface{}) error

		// Async execution, and mark the successful.
		ToSucceed(tx *sqlx.Tx, addition interface{}) error

		// Async execution, and mark canceled.
		ToCancel(tx *sqlx.Tx, addition interface{}) error

		// Async execution, and mark failure.
		ToFail(tx *sqlx.Tx, addition interface{}) error

		// Sync execution, and mark the successful.
		SyncDeal(tx *sqlx.Tx, addition interface{}) error
	}

	// handling order's action
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

// 针对 Uid2-Aid 账户，修改账户余额。
func (ctx *Context) UpdateUid2Balance() error {
	return ctx.account.UpdateBalance(ctx.Request.IOrder.GetUid2(), -ctx.Request.IOrder.GetAmount(), ctx.Request.Tx, ctx.Request.Addition)
}

// 针对 Uid2-Aid 账户，回滚账户余额。
func (ctx *Context) RollbackUid2Balance() error {
	return ctx.account.UpdateBalance(ctx.Request.IOrder.GetUid2(), ctx.Request.IOrder.GetAmount(), ctx.Request.Tx, ctx.Request.Addition)
}

// 针对 Uid-Aid2 账户，修改账户余额。
func (ctx *Context) UpdateAid2Balance() error {
	return ctx.withAccount.UpdateBalance(ctx.Request.IOrder.GetUid(), ctx.Request.IOrder.GetAmount2(), ctx.Request.Tx, ctx.Request.Addition)
}

// 针对 Uid-Aid2 账户，回滚账户余额。
func (ctx *Context) RollbackAid2Balance() error {
	return ctx.withAccount.UpdateBalance(ctx.Request.IOrder.GetUid(), -ctx.Request.IOrder.GetAmount2(), ctx.Request.Tx, ctx.Request.Addition)
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
	if req.IOrder.LastAction() == req.Action {
		return ErrReprocess
	}
	return nil
}

// Prepare the request.
func (req *Request) prepare() (respChan <-chan Response) {
	req.done = false
	req.response.Body = make(map[string]interface{})
	c := make(chan Response)
	req.respChan = (chan<- Response)(c)
	return (<-chan Response)(c)
}

// Write response body.
func (req *Request) Write(key string, value interface{}) {
	if req.done {
		log.Println("As it has been submitted, it can not be written.")
		return
	}
	req.response.Body[key] = value
}

// Complete the dealing of the request.
func (req *Request) writeback(err error) {
	if req.done {
		log.Println("repeated writeback.")
		return
	}
	req.response.Err = err
	req.respChan <- req.response
	req.done = true
	close(req.respChan)
}
