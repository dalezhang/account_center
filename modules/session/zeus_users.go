package session

import (
	"fmt"
	"net/url"

	log "gitee.com/dalezhang/account_center/logger"
	"gitee.com/dalezhang/account_center/util"
	"github.com/imiskolee/optional"
)

// LoadOneUserFromZeus err不为空时
func (s *Session) LoadOneUserFromZeus() (err error, hasFound bool, user ZeusUsersResp) {
	err, resp := loadUserFromZeus(s)
	if err != nil {
		return err, false, user
	}
	if len(resp.Users) == 0 {
		return
	}
	if len(resp.Users) > 1 {
		hasFound = true
		userSnapShord := ""
		for _, u := range resp.Users {
			userSnapShord = userSnapShord + fmt.Sprintf("{id: %d, name: %s, phone: %s}, ", u.ID, u.Name, u.Phone)
		}
		errString := fmt.Sprintf("users.size: %d, users: %s", len(resp.Users), userSnapShord)
		log.Logger.Error(errString)
		err = fmt.Errorf(errString)
		return
	}
	return nil, true, resp.Users[0]
}

func loadUserFromZeus(s *Session) (err error, resp data) {
	params := make(url.Values)
	dns := fmt.Sprintf("http://%s/%s", util.Config.ZeusDNS, "admin/sync/users")
	client := util.NewClient()

	if s.Phone.IsPresent() {
		params.Add("phone", s.Phone.V())
		err = client.Get(dns, params, &resp)
	} else if s.Wechat.UUID.IsPresent() {
		params.Add("wechat", s.Wechat.UUID.V())
		err = client.Get(dns, params, &resp)
	} else if s.QQ.UUID.IsPresent() {
		params.Add("qq", s.QQ.UUID.V())
		err = client.Get(dns, params, &resp)
	}
	fmt.Printf("=====params: %+v", params)

	return
}

type data struct {
	Users []ZeusUsersResp `json:"data"`
}

type ZeusUsersResp struct {
	ID                optional.Int64  `json:"id" gorm:"column:id; type:int(11) AUTO_INCREMENT primary_key;"`
	Name              optional.String `json:"name" gorm:"type:VARCHAR(255); not null "`
	PasswordDigest    optional.String `json:"password_digest" gorm:"type:VARCHAR(255); not null"`
	Phone             optional.String `json:"phone" gorm:"type:VARCHAR(255); DEFAULT NULL"`
	Token             optional.String `json:"token" gorm:"type:VARCHAR(255); DEFAULT NULL"`
	CreatedAt         optional.String `json:"created_at" gorm:"type:DATETIME;  "`
	UpdatedAt         optional.String `json:"updated_at" gorm:"type:DATETIME;  "`
	VerificationCode  optional.String `json:"verification_code" gorm:"type:varchar(255);  DEFAULT '' "`
	Wechat            optional.String `json:"wechat" gorm:"type:varchar(255);  DEFAULT NULL"`
	Weibo             optional.String `json:"weibo" gorm:"type:varchar(255);  DEFAULT NULL"`
	QQ                optional.String `json:"qq" gorm:"column:qq; type:varchar(255);  DEFAULT NULL"`
	Type              optional.String `json:"type" gorm:"type:varchar(255);  DEFAULT 'NormalUser'"`
	TokenExpire       optional.String `json:"token_expire" gorm:"type:DATETIME;  DEFAULT '2017-10-12 08:50:49'"`
	LastLogin         optional.String `json:"last_login" gorm:"type:DATETIME;  DEFAULT NULL"`
	Description       optional.String `json:"description" gorm:"type:varchar(255);  DEFAULT NULL"`
	IsBan             optional.Bool   `json:"is_ban" gorm:"type:tinyint(1);  DEFAULT '0'"`
	RealPhone         optional.String `json:"real_phone" gorm:"type:varchar(255);  DEFAULT '' "`
	WebToken          optional.String `json:"web_token" gorm:"type:varchar(255);  DEFAULT '' "`
	EmailConfirmed    optional.Bool   `json:"email_confirmed" gorm:"type:tinyint(1);  DEFAULT '0'"`
	ConfirmToken      optional.String `json:"confirm_token" gorm:"type:varchar(255);  DEFAULT NULL"`
	DialingCode       optional.String `json:"dialing_code" gorm:"type:varchar(255);  DEFAULT '86' "`
	PersonalWords     optional.String `json:"personal_words" gorm:"type:varchar(255);  DEFAULT NULL"`
	PlatformChannelID optional.Int64  `json:"platform_channel_id" gorm:"column:platform_channel_id;  type:int(11);  not null DEFAULT '1'"`
	RegisterIP        optional.String `json:"register_ip" gorm:"column:register_ip; type:varchar(255);  DEFAULT NULL"`
	Avatar            optional.String `json:"avatar" gorm:"type:VARCHAR(255);  DEFAULT NULL"`
}
