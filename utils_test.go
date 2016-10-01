package opay

import (
	"testing"
)

func TestFtoa(t *testing.T) {
	f := 11.12345678901234567890

	floater := NewFloater(0)
	t.Log(floater.Accuracy(), floater.Ftoa(f))

	floater = NewFloater(10)
	t.Log(floater.Accuracy(), floater.Ftoa(f))

	floater = NewFloater(15)
	t.Log(floater.Accuracy(), floater.Ftoa(f))

	floater = NewFloater(20)
	t.Log(floater.Accuracy(), floater.Ftoa(f))

	floater = NewFloater(25)
	t.Log(floater.Accuracy(), floater.Ftoa(f))
}

func TestAtof(t *testing.T) {
	s := "11.12345678901234567890"

	floater := NewFloater(0)
	f, err := floater.Atof(s, 64)
	t.Logf("%v:%0.25f %v", 0, f, err)

	floater = NewFloater(10)
	f, err = floater.Atof(s, 64)
	t.Logf("%v:%0.25f %v", floater.Accuracy(), f, err)

	floater = NewFloater(15)
	f, err = floater.Atof(s, 64)
	t.Logf("%v:%0.25f %v", floater.Accuracy(), f, err)

	floater = NewFloater(20)
	f, err = floater.Atof(s, 64)
	t.Logf("%v:%0.25f %v", floater.Accuracy(), f, err)

	floater = NewFloater(25)
	f, err = floater.Atof(s, 64)
	t.Logf("%v:%0.25f %v", floater.Accuracy(), f, err)
}
