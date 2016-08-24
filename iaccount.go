package opay

import (
	"errors"
	"sync"

	"github.com/jmoiron/sqlx"
)

type (
	// 账户操作接口
	Accounter interface {
		// 获取账户余额
		// @aid 资产ID
		// @uid 用户ID
		GetBalance(uid string, tx *sqlx.Tx) (float64, error)

		// 修改账户余额
		// @aid 资产ID
		// @amount 正加负减
		// @tx 当在一个事务中时，作为数据库的操作句柄
		UpdateBalance(uid string, amount float64, tx *sqlx.Tx) error
	}

	AccList struct {
		mu sync.RWMutex
		m  map[string]Accounter
	}
)

// 账户操作接口列表
// 默认注册空资产账户空操作接口
var globalAccList = &AccList{
	m: map[string]Accounter{
		"": new(emptyAccounter),
	},
}

// 注册账户操作接口
func (al *AccList) Account(aid string, accounter Accounter) error {
	al.mu.Lock()
	defer al.mu.Unlock()
	_, ok := al.m[aid]
	if ok {
		return errors.New("Accounter \"" + aid + "\" has been registered.")
	}
	al.m[aid] = accounter
	return nil
}

// 注册账户操作接口
func Account(aid string, acc Accounter) error {
	return globalAccList.Account(aid, acc)
}

// 获取账户操作接口
func (al *AccList) GetAccounter(aid string) (Accounter, error) {
	al.mu.RLock()
	acc, ok := al.m[aid]
	al.mu.RUnlock()
	if !ok {
		return nil, errors.New("Not Found Accounter \"" + aid + "\".")
	}
	return acc, nil
}

// 账户空操作接口
type emptyAccounter int

func (*emptyAccounter) GetBalance(uid string, tx *sqlx.Tx) (float64, error) {
	return 0, nil
}

func (*emptyAccounter) UpdateBalance(uid string, amount float64, tx *sqlx.Tx) error {
	return nil
}
