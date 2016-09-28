package opay

import (
	"sync"
)

type (
	CtxStore interface {
		Param(k string) (v interface{}, ok bool)
		Set(k string, v interface{})
		Get(k string) (v interface{}, ok bool)
	}

	// The result of dealing request.
	Response struct {
		addition map[string]interface{}
		Result   map[string]interface{}
		Err      error
		lock     sync.RWMutex
	}
)

var _ CtxStore = new(Response)

func (resp *Response) Param(k string) (interface{}, bool) {
	resp.lock.RLock()
	defer resp.lock.RUnlock()
	v, ok := resp.addition[k]
	return v, ok
}

func (resp *Response) Set(k string, v interface{}) {
	resp.lock.Lock()
	resp.Result[k] = v
	resp.lock.Unlock()
}

func (resp *Response) Get(k string) (interface{}, bool) {
	resp.lock.RLock()
	defer resp.lock.RUnlock()
	v, ok := resp.Result[k]
	return v, ok
}

// Set response error
func (resp *Response) SetError(err error) {
	resp.lock.Lock()
	resp.Err = err
	resp.lock.Unlock()
}

func (resp *Response) GetErr() error {
	resp.lock.RLock()
	defer resp.lock.RUnlock()
	return resp.Err
}
