package models

import (
	"fmt"

	log "gitee.com/dalezhang/account_center/logger"
	"gitee.com/dalezhang/account_center/util"
	"github.com/imiskolee/optional"
	"github.com/jinzhu/gorm"
)

// swagger:model Address
type Address struct {
	ID     optional.Int64 `json:"id" gorm:"column:id; type:serial;primary_key"`
	UserID optional.Int64 `json:"user_id" gorm:"column:user_id; type:bigint; index"`
	// 收件人
	Name optional.String `json:"name" gorm:"type:VARCHAR(255); "`
	// 邮编
	Postcode optional.String `json:"postcode" gorm:"type:VARCHAR(255);"`
	// 电话
	Phone optional.String `json:"phone" gorm:"type:VARCHAR(255);"`
	// 省
	Province optional.String `json:"province" gorm:"type:varchar(255);"`
	// 市
	City optional.String `json:"city" gorm:"type:varchar(255);"`
	// 区
	Zone optional.String `json:"zone" gorm:"type:varchar(255);"`
	// 详细地址
	House optional.String `json:"house" gorm:"type:varchar(255);"`
	// 标记
	Remark optional.String `json:"remark" gorm:"type:varchar(255)"`
	// 默认
	IsDefault optional.Bool `json:"is_default" gorm:"type:boolean;  default(false)"`
}

func (a *Address) BeforeCreate(scope *gorm.Scope) error {
	a.changeIsDefault()
	return nil
}

func (a *Address) BeforeUpdate(scope *gorm.Scope) error {
	a.changeIsDefault()
	return nil
}

func (a *Address) changeIsDefault() {
	if a.IsDefault.V() == false {
		return
	}
	var originDefault Address
	err := util.PG.Model(&originDefault).Where("user_id = ? and is_default = ?", a.UserID.V(), true).Find(&originDefault).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		errMsg := fmt.Sprintf("Get users address err: %v,", err)
		log.Logger.Error(errMsg)
	}
	if err == gorm.ErrRecordNotFound {
		return
	}
	originDefault.IsDefault = optional.OfBool(false)
	if err := util.PG.Model(&originDefault).Updates(&originDefault).Error; err != nil {
		errMsg := fmt.Sprintf("update address err: %v, user: %+v", err, originDefault)
		log.Logger.Error(errMsg)
	}
	return
}
