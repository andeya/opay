package opay

import (
	"errors"
	"sync"

	"github.com/jmoiron/sqlx"
)

// 账户余额操作函数
type SettleFunc func(uid string, amount float64, tx *sqlx.Tx, ctxStore CtxStore) error

// 账户余额操作函数路由
type SettleFuncMap struct {
	mu sync.RWMutex
	m  map[string]SettleFunc
}

// 获取账户余额操作函数
// @aid 资产ID
func (this *SettleFuncMap) GetSettleFunc(aid string) (SettleFunc, error) {
	this.mu.RLock()
	acc, ok := this.m[aid]
	this.mu.RUnlock()
	if !ok {
		return nil, errors.New("Not found SettleFunc '" + aid + "'.")
	}
	return acc, nil
}

// 注册账户余额操作函数
// @aid 资产ID
func (this *SettleFuncMap) RegSettleFunc(aid string, fn SettleFunc) error {
	this.mu.Lock()
	defer this.mu.Unlock()
	_, ok := this.m[aid]
	if ok {
		return errors.New("SettleFunc '" + aid + "' has been registered.")
	}
	this.m[aid] = fn
	return nil
}

// 全局账户操作接口列表，默认注册空资产账户空操作接口。
var globalSettleFuncMap = &SettleFuncMap{
	m: map[string]SettleFunc{
		"": emptySettle,
	},
}

// 注册账户余额操作函数
// @aid 资产ID
func RegSettleFunc(aid string, acc SettleFunc) error {
	return globalSettleFuncMap.RegSettleFunc(aid, acc)
}

// Empty Settle Function of empty asset.
func emptySettle(uid string, amount float64, tx *sqlx.Tx, ctxStore CtxStore) error {
	return errors.New("Empty settle function.")
}
