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

func TestRoundF(t *testing.T) {
	f0 := 11.12345678901234567890

	floater := NewFloater(0)
	f := floater.RoundF(f0)
	t.Logf("%v:%s", floater.Accuracy(), floater.Ftoa(f))

	floater = NewFloater(10)
	f = floater.RoundF(f0)
	t.Logf("%v:%s", floater.Accuracy(), floater.Ftoa(f))

	floater = NewFloater(14)
	f = floater.RoundF(f0)
	t.Logf("%v:%s", floater.Accuracy(), floater.Ftoa(f))
}

func TestRoundA(t *testing.T) {
	s0 := "11.12345678901234567890"

	floater := NewFloater(0)
	s, err := floater.RoundA(s0, 64)
	t.Logf("%v:%s %v", floater.Accuracy(), s, err)

	floater = NewFloater(10)
	s, err = floater.RoundA(s0, 64)
	t.Logf("%v:%s %v", floater.Accuracy(), s, err)

	floater = NewFloater(14)
	s, err = floater.RoundA(s0, 64)
	t.Logf("%v:%s %v", floater.Accuracy(), s, err)
}
