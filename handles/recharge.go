package handles

import (
	"github.com/henrylee2cn/opay"
)

/**
 * 充值
 */
type Recharge struct {
	BaseHandle
}

// 执行入口
func (r *Recharge) ServeOpay(ctx *opay.Context) error {
	r.SetContext(ctx)
	return r.Call(r)
}

// 新建订单，并标记为等待处理状态
func (r *Recharge) Pend(ctx *opay.Context) error {
	return nil
}

// 标记订单为正在处理状态，或有相关异步回调操作
func (r *Recharge) Do(ctx *opay.Context) error {
	return nil
}

// 处理账户并标记订单为成功状态
func (r *Recharge) Succeed(ctx *opay.Context) error {
	return nil
}

// 标记订单为撤销状态
func (r *Recharge) Cancel(ctx *opay.Context) error {
	return nil
}

// 标记订单为失败状态
func (r *Recharge) Fail(ctx *opay.Context) error {
	return nil
}

// 提现
func Withdraw(ctx *opay.Context) error {
	return nil
}

// 转账
func Transfer(ctx *opay.Context) error {

	return nil
}

// 兑换
func Exchange(ctx *opay.Context) error {
	return nil
}
