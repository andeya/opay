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
		SetCap(int)
		Push(Request) (done <-chan struct{}, err error)
		pull() Request
	}

	OrderChan struct {
		c  chan Request
		mu sync.RWMutex
	}
)

const (
	DEFAULT_QUEUE_LEN = 1024
)

func newOrderChan(queueCapacity int) Queue {
	if queueCapacity <= 0 {
		queueCapacity = DEFAULT_QUEUE_LEN
	}
	return &OrderChan{
		c: make(chan Request, queueCapacity),
	}
}

// 设置队列容量
func (oc *OrderChan) SetCap(queueCapacity int) {
	if queueCapacity <= 0 {
		queueCapacity = DEFAULT_QUEUE_LEN
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

// 推送一条订单
func (oc *OrderChan) Push(req Request) (done <-chan struct{}, err error) {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	c := make(chan struct{})
	req.done = (chan<- struct{})(c)
	done = (<-chan struct{})(c)

	timeout, err := checkTimeout(req.Deadline)

	if err != nil {
		// 已超时，取消处理
		req.writeback(err)
		return
	}

	if timeout > 0 {
		// 未超时
		select {
		case oc.c <- req:
		case <-time.After(timeout):
			err = ErrTimeout
			req.writeback(err)
		}

	} else {
		// 无超时限制
		oc.c <- req
	}

	return
}

// 读出一条订单
// 无限等待，直到取出一个有效订单
// 超时订单，自动处理
func (oc *OrderChan) pull() Request {
	var req Request
	for {
		oc.mu.RLock()
		req = <-oc.c
		if req != emptyRequest {
			oc.mu.RUnlock()

			// 超时取消对订单的处理
			if _, err := checkTimeout(req.Deadline); err != nil {
				req.writeback(err)
				continue
			}

			// 取得有效订单，返回
			break
		}

		oc.mu.RUnlock()
		runtime.Gosched()
	}

	return req
}
