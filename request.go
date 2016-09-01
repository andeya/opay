package opay

import (
	"log"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

type Request struct {
	Key      string                 //the specified handler
	Action   Action                 //the specified handler's action
	Deadline time.Time              //handle timeouts, if do not fill, no limit
	IOrder                          //instance of Order Interface
	Values   map[string]interface{} //addition params
	response *Response
	respChan chan<- Response //result signal
	*sqlx.Tx                 //the optional, database transaction
	done     bool
	lock     sync.RWMutex
}

// 检查处理行为Action是否合法
func (req *Request) ValidateAction() error {
	// 检查是否超出Action范围
	if !actions[req.Action] {
		return ErrInvalidAction
	}

	// 检查是否为重复处理
	if req.IOrder.LastAction() == req.Action {
		return ErrReprocess
	}
	return nil
}

// Prepare the request.
func (req *Request) prepare() (respChan <-chan Response) {
	req.done = false
	if req.Values == nil {
		req.Values = make(map[string]interface{})
	}
	req.response = &Response{Values: req.Values}
	c := make(chan Response)
	req.respChan = (chan<- Response)(c)
	return (<-chan Response)(c)
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
