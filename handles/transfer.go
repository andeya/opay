package handles

import (
	"github.com/henrylee2cn/opay"
)

/*
 * 转账
 */
type Transfer struct {
	Background
}

// 编译期检查接口实现
var _ Handler = (*Transfer)(nil)

// 执行入口
func (t *Transfer) ServeOpay(ctx *opay.Context) error {
	if !ctx.HasStakeholder() {
		return opay.ErrStakeholderNotExist
	}
	if opay.GreaterOrEqual(ctx.Request.Initiator.GetAmount(), 0) ||
		opay.SmallerOrEqual(ctx.Request.Stakeholder.GetAmount(), 0) {
		return opay.ErrIncorrectAmount
	}
	return t.Call(t, ctx)
}

// 处理账户并标记订单为成功状态，
// IOrder.ToSucceed()中应包含Uid2的订单创建与标记成功
func (t *Transfer) ToSucceed() error {
	// 操作账户
	err := t.Background.Context.UpdateBalance()
	if err != nil {
		return err
	}

	// 更新订单
	return t.Background.Context.ToSucceed()
}

// 实时转账
func (t *Transfer) SyncDeal() error {
	// 操作账户
	err := t.Background.Context.UpdateBalance()
	if err != nil {
		return err
	}

	// 更新订单
	return t.Background.Context.ToSucceed()
}
