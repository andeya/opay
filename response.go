package opay

import (
	"log"
	"sync"
)

// The result of dealing respuest.
type Response struct {
	Err      error
	respChan chan<- *Response //result signal
	done     bool
	lock     sync.RWMutex
}

// Set response error
func (resp *Response) setError(err error) {
	resp.lock.Lock()
	resp.Err = err
	resp.lock.Unlock()
}

// Complete the dealing of the respuest.
func (resp *Response) writeback() {
	resp.lock.Lock()
	defer resp.lock.Unlock()
	if resp.done {
		log.Println("repeated writeback.")
		return
	}
	resp.respChan <- resp
	resp.done = true
	close(resp.respChan)
}
