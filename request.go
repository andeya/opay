package opay

import (
	"errors"
	"log"
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
	respChan    chan<- Response //result signal
	*sqlx.Tx                    //the optional, database transaction
	operator    string
	action      Action
	done        bool
	lock        sync.RWMutex
}

var (
	ErrStakeholderNotExist = errors.New("Stakeholder Order is not exist.")
	ErrExtraStakeholder    = errors.New("Stakeholder Order is extra.")
	ErrIncorrectAmount     = errors.New("Account operation amount is incorrect.")
)

// 获取指定的订单处理操作符
func (req *Request) Operator() string {
	req.lock.RLock()
	defer req.lock.RUnlock()
	return req.operator
}

// 获取订单处理的行为目标
func (req *Request) Action() Action {
	req.lock.RLock()
	defer req.lock.RUnlock()
	return req.action
}

// Prepare the request.
func (req *Request) prepare(engine *Engine) (respChan <-chan Response, err error) {
	req.lock.Lock()
	defer req.lock.Unlock()
	req.done = false
	if req.Addition == nil {
		req.Addition = make(map[string]interface{})
	}
	req.response = &Response{
		addition: req.Addition,
		Result:   make(map[string]interface{}),
	}
	c := make(chan Response)
	req.respChan = (chan<- Response)(c)
	respChan = (<-chan Response)(c)
	err = req.validate(engine)
	if err == nil {
		req.operator = req.Initiator.Operator()
		req.action = req.Initiator.TargetAction()
	}
	return
}

// 检查请求的合法性
func (req *Request) validate(engine *Engine) error {
	// The main order can not be empty.
	if req.Initiator == nil {
		return errors.New("Request.Initiator Can not be nil.")
	}

	// 检查操作是否存在
	if err := engine.CheckOperator(req.Initiator.Operator()); err != nil {
		return err
	}

	// 检查是否为重复处理行为
	if req.Initiator.TargetAction() == req.Initiator.RecentAction() {
		return ErrReprocess
	}

	// 必须设定目标处理行为
	if req.Initiator.TargetAction() == UNSET {
		return ErrUnsetAction
	}

	// 检查处理行为是否超出范围
	if !actions[req.Initiator.TargetAction()] || !actions[req.Initiator.RecentAction()] {
		return ErrInvalidAction
	}

	// 非待处理状态的订单不可撤销
	if req.Initiator.TargetAction() == CANCEL && req.Initiator.RecentAction() != PEND {
		return ErrCancelAction
	}

	// 主订单操作金额不能为0
	if engine.Equal(req.Initiator.GetAmount(), 0) {
		return ErrIncorrectAmount
	}

	// 检查从属订单
	if req.Stakeholder != nil {
		// 检查主从订单操作是否一致
		if req.Stakeholder.Operator() != req.Initiator.Operator() {
			return ErrDifferentOperator
		}

		// 允许从属订单不设定目标行为
		if req.Stakeholder.TargetAction() != UNSET {
			// 检查主从订单行为是否一致
			if req.Stakeholder.TargetAction() != req.Initiator.TargetAction() ||
				req.Stakeholder.RecentAction() != req.Initiator.RecentAction() {
				return ErrDifferentAction
			}
		}

		// 从属订单操作金额不能为0
		if engine.Equal(req.Stakeholder.GetAmount(), 0) {
			return ErrIncorrectAmount
		}
	}

	return nil
}

// Write response body.
func (req *Request) write(k string, v interface{}) {
	req.lock.Lock()
	defer req.lock.Unlock()
	if req.done {
		log.Println("As it has been submitted, it can not be written.")
		return
	}
	req.response.Set(k, v)
}

// Set response error
func (req *Request) setError(err error) {
	req.response.SetError(err)
}

// Complete the dealing of the request.
func (req *Request) writeback() {
	req.lock.Lock()
	defer req.lock.Unlock()
	if req.done {
		log.Println("repeated writeback.")
		return
	}
	req.respChan <- *req.response
	req.done = true
	close(req.respChan)
}

func (req *Request) isNil() bool {
	req.lock.RLock()
	defer req.lock.RUnlock()
	return req.response == nil
}
