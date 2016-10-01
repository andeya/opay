package opay

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

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

type Floater struct {
	numOfDecimalPlaces uint8
	accuracy           float64
	format             string
}

func NewFloater(numOfDecimalPlaces uint8) *Floater {
	var accuracy float64 = 1
	if numOfDecimalPlaces > 0 {
		accuracyString := "0." + strings.Repeat("0", int(numOfDecimalPlaces-1)) + "1"
		accuracy, _ = strconv.ParseFloat(accuracyString, 64)
	}
	return &Floater{
		numOfDecimalPlaces: numOfDecimalPlaces,
		accuracy:           accuracy,
		format:             "%0." + strconv.Itoa(int(numOfDecimalPlaces)) + "f",
	}
}

func (this *Floater) NumOfDecimalPlaces() uint8 {
	return this.numOfDecimalPlaces
}

func (this *Floater) Accuracy() float64 {
	return this.accuracy
}

func (this *Floater) Format() string {
	return this.format
}

func (this *Floater) Ftoa(f float64) string {
	return fmt.Sprintf(this.format, f)
}

func (this *Floater) Atof(s string, bitSize int) (float64, error) {
	f, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		return f, err
	}
	return strconv.ParseFloat(fmt.Sprintf(this.format, f), bitSize)
}

func (this *Floater) Equal(a, b float64) bool {
	return math.Abs(a-b) < this.accuracy
}

func (this *Floater) Greater(a, b float64) bool {
	return math.Max(a, b) == a && math.Abs(a-b) > this.accuracy
}

func (this *Floater) Smaller(a, b float64) bool {
	return math.Max(a, b) == b && math.Abs(a-b) > this.accuracy
}

func (this *Floater) GreaterOrEqual(a, b float64) bool {
	return math.Max(a, b) == a || math.Abs(a-b) < this.accuracy
}

func (this *Floater) SmallerOrEqual(a, b float64) bool {
	return math.Max(a, b) == b || math.Abs(a-b) < this.accuracy
}
