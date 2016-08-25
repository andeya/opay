package opay

import (
	"errors"
	"time"
)

var ErrTimeout = errors.New("Add to queue timeout.")

// 检查并处理超时订单
func checkTimeout(deadline time.Time) (timeout time.Duration, errTimeout error) {
	// 无超时限制
	if deadline.IsZero() {
		return
	}

	timeout = deadline.Sub(time.Now())

	// 已超时，取消订单处理
	if timeout <= 0 {
		return timeout, ErrTimeout
	}

	// 未超时
	return
}
