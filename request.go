package opay

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

type Request struct {
	Key         string                 //the specified handler
	Action      Action                 //the specified handler's action
	Deadline    time.Time              //handle timeouts, if do not fill, no limit
	Initiator   IOrder                 //master order
	Stakeholder IOrder                 //the optional, slave order
	Values      map[string]interface{} //addition params
	response    *Response
	respChan    chan<- Response //result signal
	*sqlx.Tx                    //the optional, database transaction
	done        bool
	lock        sync.RWMutex
}

var (
	ErrStakeholderNotExist = errors.New("Stakeholder Order is not exist.")
	ErrExtraStakeholder    = errors.New("Stakeholder Order is extra.")
	ErrIncorrectAmount     = errors.New("Account operation amount is incorrect.")
)

// 检查处理行为Action是否合法
func (req *Request) ValidateAction() error {
	// 检查是否超出Action范围
	if !actions[req.Action] {
		return ErrInvalidAction
	}

	// 检查是否为重复处理
	if req.Initiator.LastAction() == req.Action {
		return ErrReprocess
	}
	if req.Stakeholder != nil {
		if req.Stakeholder.LastAction() == req.Action {
			return ErrDifferentAction
		}
	}
	return nil
}

// Prepare the request.
func (req *Request) prepare(a Accuracy) (respChan <-chan Response, err error) {
	req.done = false
	if req.Values == nil {
		req.Values = make(map[string]interface{})
	}
	req.response = &Response{Values: req.Values}
	c := make(chan Response)
	req.respChan = (chan<- Response)(c)
	respChan = (<-chan Response)(c)

	if req.Initiator == nil {
		err = errors.New("Request.Initiator Can not be nil.")
		return
	}

	if a.Equal(req.Initiator.GetAmount(), 0) {
		err = ErrIncorrectAmount
		return
	}

	if req.Stakeholder != nil && a.Equal(req.Stakeholder.GetAmount(), 0) {
		err = ErrIncorrectAmount
	}

	return
}

// Write response body.
func (req *Request) write(k string, v interface{}) {
	req.lock.RLock()
	defer req.lock.RUnlock()
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
	return req.response == nil
}
