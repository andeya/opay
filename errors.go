package opay

import (
	"errors"
)

var (
	ErrTimeout = errors.New("opay: add to queue timeout.")

	ErrInvalidStatus       = errors.New("opay: order status is invalid.")
	ErrStakeholderNotExist = errors.New("opay: stakeholder Order is not exist.")
	ErrExtraStakeholder    = errors.New("opay: stakeholder Order is extra.")
	ErrIncorrectAmount     = errors.New("opay: account operation amount is incorrect.")
	ErrInitiatorNil        = errors.New("opay: request.Initiator Can not be nil.")

	ErrIllegalStep       = errors.New("opay: illegal Step.")
	ErrInvalidOperation  = errors.New("opay: invalid operation.")
	ErrCancelStep        = errors.New("opay: the Order cannot be canceled.")
	ErrReprocess         = errors.New("opay: repeat process order.")
	ErrDifferentOperator = errors.New("opay: initiator's Operator and Stakeholder must be same.")
	ErrDifferentStep     = errors.New("opay: initiator's Step and Stakeholder must be same.")
)
