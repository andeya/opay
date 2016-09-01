package opay

import (
	"errors"
	"sync"

	"github.com/jmoiron/sqlx"
)

type (
	// 账户操作接口
	IAccount interface {
		// 修改账户余额
		// @amount 正加负减
		// @tx 数据库事务操作实例
		UpdateBalance(uid string, amount float64, tx *sqlx.Tx, values Values) error
	}

	// 账户操作接口函数
	IAccountFunc func(uid string, amount float64, tx *sqlx.Tx, values Values) error
)

var _ IAccount = IAccountFunc(nil)

func (af IAccountFunc) UpdateBalance(uid string, amount float64, tx *sqlx.Tx, values Values) error {
	return af(uid, amount, tx, values)
}

// 账户操作接口列表
type AccList struct {
	mu sync.RWMutex
	m  map[string]IAccount
}

// 全局账户操作接口列表，默认注册空资产账户空操作接口。
var globalAccList = &AccList{
	m: map[string]IAccount{
		"": new(emptyIAccount),
	},
}

// 注册账户操作接口
// @aid 资产ID
func (al *AccList) Account(aid string, accounter IAccount) error {
	al.mu.Lock()
	defer al.mu.Unlock()
	_, ok := al.m[aid]
	if ok {
		return errors.New("IAccount \"" + aid + "\" has been registered.")
	}
	al.m[aid] = accounter
	return nil
}

// 注册账户操作接口
// @aid 资产ID
func (al *AccList) AccountFunc(
	aid string,
	fn func(uid string, amount float64, tx *sqlx.Tx, values Values) error,
) error {
	return al.Account(aid, IAccountFunc(fn))
}

// 注册账户操作接口
// @aid 资产ID
func Account(aid string, acc IAccount) error {
	return globalAccList.Account(aid, acc)
}

// 注册账户操作接口
// @aid 资产ID
func AccountFunc(
	aid string,
	fn func(uid string, amount float64, tx *sqlx.Tx, values Values) error,
) error {
	return globalAccList.AccountFunc(aid, fn)
}

// 获取账户操作接口
// @aid 资产ID
func (al *AccList) GetIAccount(aid string) (IAccount, error) {
	al.mu.RLock()
	acc, ok := al.m[aid]
	al.mu.RUnlock()
	if !ok {
		return nil, errors.New("Not Found IAccount \"" + aid + "\".")
	}
	return acc, nil
}

// 账户空操作接口
type emptyIAccount int

func (*emptyIAccount) GetBalance(uid string, tx *sqlx.Tx, values Values) (float64, error) {
	return 0, nil
}

func (*emptyIAccount) UpdateBalance(uid string, amount float64, tx *sqlx.Tx, values Values) error {
	return nil
}
