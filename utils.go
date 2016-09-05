package opay

import (
	"errors"
	"math"
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

const (
	ACCURACY = 0.0000000001
)

func Equal(a, b float64) bool {
	return math.Abs(a-b) < ACCURACY
}

func Greater(a, b float64) bool {
	return math.Max(a, b) == a && math.Abs(a-b) > ACCURACY
}

func Smaller(a, b float64) bool {
	return math.Max(a, b) == b && math.Abs(a-b) > ACCURACY
}

func GreaterOrEqual(a, b float64) bool {
	return math.Max(a, b) == a || math.Abs(a-b) < ACCURACY
}

func SmallerOrEqual(a, b float64) bool {
	return math.Max(a, b) == b || math.Abs(a-b) < ACCURACY
}
