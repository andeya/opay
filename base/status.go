package base

import (
	"sync"
)

// 订单类型的状态文本描述
var ost = &struct {
	text map[uint8]map[int32]string //Type:Status:Text
	lock sync.RWMutex
}{text: map[uint8]map[int32]string{}}

// 获取订单类型的状态文本描述
func GetStatusText(typ uint8, status int32) string {
	ost.lock.RLock()
	defer ost.lock.RUnlock()
	return ost.text[typ][status]
}

// 设置订单类型的状态文本描述
func SetStatusText(typ uint8, status int32, text string) {
	ost.lock.Lock()
	defer ost.lock.Unlock()
	if _, ok := ost.text[typ]; !ok {
		ost.text[typ] = map[int32]string{}
	}
	ost.text[typ][status] = text
}
