package base

import (
	"testing"
)

func TestCreateOrderid(t *testing.T) {
	SetTimeZone("CST", 8)
	orderid := CreateOrderid("a")
	t.Log(orderid)
	t.Log(GetAidFromOrderid(orderid))
}

func TestGetTimeFromOrderid(t *testing.T) {
	t.Log(GetTimeFromOrderid("1612011008581826898744368413960e"))
}
