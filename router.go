package opay

import (
	"errors"
	"sync"
)

type (
	// 订单处理接口
	Handler interface {
		ServeOpay(*Context) error
	}

	// 订单处理接口函数
	HandlerFunc func(*Context) error
)

func (hf HandlerFunc) ServeOpay(ctx *Context) error {
	return hf(ctx)
}

type (
	// 订单操作接口路由
	ServeMux struct {
		mu sync.RWMutex
		m  map[string]Handler
	}
)

// 注册订单处理接口
func (mux *ServeMux) Handle(key string, handler Handler) error {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	_, ok := mux.m[key]
	if ok {
		return errors.New("Handler \"" + key + "\" has been registered.")
	}
	mux.m[key] = handler
	return nil
}

// 注册订单处理接口
func (mux *ServeMux) HandleFunc(key string, handler func(*Context) error) error {
	return mux.Handle(key, HandlerFunc(handler))
}

// 通过路由执行订单处理
func (mux *ServeMux) serve(ctx *Context) error {
	mux.mu.RLock()
	h, ok := mux.m[ctx.Request.Key]
	mux.mu.RUnlock()

	if !ok {
		return errors.New("Not Found Handler")
	}

	return h.ServeOpay(ctx)
}

// 订单操作接口的全局路由
var globalServeMux = &ServeMux{
	m: make(map[string]Handler),
}

// 向全局路由注册订单处理接口
func Handle(key string, handler Handler) error {
	return globalServeMux.Handle(key, handler)
}

// 向全局路由注册订单处理接口
func HandleFunc(key string, handler func(*Context) error) error {
	return globalServeMux.HandleFunc(key, handler)
}
