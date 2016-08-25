package handles

import (
	"github.com/henrylee2cn/opay"
)

/*
 * 提现
 */
type Withdraw struct {
	Background
}

// 编译期检查接口实现
var _ Handler = (*Withdraw)(nil)

// 执行入口
func (w *Withdraw) ServeOpay(ctx *opay.Context) error {
	return w.Call(w, ctx)
}

// 处理账户并标记订单为成功状态
func (w *Withdraw) ToSucceed() error {
	// 操作账户
	err := w.UpdateBalance()
	if err != nil {
		return err
	}

	// 更新订单
	return w.ToSucceed()
}
