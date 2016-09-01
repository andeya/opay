package opay

import (
	"sync"
)

type (
	Values interface {
		Set(k string, v interface{})
		Get(k string) (v interface{}, ok bool)
	}

	// The result of dealing request.
	Response struct {
		Values map[string]interface{}
		Err    error
		lock   sync.RWMutex
	}
)

var _ Values = new(Response)

func (resp *Response) Set(k string, v interface{}) {
	resp.lock.Lock()
	resp.Values[k] = v
	resp.lock.Unlock()
}

func (resp *Response) Get(k string) (interface{}, bool) {
	resp.lock.RLock()
	defer resp.lock.RUnlock()
	v, ok := resp.Values[k]
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
