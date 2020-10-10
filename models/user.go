package models

import (
	"fmt"
	"time"

	. "gitee.com/dalezhang/account_center/logger"
	"gitee.com/dalezhang/account_center/util"
	"github.com/imiskolee/optional"
	"github.com/jinzhu/gorm"
)

// User 定义用户表结构
type User struct {
	ID                optional.Int64  `json:"id" gorm:"column:id; type:serial;primary_key;"`
	MysqlUserID       optional.Int64  `json:"mysql_user_id" gorm:"column:mysql_user_id; type:bigint; index"`
	Name              optional.String `json:"name" gorm:"type:VARCHAR(255); "`
	PasswordDigest    optional.String `json:"password_digest" gorm:"type:VARCHAR(255);"`
	Phone             optional.String `json:"phone" gorm:"type:VARCHAR(255);"`
	CreatedAt         time.Time       `json:"created_at" gorm:"type:timestamp;"`
	UpdatedAt         time.Time       `json:"updated_at" gorm:"type:timestamp;"`
	VerificationCode  optional.String `json:"verification_code" gorm:"type:varchar(255);"`
	Wechat            optional.String `json:"wechat" gorm:"type:varchar(255);"`
	Weibo             optional.String `json:"weibo" gorm:"type:varchar(255);"`
	QQ                optional.String `json:"qq" gorm:"column:qq; type:varchar(255);"`
	Type              optional.String `json:"type" gorm:"type:varchar(255);  DEFAULT 'NormalUser'"`
	DeletedAt         *time.Time      `gorm:"type:timestamp;"`
	LastLogin         *time.Time      `json:"last_login" gorm:"type:timestamp;"`
	Description       optional.String `json:"description" gorm:"type:varchar(255);"`
	Status            optional.Int64  `json:"Status" gorm:"type:int;  DEFAULT '0'"`
	EmailConfirmed    optional.Bool   `json:"email_confirmed" gorm:"type:boolean;"`
	DialingCode       optional.String `json:"dialing_code" gorm:"type:varchar(255);  DEFAULT '86' "`
	PersonalWords     optional.String `json:"personal_words" gorm:"type:varchar(255);"`
	PlatformChannelID optional.Int64  `json:"platform_channel_id" gorm:"column:platform_channel_id;  type:int;  not null DEFAULT '1'"`
	RegisterIP        optional.String `json:"register_ip" gorm:"column:register_ip; type:varchar(255);"`
	Avatar            optional.String `json:"avatar" gorm:"type:VARCHAR(2040);"`
	RecommendUserID   optional.Int64  `json:"user_id" gorm:"column:user_id; type:bigint;"`
	FirstBindingPhone optional.Bool   `json:"first_binding_phone" gorm:"-"`
}

// TableName 如果定义此方法，则表名称为users，否则根据命名规则自动生成表名
func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	u.setBindingPhone()
	return nil
}

func (u *User) AfterCreate(scope *gorm.Scope) error {
	u.buildSelfAccount()
	u.checkBindingPhone()
	return nil
}

func (u *User) BeforeUpdate(scope *gorm.Scope) error {
	u.setBindingPhone()
	return nil
}

func (u *User) AfterUpdate(scope *gorm.Scope) error {
	u.checkBindingPhone()
	return nil
}

func (u *User) Validate() (errors []string) {
	var i int
	var count int
	if u.ID.V() == 0 {
		err := util.PG.Model(&u).Where("phone = ? and dialing_code = ?", u.Phone, u.DialingCode).Count(&count).Error
		if err != nil {
			i = len(errors)
			errors[i] = err.Error()
		}
		if count > 0 {
			i = len(errors)
			errors[i] = "phone已经存在"
		}
	}
	return
}

func (u *User) buildSelfAccount() {
	var account Account
	account.Type = optional.OfString("NormalAccount")
	account.UserID = u.ID
	if err := util.PG.Model(&account).Create(&account).Error; err != nil || account.ID.V() == 0 {
		err = fmt.Errorf("\n account_center创建失败：%+v, account_center: %+v", err, account)
		Logger.Error(err)
	}
}

func (u *User) riskDetectionPassed() bool {
	// def similar_risk_users
	// begin_time = 1.days.before(self.created_at).beginning_of_day
	// after_time = 1.days.after(self.created_at).end_of_day

	// User.where.not(id: self.id).where("created_at > ? AND created_at < ?", begin_time, after_time).where(register_ip: self.register_ip).where.not(phone: nil)
	// end
	begingTime := u.CreatedAt.AddDate(0, 0, -1)
	afterTime := u.CreatedAt.AddDate(0, 0, 1)
	var total int
	err := util.PG.Model(&u).Where("created_at > ? and created_at < ? and register_ip = ? and phone is not null", begingTime, afterTime, u.RegisterIP).Count(&total).Error
	if err != nil {
		errMsg := fmt.Sprintf("Count users err: %v,", err)
		Logger.Error(errMsg)
		return true
	}
	if total < 4 {
		return true
	}
	return false
}

func (u *User) setBindingPhone() {
	// def set_binding_phone
	// if self.changes["phone"].present? && self.changes["phone"].first.blank? && self.changes["phone"].last.present?
	// 	self.created_at = Time.now
	// 	self.first_binding_phone = true
	// end
	// end
	if u.ID.V() == 0 && u.Phone.V() != "" {
		u.FirstBindingPhone = optional.OfBool(true)
	} else {
		if u.Phone.V() != "" && u.ID.V() != 0 {
			var originUser User
			err := util.PG.Model(&originUser).Find(&originUser, u.ID.V()).Error
			if err != nil {
				errMsg := fmt.Sprintf("Get users err: %v,", err)
				Logger.Error(errMsg)
			}
			if originUser.ID.V() != 0 && originUser.Phone.V() == "" {
				u.FirstBindingPhone = optional.OfBool(true)
			}
		}
	}
}

func (u *User) checkBindingPhone() {
	// def check_binding_phone
	// if self.first_binding_phone && risk_detection_passed?
	// 	release_registration_red_packet

	// 	self.first_binding_phone = nil
	// end
	// end

	if u.FirstBindingPhone.V() && u.riskDetectionPassed() {
		// TODO release_registration_red_packet
	}
}
func (u *User) Account() (account Account) {
	err := util.PG.Model(&account).Where("user_id = ?", u.ID).Find(&account).Error
	if err != nil {
		errMsg := fmt.Sprintf("Get users account err: %v,", err)
		Logger.Error(errMsg)
	}
	return
}

func (u *User) GetStatus() string {
	switch u.Status.V() {
	case 0:
		return "normal"
	case 1:
		return "is_ban"
	}
	return ""
}

func (u *User) SetStatus(st string) int {
	switch st {
	case "normal":
		u.Status = optional.OfInt64(0)
		return 0
	case "is_ban":
		u.Status = optional.OfInt64(1)
		return 1
	}
	return 999
}
