package opay

import (
	"log"
	"runtime"
	"sync"
	"time"
)

type (
	// 订单队列
	Queue interface {
		GetCap() int
		SetCap(int)
		Push(Request) (respChan <-chan *Response)
		Pull() Request
		GetOpay() *Opay
	}

	OrderChan struct {
		c    chan Request
		mu   sync.RWMutex
		opay *Opay
	}
)

const (
	DEFAULT_QUEUE_CAP = 1024 // DEFAULT_QUEUE_CAP is the queue default capacity
)

func newOrderChan(queueCapacity int, opay *Opay) Queue {
	if queueCapacity <= 0 {
		queueCapacity = DEFAULT_QUEUE_CAP
	}
	return &OrderChan{
		c:    make(chan Request, queueCapacity),
		opay: opay,
	}
}

// GetCap returns queue capacity.
func (oc *OrderChan) GetCap() int {
	oc.mu.RLock()
	defer oc.mu.RUnlock()
	return cap(oc.c)
}

// SetCap sets the queue capacity.
func (oc *OrderChan) SetCap(queueCapacity int) {
	if queueCapacity <= 0 {
		queueCapacity = DEFAULT_QUEUE_CAP
	}
	close(oc.c)
	if len(oc.c) > 0 {
		log.Println("Waiting for the completion of the remaining order processing...")
		for len(oc.c) > 0 {
			runtime.Gosched()
		}
	}
	oc.mu.Lock()
	oc.c = make(chan Request, queueCapacity)
	oc.mu.Unlock()

	log.Println("Successfully set the queue capacity.")
}

// Push an order
func (oc *OrderChan) Push(req Request) (respChan <-chan *Response) {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	respChan, err := req.prepare(oc.GetOpay())
	if err != nil {
		req.setError(err)
		req.writeback()
		return
	}

	timeout, err := checkTimeout(req.Deadline)

	if err != nil {
		// Time out, cancel processing
		req.setError(err)
		req.writeback()
		return
	}

	if timeout > 0 {
		// Not timed out
		select {
		case oc.c <- req:
		case <-time.After(timeout):
			err = ErrTimeout
			req.setError(err)
			req.writeback()
		}

	} else {
		// No timeout limit
		oc.c <- req
	}

	return
}

// Read an order.
// Wait indefinitely until a valid order is taken.
// Automatically processes overtime orders.
func (oc *OrderChan) Pull() Request {
	var (
		req Request
		c   chan Request
	)

	for {
		oc.mu.RLock()
		c = oc.c
		oc.mu.RUnlock()

		req = <-c
		if req.isNil() {
			continue
		}

		// If timeout, cancel the order.
		if _, err := checkTimeout(req.Deadline); err != nil {
			req.setError(err)
			req.writeback()
			continue
		}
		break
	}

	return req
}

// GetOpay returns Opay
func (oc *OrderChan) GetOpay() *Opay {
	oc.mu.RLock()
	defer oc.mu.RUnlock()
	return oc.opay
}
