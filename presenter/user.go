package presenter

import (
	"time"

	"gitee.com/dalezhang/account_center/models"
	"github.com/imiskolee/optional"
)

// swagger:model UserPresenter
type UserPresenter struct {
	dbUser    *models.User
	ID        optional.Int64  `json:"id"`
	Name      optional.String `json:"name"`
	Phone     optional.String `json:"phone"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Type      optional.String `json:"type"`
	LastLogin *time.Time      `json:"last_login"`
	// 分析师简介
	Description optional.String `json:"description"`
	// 用户状态： is_ban(冻结)，normal(正常)
	Status optional.String `json:"status"`
	// 用户个人简介
	PersonalWords     optional.String `json:"personal_words"`
	PlatformChannelID optional.Int64  `json:"platform_channel_id"`
	RegisterIP        optional.String `json:"register_ip"`
	// 头像链接
	Avatar optional.String `json:"avatar"`
	// 是否有密码
	HasPassword optional.Bool `json:"has_password"`
	// 第三方登陆来源
	OmniauthProvider optional.String `json:"omniauth_provider"`
}

func (p *UserPresenter) Present(user *models.User) {
	p.dbUser = user
	p.ID = user.ID
	p.Name = user.Name
	p.Phone = user.Phone
	p.CreatedAt = user.CreatedAt
	p.UpdatedAt = user.UpdatedAt
	p.Avatar = user.Avatar
	p.PersonalWords = user.PersonalWords
	if user.PasswordDigest.V() != "" {
		p.HasPassword = optional.OfBool(true)
	}
	p.Description = user.Description
	if user.Wechat.V() != "" {
		p.OmniauthProvider = optional.OfString("wechat")
	} else if user.QQ.V() != "" {
		p.OmniauthProvider = optional.OfString("qq")
	}
	p.Status = optional.OfString(user.GetStatus())
}
