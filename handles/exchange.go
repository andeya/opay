package handles

import (
	"github.com/henrylee2cn/opay"
)

/*
 * 兑换
 */
type Exchange struct {
	Background
}

// 编译期检查接口实现
var _ Handler = (*Exchange)(nil)

// 执行入口
func (e *Exchange) ServeOpay(ctx *opay.Context) error {
	return e.Call(e, ctx)
}

// 处理账户并标记订单为成功状态
func (e *Exchange) ToSucceed() error {
	// 操作账户
	err := e.Background.Context.UpdateBalance()
	if err != nil {
		return err
	}

	err = e.Background.Context.UpdateAid2Balance()
	if err != nil {
		return err
	}

	// 更新订单
	return e.Background.Context.ToSucceed()
}

// 实时兑换
func (e *Exchange) SyncDeal() error {
	// 操作账户
	err := e.Background.Context.UpdateBalance()
	if err != nil {
		return err
	}

	err = e.Background.Context.UpdateAid2Balance()
	if err != nil {
		return err
	}

	return e.Background.Context.SyncDeal()
}
