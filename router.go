package opay

import (
	"errors"
	"reflect"
	"sync"
)

type (
	// 订单处理接口
	// 只允许函数或结构体类型
	Handler interface {
		ServeOpay(*Context) error
	}

	// 订单处理接口函数
	HandlerFunc func(*Context) error
)

var _ Handler = HandlerFunc(nil)

func (hf HandlerFunc) ServeOpay(ctx *Context) error {
	return hf(ctx)
}

type (
	// 订单操作接口路由
	ServeMux struct {
		mu sync.RWMutex
		m  map[string]reflect.Value
	}
)

// 注册订单处理接口
func (mux *ServeMux) Handle(operator string, handler Handler) error {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	_, ok := mux.m[operator]
	if ok {
		return errors.New("Handler \"" + operator + "\" has been registered.")
	}

	v := reflect.ValueOf(handler)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 过滤不允许的类型
	if !(v.Kind() == reflect.Struct || v.Kind() == reflect.Func) {
		return errors.New("Handler must be func or struct type.")
	}

	mux.m[operator] = v
	return nil
}

// 注册订单处理接口
func (mux *ServeMux) HandleFunc(operator string, fn func(*Context) error) error {
	return mux.Handle(operator, HandlerFunc(fn))
}

// 通过路由执行订单处理
func (mux *ServeMux) serve(ctx *Context) error {
	mux.mu.RLock()
	v, ok := mux.m[ctx.Operator()]
	mux.mu.RUnlock()

	if !ok {
		return errors.New("Not Found Handler")
	}

	// 若为结构体类型，则创建新实例
	if v.Kind() == reflect.Struct {
		v = reflect.New(v.Type())
	}
	return v.Interface().(Handler).ServeOpay(ctx)
}

// 检查指定操作符是否存在
func (mux *ServeMux) CheckOperator(operator string) error {
	mux.mu.RLock()
	_, ok := mux.m[operator]
	mux.mu.RUnlock()
	if !ok {
		return errors.New("Not Found Handler")
	}
	return nil
}

// 订单操作接口的全局路由
var globalServeMux = &ServeMux{
	m: make(map[string]reflect.Value),
}

// 向全局路由注册订单处理接口
func Handle(operator string, handler Handler) error {
	return globalServeMux.Handle(operator, handler)
}

// 向全局路由注册订单处理接口
func HandleFunc(operator string, handler func(*Context) error) error {
	return globalServeMux.HandleFunc(operator, handler)
}
