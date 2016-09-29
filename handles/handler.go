package handles

import (
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
		Pend() error

		// 标记订单为正在处理状态，或有相关异步回调操作
		Do() error

		// 处理账户并标记订单为成功状态
		Succeed() error

		// 标记订单为撤销状态
		Cancel() error

		// 标记订单为失败状态
		Fail() error

		// 同步处理订单，并标记为成功状态
		SyncDeal() error
	}

	// 实现基本操作接口
	Background struct {
		*opay.Context
	}
)

// 执行入口
func (b *Background) Call(handler Handler, ctx *opay.Context) error {
	b.Context = ctx
	switch b.Step() {
	case opay.FAIL:
		return handler.Fail()
	case opay.CANCEL:
		return handler.Cancel()
	case opay.PEND:
		return handler.Pend()
	case opay.DO:
		return handler.Do()
	case opay.SUCCEED:
		return handler.Succeed()
	case opay.SYNC_DEAL:
		return handler.SyncDeal()
	}
	return opay.ErrIllegalStep
}

// 处理账户并标记订单为成功状态
func (b *Background) SyncDeal() error {
	return opay.ErrIllegalStep
}
