package opay

import (
	"testing"
)

func TestFtoa(t *testing.T) {
	f := 11.1234567890123350

	floater := NewFloater(0)
	t.Log(floater.Accuracy(), floater.Ftoa(f))

	floater = NewFloater(10)
	t.Log(floater.Accuracy(), floater.Ftoa(f))

	floater = NewFloater(14)
	t.Log(floater.Accuracy(), floater.Ftoa(f))
}

func TestAtof(t *testing.T) {
	s := "11.12345678901234567890"

	floater := NewFloater(0)
	f, err := floater.Atof(s, 64)
	t.Logf("%v:%s %v", floater.Accuracy(), floater.Ftoa(f), err)

	floater = NewFloater(10)
	f, err = floater.Atof(s, 64)
	t.Logf("%v:%s %v", floater.Accuracy(), floater.Ftoa(f), err)

	floater = NewFloater(14)
	f, err = floater.Atof(s, 64)
	t.Logf("%v:%s %v", floater.Accuracy(), floater.Ftoa(f), err)
	t.Log(11.12345678901235 == f)
}

func TestFtof(t *testing.T) {
	f0 := 11.12345678901234567890

	floater := NewFloater(0)
	f := floater.Ftof(f0)
	t.Logf("%v:%s", floater.Accuracy(), floater.Ftoa(f))

	floater = NewFloater(10)
	f = floater.Ftof(f0)
	t.Logf("%v:%s", floater.Accuracy(), floater.Ftoa(f))

	floater = NewFloater(14)
	f = floater.Ftof(f0)
	t.Logf("%v:%s", floater.Accuracy(), floater.Ftoa(f))
}

func TestAtoa(t *testing.T) {
	s0 := "11.12345678901234567890"

	floater := NewFloater(0)
	s, err := floater.Atoa(s0, 64)
	t.Logf("%v:%s %v", floater.Accuracy(), s, err)

	floater = NewFloater(10)
	s, err = floater.Atoa(s0, 64)
	t.Logf("%v:%s %v", floater.Accuracy(), s, err)

	floater = NewFloater(14)
	s, err = floater.Atoa(s0, 64)
	t.Logf("%v:%s %v", floater.Accuracy(), s, err)
}

func TestZeroString(t *testing.T) {
	floater := NewFloater(0)
	t.Logf("%v:%v", floater.Accuracy(), floater.IsZero(floater.Accuracy()))

	floater = NewFloater(1)
	t.Logf("%v:%v", floater.Accuracy(), floater.IsZero(0.1))
	t.Logf("%v:%v", floater.Accuracy(), floater.IsZero(0.05))

	floater = NewFloater(14)
	t.Logf("%v:%v", floater.Accuracy(), floater.IsZero(floater.Accuracy()))
}

func TestCompare(t *testing.T) {
	floater := NewFloater(1)
	t.Logf("%v:%v", floater.Accuracy(), floater.Equal(0.1, 0.15))
}
