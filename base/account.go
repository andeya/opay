package base

import (
	"errors"

	"github.com/henrylee2cn/opay"
	"github.com/jmoiron/sqlx"
)

type BaseAccount struct {
	tableName string //database table name
}

var _ opay.IAccount = new(BaseAccount)

// 账户总额变动
// amount正负代表收支
// 调用者需为字段完整的对象指针
func (this *BaseAccount) UpdateBalance(uid string, amount float64, tx *sqlx.Tx, values opay.Values) error {
	return errors.New("*BaseAccount does not implement opay.IAccount (missing UpdateBalance method).")
}
