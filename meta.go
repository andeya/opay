package opay

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

type (
	Meta struct {
		orderType string
		handler   reflect.Value
		statuses  map[int64]Status
		unsetCode int64
	}
	Status struct {
		Code int64
		Note string
		Step Step
	}
)

func (o *Opay) RegMeta(orderType string, handler Handler, statuses []Status) (*Meta, error) {
	o.metasLock.Lock()
	defer o.metasLock.Unlock()
	_, ok := o.metas[orderType]
	if ok {
		return nil, errors.New("opay: repeat regester order meta: " + orderType)
	}

	v := reflect.ValueOf(handler)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 过滤不允许的类型
	if !(v.Kind() == reflect.Struct || v.Kind() == reflect.Func) {
		return nil, errors.New("opay: handler must be func or struct type.")
	}

	meta := &Meta{
		orderType: orderType,
		handler:   v,
		statuses:  make(map[int64]Status, len(statuses)),
	}
	for _, status := range statuses {
		_, ok := steps[status.Step]
		if !ok {
			return nil, fmt.Errorf("opay: invalid Step: %d", status.Step)
		}
		meta.statuses[status.Code] = status
	}

	for i := int64(math.MinInt64); i <= math.MaxInt64; i++ {
		if _, ok := meta.statuses[i]; !ok {
			meta.statuses[i] = Status{
				Code: i,
				Step: UNSET,
			}
			meta.unsetCode = i
			break
		}
	}
	o.metas[orderType] = meta
	return meta, nil
}

func (o *Opay) Meta(orderType string) (*Meta, bool) {
	o.metasLock.RLock()
	defer o.metasLock.RUnlock()
	meta, ok := o.metas[orderType]
	return meta, ok
}

// func (o *Opay) MetaStatus(orderType string, code int64) (Status, bool) {
// 	o.metasLock.RLock()
// 	defer o.metasLock.RUnlock()
// 	meta, ok := o.metas[orderType]
// 	if !ok {
// 		return Status{}, ok
// 	}
// 	status, ok := meta.statuses[code], ok
// 	return status, ok
// }

// func (o *Opay) MetaUnsetCode(orderType string) (int64, bool) {
// 	o.metasLock.RLock()
// 	defer o.metasLock.RUnlock()
// 	meta, ok := o.metas[orderType]
// 	if !ok {
// 		return 0, ok
// 	}
// 	return meta.unsetCode, ok
// }

// func (o *Opay) MetaStep(orderType string, code int64) Step {
// 	o.metasLock.RLock()
// 	defer o.metasLock.RUnlock()
// 	meta, ok := o.metas[orderType]
// 	if !ok {
// 		return UNSET
// 	}
// 	status, ok := meta.statuses[code]
// 	if !ok {
// 		return UNSET
// 	}
// 	return status.Step
// }

// func (o *Opay) MetaNote(orderType string, code int64) string {
// 	o.metasLock.RLock()
// 	defer o.metasLock.RUnlock()
// 	meta, ok := o.metas[orderType]
// 	if !ok {
// 		return ""
// 	}
// 	status, ok := meta.statuses[code]
// 	if !ok {
// 		return ""
// 	}
// 	return status.Note
// }

func (m *Meta) OrderType() string {
	return m.orderType
}

func (m *Meta) UnsetCode() int64 {
	return m.unsetCode
}

func (m *Meta) Status(code int64) (Status, bool) {
	status, ok := m.statuses[code]
	return status, ok
}

func (m *Meta) Note(code int64) string {
	status, ok := m.Status(code)
	if !ok {
		return ""
	}
	return status.Note
}

// 执行订单处理
func (m *Meta) serve(ctx *Context) error {
	// 若为结构体类型，则创建新实例
	if m.handler.Kind() == reflect.Struct {
		m.handler = reflect.New(m.handler.Type())
	}
	return m.handler.Interface().(Handler).ServeOpay(ctx)
}
