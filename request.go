package opay

import (
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

type Request struct {
	Deadline    time.Time              //handle timeouts, if do not fill, no limit
	Addition    map[string]interface{} //additional params
	Initiator   IOrder                 //master order
	Stakeholder IOrder                 //the optional, slave order
	response    *Response
	*sqlx.Tx    //the optional, database transaction
	operator    string
	step        Step
	lock        sync.RWMutex
}

// 获取指定的订单处理操作符
func (req *Request) Operator() string {
	req.lock.RLock()
	defer req.lock.RUnlock()
	return req.operator
}

// 获取订单处理的行为目标
func (req *Request) Step() Step {
	req.lock.RLock()
	defer req.lock.RUnlock()
	return req.step
}

// Prepare the request.
func (req *Request) prepare(opay *Opay) (respChan <-chan *Response, err error) {
	req.lock.Lock()
	defer req.lock.Unlock()

	c := make(chan *Response, 1)
	respChan = (<-chan *Response)(c)

	req.response = &Response{
		Result:   make(map[string]interface{}),
		respChan: (chan<- *Response)(c),
	}

	// The main order can not be empty.
	if req.Initiator == nil {
		err = ErrInitiatorNil
		return
	}

	meta := req.Initiator.GetMeta()

	// 检查订单状态是否已注册
	preStatus, ok := meta.Status(req.Initiator.PreStatus())
	if !ok {
		err = ErrInvalidStatus
		return
	}
	targetStatus, ok := meta.Status(req.Initiator.TargetStatus())
	if !ok {
		err = ErrInvalidStatus
		return
	}

	// 检查是否为重复处理行为
	if preStatus.Code == targetStatus.Code {
		err = ErrReprocess
		return
	}

	// 设置订单的处理阶段
	req.step = targetStatus.Step

	// 必须设定目标处理行为
	if req.step == UNSET {
		err = ErrInvalidOperation
		return
	}

	curStep := preStatus.Step
	// 不可操作已撤销或已完成的订单
	if curStep == CANCEL ||
		curStep == FAIL ||
		curStep == SUCCEED ||
		curStep == SYNC_DEAL {
		err = ErrInvalidOperation
		return
	}

	// 非待处理状态的订单不可撤销
	if curStep != PEND && req.step == CANCEL {
		err = ErrCancelStep
		return
	}

	// 主订单操作金额不能为0
	if opay.Equal(req.Initiator.GetAmount(), 0) {
		err = ErrIncorrectAmount
		return
	}

	// 检查从属订单
	if req.Stakeholder != nil {
		// 检查主从订单操作是否一致
		if req.Stakeholder.GetMeta() != meta {
			err = ErrDifferentOperator
			return
		}

		// 检查订单状态是否已注册
		preStatus2, ok := meta.Status(req.Stakeholder.PreStatus())
		if !ok {
			err = ErrInvalidStatus
			return
		}
		targetStatus2, ok := meta.Status(req.Stakeholder.TargetStatus())
		if !ok {
			err = ErrInvalidStatus
			return
		}

		// 检查主从订单行为是否一致
		if preStatus2.Step != curStep ||
			targetStatus2.Step != req.step {
			err = ErrDifferentStep
			return
		}

		// 从属订单操作金额不能为0
		if opay.Equal(req.Stakeholder.GetAmount(), 0) {
			err = ErrIncorrectAmount
			return
		}
	}

	if req.Addition == nil {
		req.Addition = make(map[string]interface{})
	}
	req.operator = meta.OrderType()

	return
}

func (req *Request) param(k string) interface{} {
	req.lock.RLock()
	defer req.lock.RUnlock()
	return req.Addition[k]
}

// Write response body.
func (req *Request) write(k string, v interface{}) {
	req.response.write(k, v)
}

// Set response error
func (req *Request) setError(err error) {
	req.response.setError(err)
}

// Complete the dealing of the request.
func (req *Request) writeback() {
	req.response.writeback()
}

func (req *Request) isNil() bool {
	req.lock.RLock()
	defer req.lock.RUnlock()
	return req.response == nil
}
