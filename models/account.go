package models

import (
	"time"

	"github.com/imiskolee/optional"
)

// Account 定义用户表结构
type Account struct {
	ID          optional.Int64   `json:"id" gorm:"column:id; type:serial;primary_key"`
	UserID      optional.Int64   `json:"user_id" gorm:"column:user_id; type:bigint; index"`
	Type        optional.String  `json:"type" gorm:"type:varchar(255);"`
	Coin        optional.Float64 `json:"coin" gorm:"type:decimal; default 0.0"`
	FrozenCoins optional.Float64 `json:"frozen_coins" gorm:"type:decimal;  default(0.0)"`
	LockCoin    optional.Float64 `json:"lock_coin" gorm:"type:decimal;  default(0.0)"`
	LimitCoin   optional.Float64 `json:"limit_coin" gorm:"type:decimal;  default(0.0)"`
	CreatedAt   time.Time        `json:"created_at" gorm:"type:timestamp; "`
	UpdatedAt   time.Time        `json:"updated_at" gorm:"type:timestamp; "`
}

// TableName 如果定义此方法，则表名称为users，否则根据命名规则自动生成表名
func (Account) TableName() string {
	return "accounts"
}
