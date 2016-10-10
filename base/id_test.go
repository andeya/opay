package base

import (
	"testing"
)

func TestCreateOrderid(t *testing.T) {
	SetTimeZone("CST", 8)
	orderid := createOrderid("a")
	t.Log(orderid)
	t.Log(GetAidFromOrderid(orderid))
}
