package base

import (
	"log"
	"sync"

	"github.com/henrylee2cn/opay"
)

// 订单类型的状态行为信息
var orderMetaInfos = &struct {
	operator map[uint8]string                //Type:operator
	action   map[uint8]map[int32]opay.Action //Type:Status:Action
	text     map[uint8]map[int32]string      //Type:Status:text
	lock     sync.RWMutex
}{
	operator: map[uint8]string{},
	action:   map[uint8]map[int32]opay.Action{},
	text:     map[uint8]map[int32]string{},
}

// 为订单类型绑定处理操作符
func BindOrderOperator(typ uint8, operator string) {
	orderMetaInfos.lock.Lock()
	defer orderMetaInfos.lock.Unlock()
	if _, ok := orderMetaInfos.operator[typ]; !ok {
		orderMetaInfos.operator[typ] = operator
	} else {
		log.Printf("repeat binding order operator: %d - %s\n", typ, operator)
	}
}

// 绑定订单类型的各状态的行为与描述信息
func BindOrderAboutStatus(typ uint8, status int32, action opay.Action, text string) {
	orderMetaInfos.lock.Lock()
	defer orderMetaInfos.lock.Unlock()
	if _, ok := orderMetaInfos.text[typ]; !ok {
		orderMetaInfos.text[typ] = map[int32]string{}
		orderMetaInfos.action[typ] = map[int32]opay.Action{}
	}
	if _, ok := orderMetaInfos.text[typ][status]; !ok {
		orderMetaInfos.text[typ][status] = text
		orderMetaInfos.action[typ][status] = action
	} else {
		log.Printf("repeat binding order status information: %d - %d - %d - %s\n", typ, status, action, text)
	}
}

// 获取订单类型的绑定的处理操作符
func OrderOperator(typ uint8) string {
	orderMetaInfos.lock.RLock()
	defer orderMetaInfos.lock.RUnlock()
	return orderMetaInfos.operator[typ]
}

// 获取订单类型的状态的行为
func OrderAction(typ uint8, status int32) opay.Action {
	orderMetaInfos.lock.RLock()
	defer orderMetaInfos.lock.RUnlock()
	return orderMetaInfos.action[typ][status]
}

// 获取订单类型的状态文本描述
func OrderStatusText(typ uint8, status int32) string {
	orderMetaInfos.lock.RLock()
	defer orderMetaInfos.lock.RUnlock()
	return orderMetaInfos.text[typ][status]
}
