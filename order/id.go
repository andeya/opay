package order

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var timeZone = time.UTC

func SetTimeZone(name string, hourOffset int) {
	timeZone = time.FixedZone("CST", hourOffset*60*60)
}

var orderid = &struct {
	salt int32
	lock sync.Mutex
}{salt: rand.Int31n(1000000)}

// 生成32字节(时间23+类型3+随机6)的订单ID，含第8区时间，纯数字
// 可保证同一进程内全局唯一，重复概率为0
// 不同进程生成的ID几乎不会重复，但仍有重复概率
// 建议：全部产品使用同一个进程生成ID
func CreateOrderid(typ uint8) string {
	orderid.lock.Lock()
	t := time.Now().In(timeZone)
	if orderid.salt >= 1000000 {
		orderid.salt = rand.Int31n(1000000)
	} else {
		orderid.salt++
	}
	salt := orderid.salt
	orderid.lock.Unlock()
	return fmt.Sprintf("%s%09d%03d%06d", t.Format("20060102150405"), t.Nanosecond(), typ, salt)
}
