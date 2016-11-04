package opay

import (
	"errors"
)

var (
	// ErrTimeout = errors.New("opay: add to queue timeout.")
	ErrTimeout = errors.New("加入交易队列超时")

	// ErrInvalidStatus       = errors.New("opay: order status is invalid.")
	ErrInvalidStatus = errors.New("无效的交易订单状态")
	// ErrStakeholderNotExist = errors.New("opay: stakeholder order is not exist.")
	ErrStakeholderNotExist = errors.New("关联订单不存在")
	// ErrExtraStakeholder    = errors.New("opay: stakeholder order is extra.")
	ErrExtraStakeholder = errors.New("多余的关联订单")
	// ErrIncorrectAmount     = errors.New("opay: account operation amount is incorrect.")
	ErrIncorrectAmount = errors.New("交易金额不正确")
	// ErrInitiatorNil        = errors.New("opay: request.Initiator can not be nil.")
	ErrInitiatorNil = errors.New("交易订单为空")

	// ErrIllegalStep       = errors.New("opay: illegal step.")
	ErrIllegalStep = errors.New("非法的交易订单操作")
	// ErrInvalidStep  = errors.New("opay: invalid operation.")
	ErrInvalidStep = errors.New("无效的交易订单操作")
	// ErrCancelStep        = errors.New("opay: the order cannot be canceled.")
	ErrCancelStep = errors.New("交易订单不可撤销")
	// ErrReprocess         = errors.New("opay: repeat process order.")
	ErrReprocess = errors.New("重复操作交易订单")
	// ErrDifferentStep     = errors.New("opay: initiator's step and stakeholder's must be same.")
	ErrDifferentStep = errors.New("关联订单的操作不一致")
	// ErrDifferentOperator = errors.New("opay: initiator's type and stakeholder's must be same.")
	ErrDifferentType = errors.New("关联订单的类型不一致")
)
