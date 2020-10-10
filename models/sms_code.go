package models

import (
	"fmt"
	"regexp"
	"time"

	. "gitee.com/dalezhang/account_center/logger"
	"gitee.com/dalezhang/account_center/util"
	"github.com/imiskolee/optional"
	"github.com/jinzhu/gorm"
)

type SmsCode struct {
	ID          optional.Int64  `json:"id" gorm:"column:id; type:serial;primary_key"`
	InvalidTime *time.Time      `json:"invalid_time" gorm:"column:invalid_time; type:timestamp; "`
	Code        optional.String `json:"code" gorm:"column:code; type:varchar(255);"`
	Phone       optional.String `json:"phone" gorm:"column:phone; type:varchar(255);"`
	DialingCode optional.String `json:"dialing_code" gorm:"column:dialing_code; type:varchar(255);  DEFAULT '86' "`
	CreatedAt   time.Time       `gorm:"column:created_at; type:timestamp ;" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"column:updated_at; type:timestamp ;" json:"updated_at"`
}

func (SmsCode) TableName() string {
	return "sms_codes"
}

func (s *SmsCode) BeforeCreate(scope *gorm.Scope) error {
	reg1 := regexp.MustCompile(`\+(9[976]\d|8[987530]\d|6[987]\d|5[90]\d|42\d|3[875]\d|2[98654321]\d|9[8543210]|8[6421]|6[6543210]|5[87654321]|4[987654310]|3[9643210]|2[70]|7|1)\d{1,14}$`)
	realPhone := fmt.Sprintf("+%s%s", s.DialingCode, s.Phone)
	if reg1.MatchString(realPhone) {
		return nil
	}
	err := fmt.Errorf("手机格式错误: %s", realPhone)
	return err
}

func (c *SmsCode) FindSmsCode(code, phone, dialingCode optional.String) error {
	err := util.PG.Model(&c).Where("code = ? and phone = ? and dialing_code = ?", code, phone, dialingCode).First(&c).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		errMsg := fmt.Sprintf("find code err: %v, code: %s,phone: %s, dialing_code: %s", err, code, phone, dialingCode)
		Logger.Error(errMsg)
	}
	if err == gorm.ErrRecordNotFound {
		return fmt.Errorf("验证失败")
	}
	return nil
}

func (c *SmsCode) FindAndDestorySmsCode(code, phone, dialingCode optional.String) {
	err := util.PG.Model(&c).Where("code = ? and phone = ? and dialing_code = ?", code, phone, dialingCode).Find(&c).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		errMsg := fmt.Sprintf("find code err: %v, code: %s,phone: %s, dialing_code: %s", err, code, phone, dialingCode)
		Logger.Error(errMsg)
	}
	if err == gorm.ErrRecordNotFound {
		return
	}
	err = util.PG.Model(&c).Delete(&c).Error
	if err != nil {
		errMsg := fmt.Sprintf("Delete code err: %v, code: %s,phone: %s, dialing_code: %s", err, code, phone, dialingCode)
		Logger.Error(errMsg)
	}
	return
}
