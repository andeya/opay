package opay

import (
	"errors"
	"sync"

	"github.com/jmoiron/sqlx"
)

// SettleFunc: Account balance operation function.
type SettleFunc func(uid string, amount float64, tx *sqlx.Tx) error

// SettleFuncMap: Account Balance Operations Function Router.
type SettleFuncMap struct {
	mu sync.RWMutex
	m  map[string]SettleFunc
}

// GetSettleFunc gets the account balance operation function
// @aid Assets ID
func (this *SettleFuncMap) GetSettleFunc(aid string) (SettleFunc, error) {
	this.mu.RLock()
	acc, ok := this.m[aid]
	this.mu.RUnlock()
	if !ok {
		return nil, errors.New("opay: not found SettleFunc '" + aid + "'.")
	}
	return acc, nil
}

// RegSettleFunc registers the account balance operation function.
// @aid Assets ID
func (this *SettleFuncMap) RegSettleFunc(aid string, fn SettleFunc) error {
	this.mu.Lock()
	defer this.mu.Unlock()
	_, ok := this.m[aid]
	if ok {
		return errors.New("opay: settleFunc '" + aid + "' has been registered.")
	}
	this.m[aid] = fn
	return nil
}

// Global account operation interface list, the default registered empty asset account empty operation interface.
var globalSettleFuncMap = &SettleFuncMap{
	m: map[string]SettleFunc{
		"": emptySettle,
	},
}

// RegSettleFunc registers the account balance operation function.
// @aid Assets ID
func RegSettleFunc(aid string, acc SettleFunc) error {
	return globalSettleFuncMap.RegSettleFunc(aid, acc)
}

// Empty Settle Function of empty asset.
func emptySettle(uid string, amount float64, tx *sqlx.Tx) error {
	return errors.New("opay: empty settle function.")
}
