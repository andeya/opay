package base

import (
	"errors"

	"github.com/henrylee2cn/opay"
	"github.com/jmoiron/sqlx"
)

type BaseAccount struct {
	Uid       string `json:"uid" db:"uid"`     //user id
	Aid       string `json:"aid" db:"aid"`     //asset id
	Total     string `json:"total" db:"total"` //account balance
	Ip        string `json:"ip" db:"ip"`
	UpdatedAt int64  `json:"updated_at" db:"updated_at"`
	tableName string //database table name
}

var _ opay.IAccount = new(BaseAccount)

func (this *BaseAccount) UpdateBalance(uid string, amount float64, tx *sqlx.Tx, values opay.Values) error {
	return errors.New("*BaseAccount does not implement opay.IAccount (missing UpdateBalance method).")
}

func (this *BaseAccount) TableName() string {
	return this.tableName
}

func (this *BaseAccount) SetTableName(tableName string) {
	this.tableName = tableName
}
