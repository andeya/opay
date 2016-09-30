package opay

import (
	"errors"
)

var (
	ErrTimeout = errors.New("opay: add to queue timeout.")

	ErrInvalidStatus       = errors.New("opay: order status is invalid.")
	ErrStakeholderNotExist = errors.New("opay: stakeholder order is not exist.")
	ErrExtraStakeholder    = errors.New("opay: stakeholder order is extra.")
	ErrIncorrectAmount     = errors.New("opay: account operation amount is incorrect.")
	ErrInitiatorNil        = errors.New("opay: request.Initiator can not be nil.")

	ErrIllegalStep       = errors.New("opay: illegal step.")
	ErrInvalidOperation  = errors.New("opay: invalid operation.")
	ErrCancelStep        = errors.New("opay: the order cannot be canceled.")
	ErrReprocess         = errors.New("opay: repeat process order.")
	ErrDifferentOperator = errors.New("opay: initiator's operator and stakeholder's must be same.")
	ErrDifferentStep     = errors.New("opay: initiator's step and stakeholder's must be same.")
)
