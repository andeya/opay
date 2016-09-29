package handles

import (
	"github.com/henrylee2cn/opay"
)

/*
 * 充值
 */
type Recharge struct {
	Background
}

// 编译期检查接口实现
var _ Handler = (*Recharge)(nil)

// 执行入口
func (r *Recharge) ServeOpay(ctx *opay.Context) error {
	if ctx.HasStakeholder() {
		return opay.ErrExtraStakeholder
	}
	if ctx.SmallerOrEqual(ctx.Request.Initiator.GetAmount(), 0) {
		return opay.ErrIncorrectAmount
	}
	return r.Call(r, ctx)
}

// 处理账户并标记订单为成功状态
func (r *Recharge) Succeed() error {
	// 操作账户
	err := r.Background.Context.UpdateBalance()
	if err != nil {
		return err
	}

	// 更新订单
	return r.Background.Context.Succeed()
}

// 实时充值
func (r *Recharge) SyncDeal() error {
	// 操作账户
	err := r.Background.Context.UpdateBalance()
	if err != nil {
		return err
	}

	// 更新订单
	return r.Background.Context.SyncDeal()
}
