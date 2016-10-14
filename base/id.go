package base

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var timeZone = time.UTC

func SetTimeZone(name string, hourOffset int) {
	timeZone = time.FixedZone(name, hourOffset*60*60)
}

var orderid = &struct {
	salt int32
	lock sync.Mutex
}{salt: rand.Int31n(1000000000)}

// 生成32字节(时间23+随机9)的订单ID，含第8区时间，纯数字
// 可保证同一进程内全局唯一，重复概率为0
// 不同进程生成的ID几乎不会重复，但仍有重复概率
// 建议：全部产品使用同一个进程生成ID
func CreateOrderid(aid string) string {
	switch len(aid) {
	case 0:
		aid = "00"
	case 1:
		aid = "0" + aid
	case 2:
	default:
		aid = aid[:2]
	}
	orderid.lock.Lock()
	t := time.Now().In(timeZone)
	if orderid.salt >= 1000000000 {
		orderid.salt = rand.Int31n(1000000000)
	} else {
		orderid.salt++
	}
	salt := orderid.salt
	orderid.lock.Unlock()
	return fmt.Sprintf("%s%s%09d%09d", aid, t.Format("060102150405"), t.Nanosecond(), salt)
}

func GetAidFromOrderid(orderid string) string {
	if len(orderid) < 2 {
		return ""
	}
	if strings.HasPrefix(orderid, "0") {
		return orderid[1:2]
	}
	return orderid[:2]
}

func CheckOrderid(orderid string) (aid string, err error) {
	if len(orderid) != 32 {
		return "", errors.New("orderid is not the correct length.")
	}
	aid = GetAidFromOrderid(orderid)
	if len(aid) == 0 {
		return "", errors.New("orderid's 'aid' section is incorrect.")
	}
	return aid, nil
}
