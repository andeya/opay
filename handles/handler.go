package handles

import (
	"errors"

	"github.com/henrylee2cn/opay"
)

/*
 * 操作的基础结构体
 */
type (
	// 订单处理接口
	Handler interface {
		// 订单处理核心接口
		opay.Handler

		// 新建订单，并标记为等待处理状态
		ToPend() error

		// 标记订单为正在处理状态，或有相关异步回调操作
		ToDo() error

		// 处理账户并标记订单为成功状态
		ToSucceed() error

		// 标记订单为撤销状态
		ToCancel() error

		// 标记订单为失败状态
		ToFail() error

		// 同步处理订单，并标记为成功状态
		SyncDeal() error
	}

	// 实现基本操作接口
	Background struct {
		*opay.Context
	}
)

var ErrAction = errors.New("Action not supported.")

// 执行入口
func (b *Background) Call(handler Handler, ctx *opay.Context) error {
	b.Context = ctx
	switch b.Action() {
	case opay.FAIL:
		return handler.ToFail()
	case opay.CANCEL:
		return handler.ToCancel()
	case opay.PEND:
		return handler.ToPend()
	case opay.DO:
		return handler.ToDo()
	case opay.SUCCEED:
		return handler.ToSucceed()
	case opay.SYNC_DEAL:
		return handler.SyncDeal()
	}
	return ErrAction
}

// 处理账户并标记订单为成功状态
func (b *Background) SyncDeal() error {
	return ErrAction
}
