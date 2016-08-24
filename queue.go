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
		Push(IOrder) error
		pull() IOrder
	}

	OrderChan struct {
		c  chan IOrder
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
		c: make(chan IOrder, queueCapacity),
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
	oc.c = make(chan IOrder, queueCapacity)
	oc.mu.Unlock()

	log.Println("Successfully set the queue capacity.")
}

// 推送一条订单
func (oc *OrderChan) Push(iOrd IOrder) error {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	timeout, err := dealTimeout(iOrd)

	// 已超时，取消处理
	if err != nil {
		return err
	}

	// 未超时
	if timeout > 0 {
		select {
		case oc.c <- iOrd:
			return nil
		case <-time.After(timeout):
			iOrd.Writeback(ErrTimeout)
			return ErrTimeout
		}
	}

	// 无超时限制
	oc.c <- iOrd
	return nil
}

// 读出一条订单
// 无限等待，直到取出一个有效订单
// 超时订单，自动处理
func (oc *OrderChan) pull() IOrder {
	var iOrd IOrder
	for {
		oc.mu.RLock()
		iOrd = <-oc.c
		if iOrd != nil {
			oc.mu.RUnlock()

			// 超时取消对订单的处理
			if _, err := dealTimeout(iOrd); err != nil {
				continue
			}

			// 取得有效订单，返回
			break
		}

		oc.mu.RUnlock()
		runtime.Gosched()
	}

	return iOrd
}
