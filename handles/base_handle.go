package handles

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/henrylee2cn/opay"
)

/**
 * 操作的基础结构体
 */

type BaseHandle struct {
	*opay.Context
}

var ErrAction = errors.New("Action not supported.")

func (bh *BaseHandle) SetContext(ctx *opay.Context) {
	bh.Context = ctx
}

func (bh *BaseHandle) Call(handler opay.Handler) error {
	var name string
	switch bh.Action() {
	case opay.FAIL:
		name = "Fail"
	case opay.CANCEL:
		name = "Cancel"
	case opay.PEND:
		name = "Pend"
	case opay.DO:
		name = "Do"
	case opay.SUCCEED:
		name = "Succeed"
	default:
		return ErrAction
	}

	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	v := reflect.ValueOf(handler).MethodByName(name)
	if v == (reflect.Value{}) {
		return ErrAction
	}

	err = v.Interface().(func(*opay.Context) error)(bh.Context)

	return err
}
