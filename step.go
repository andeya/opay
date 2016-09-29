package opay

type (
	// handling order's action
	Step int
)

// 六种订单处理行为状态
const (
	FAIL      Step = UNSET - 2 //处理失败
	CANCEL    Step = UNSET - 1 //取消订单
	UNSET     Step = 0         //未设置
	PEND      Step = UNSET + 1 //等待处理
	DO        Step = UNSET + 2 //正在处理
	SUCCEED   Step = UNSET + 3 //处理成功
	SYNC_DEAL Step = UNSET + 4 //同步处理至成功
)

var (
	steps = map[Step]bool{
		FAIL:      true,
		CANCEL:    true,
		UNSET:     true,
		PEND:      true,
		DO:        true,
		SUCCEED:   true,
		SYNC_DEAL: true,
	}
)
