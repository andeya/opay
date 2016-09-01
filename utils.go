package opay

import (
	"errors"
	"time"
)

var ErrTimeout = errors.New("Add to queue timeout.")

func checkTimeout(deadline time.Time) (timeout time.Duration, errTimeout error) {
	// No timeout
	if deadline.IsZero() {
		return
	}

	timeout = deadline.Sub(time.Now())

	// Timeout, cancel order.
	if timeout <= 0 {
		return timeout, ErrTimeout
	}

	// No timeout
	return
}
