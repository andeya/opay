package opay

import (
	"errors"
	"sync"

	"github.com/jmoiron/sqlx"
)

// 账户余额操作函数
type SettlementFunc func(uid string, amount float64, tx *sqlx.Tx, values Values) error

// 账户余额操作函数路由
type SettlementMux struct {
	mu sync.RWMutex
	m  map[string]SettlementFunc
}

// 获取账户余额操作函数
// @aid 资产ID
func (this *SettlementMux) GetSettlementFunc(aid string) (SettlementFunc, error) {
	this.mu.RLock()
	acc, ok := this.m[aid]
	this.mu.RUnlock()
	if !ok {
		return nil, errors.New("Not Found SettlementFunc \"" + aid + "\".")
	}
	return acc, nil
}

// 注册账户余额操作函数
// @aid 资产ID
func (this *SettlementMux) RegSettlementFunc(aid string, fn SettlementFunc) error {
	this.mu.Lock()
	defer this.mu.Unlock()
	_, ok := this.m[aid]
	if ok {
		return errors.New("SettlementFunc \"" + aid + "\" has been registered.")
	}
	this.m[aid] = fn
	return nil
}

// 全局账户操作接口列表，默认注册空资产账户空操作接口。
var globalSettlementMux = &SettlementMux{
	m: map[string]SettlementFunc{
		"": emptySettlement,
	},
}

// 注册账户余额操作函数
// @aid 资产ID
func RegSettlementFunc(aid string, acc SettlementFunc) error {
	return globalSettlementMux.RegSettlementFunc(aid, acc)
}

// Empty Settlement Function of empty asset.
func emptySettlement(uid string, amount float64, tx *sqlx.Tx, values Values) error {
	return errors.New("Empty Settlement Function.")
}
