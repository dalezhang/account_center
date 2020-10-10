package session

import (
	"fmt"
	"time"

	. "gitee.com/dalezhang/account_center/logger"
	"gitee.com/dalezhang/account_center/models"
	"gitee.com/dalezhang/account_center/modules/response"
	"gitee.com/dalezhang/account_center/util"
	"github.com/imiskolee/optional"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// swagger:model Session
type Session struct {
	// swagger:ignore
	Channel optional.String `json:"channel"`
	// 区号
	DialingCode optional.String `json:"dialing_code"`
	// 姓名
	Name optional.String `json:"name"`
	// 密码
	Password optional.String `json:"password"`
	// 电话
	Phone optional.String `json:"phone"`
	// qq
	QQ thirdParty `json:"qq"`
	// 微信
	Wechat thirdParty `json:"wechat"`
	// 客户端ip
	RemoteIP optional.String `json:"remote_ip"`
	// 推荐人ID
	RecommendUserID optional.String   `json:"recommend_user_id"`
	Header          map[string]string `json:"-"`
}

type thirdParty struct {
	UUID optional.String `json:"uuid"`
	// 头像连接
	Avatar optional.String `json:"avatar"`
}

func (session *Session) Authenticate() (user models.User, errResponse response.ErrResponse, hasError bool) {
	user, errResponse, hasError = loadUser(session)

	if hasError {
		return
	}
	if user.ID.V() == 0 { //  super code: 9001, text: '用户不存在', status: 200
		errResponse.Status = 200
		errResponse.ErrMessage.Code = 9001
		errResponse.ErrMessage.Text = "用户不存在"
		hasError = true
		return
	}
	if user.GetStatus() == "is_ban" { //  super code: 1099, text: '您的账户存在一定风险，请联系客服及时处理', status: 200
		errResponse.Status = 200
		errResponse.ErrMessage.Code = 1099
		errResponse.ErrMessage.Text = "您的账户存在一定风险，请联系客服及时处理"
		hasError = true
		return
	}
	errResponse, hasError = verifyUser(&user, session)
	return
}

func loadUser(session *Session) (user models.User, errResponse response.ErrResponse, hasError bool) {
	var err error
	// 有电话号码时，建立或者关联旧user，没有时创建一个没有phone的user
	if session.Phone.V() != "" {
		// 从项目数据库找
		err = util.PG.Model(&user).Where("phone = ?", session.Phone).Find(&user).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			errMsg := fmt.Sprintf("find user err: %v, phone: %s ", err, session.Phone)
			Logger.Error(errMsg)
			//   super code: 9003, text: '用户不存在', status: 200
			errResponse.Status = 200
			errResponse.ErrMessage.Code = 9003
			errResponse.ErrMessage.Text = errMsg
			// 有异常直接返回
			return user, errResponse, true
		}
		if user.ID.V() != 0 { // 找到后返回
			return
		}
	}
	if session.Wechat.UUID.V() != "" {
		err = util.PG.Model(&user).Where("wechat = ?", session.Wechat.UUID).Find(&user).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			errMsg := fmt.Sprintf("find user err: %v, wechat uuid: %s", err, session.Wechat.UUID)
			Logger.Error(errMsg)
			errResponse.Status = 200
			errResponse.ErrMessage.Code = 9003
			errResponse.ErrMessage.Text = errMsg
			return user, errResponse, true
		}
		if user.ID.V() != 0 {
			return
		}
	}
	if session.QQ.UUID.V() != "" {
		err = util.PG.Model(&user).Where("qq = ?", session.QQ.UUID).Find(&user).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			errMsg := fmt.Sprintf("find user err: %v, qq uuid: %s", err, session.QQ.UUID)
			Logger.Error(errMsg)
			errResponse.Status = 200
			errResponse.ErrMessage.Code = 9003
			errResponse.ErrMessage.Text = errMsg
			return user, errResponse, true
		}
		if user.ID.V() != 0 {
			return
		}
	}
	err, hasFound, user := loadFromZeus(session)
	if err != nil {
		errResponse.Status = 200
		errResponse.ErrMessage.Code = 9003
		errResponse.ErrMessage.Text = err.Error()
		return user, errResponse, true
	}
	if hasFound {
		return
	}
	user, errResponse, hasError = createNewUser(session)
	return
}

func loadFromZeus(session *Session) (err error, hasFound bool, user models.User) {
	err, hasfound, zeusUser := session.LoadOneUserFromZeus()
	if err != nil { // 有异常直接返回
		err = fmt.Errorf("find user err: %v, phone: %s ", err, session.Phone)
		return err, false, user
	}
	if hasfound { // 没找到,继续用其他条件找，或者在最后应该新建一个
		user.Name = zeusUser.Name
		user.MysqlUserID = zeusUser.ID
		user.PasswordDigest = zeusUser.PasswordDigest
		user.Phone = zeusUser.Phone
		user.VerificationCode = zeusUser.VerificationCode
		user.Wechat = zeusUser.Wechat
		user.Weibo = zeusUser.Weibo
		user.QQ = zeusUser.QQ
		user.Type = zeusUser.Type
		user.Description = zeusUser.Description
		if zeusUser.IsBan.V() {
			user.SetStatus("is_ban")
		}
		user.DialingCode = zeusUser.DialingCode
		user.PersonalWords = zeusUser.PersonalWords
		user.PlatformChannelID = zeusUser.PlatformChannelID
		user.RegisterIP = zeusUser.RegisterIP
		user.Avatar = zeusUser.Avatar
		err := util.PG.Model(&user).Create(&user).Error
		if err != nil || user.ID.V() == 0 {
			errMsg := fmt.Sprintf("Insert user err: %v, user: %+v", err, user)
			Logger.Error(errMsg)
			return err, false, user
		}
		hasFound = true
	}
	return
}

func createNewUser(session *Session) (user models.User, errResponse response.ErrResponse, hasError bool) {
	var total int
	// recent_register_count = ::User.where(register_ip: remote_ip).where('created_at > ?', Time.now - 1.day).count
	if session.RemoteIP.V() != "" {
		n := time.Now()
		yesterday := n.AddDate(0, 0, -1)
		err := util.PG.Model(&user).Where("register_ip = ? and created_at > ?", session.RemoteIP, yesterday).Count(&total).Error
		if err != nil {
			errMsg := fmt.Sprintf("find user err: %v, remote_ip: %s", err, session.RemoteIP)
			Logger.Error(errMsg)
			errResponse.Status = 200
			errResponse.ErrMessage.Code = 2014
			errResponse.ErrMessage.Text = errMsg
			return user, errResponse, true
		}
		if total > 5 { //   super code: 2014, text: '今天注册的太多了', status: 200
			errResponse.Status = 200
			errResponse.ErrMessage.Code = 2014
			errResponse.ErrMessage.Text = "今天注册的太多了"
			return user, errResponse, true
		}
	}
	if session.Password.V() != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(session.Password.V()), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
		}

		encodePW := string(hash) // 保存在数据库的密码，虽然每次生成都不同，只需保存一份即
		user.PasswordDigest = optional.OfString(encodePW)
	}
	if session.Phone.V() != "" { //   code = self.find_by(code: code, phone: phone, dialing_code: dialing_code)
		user.DialingCode = validDialingCode(session)
		user.Name = session.Name
		user.Phone = session.Phone
		user.PlatformChannelID = getChanneID(session)
		user.RegisterIP = session.RemoteIP

		err := util.PG.Model(&user).Create(&user).Error
		if err != nil || user.ID.V() == 0 {
			errMsg := fmt.Sprintf("Insert user err: %v, phone: %s, dialing_code: %s", err, session.Phone, session.DialingCode)
			Logger.Error(errMsg)
		} else {
			return user, errResponse, false
		}
	}
	if session.QQ.UUID.V() != "" {
		user.QQ = session.QQ.UUID
		user.Name = session.Name
		user.PlatformChannelID = getChanneID(session)
		user.RegisterIP = session.RemoteIP
		err := util.PG.Model(&user).Create(&user).Error
		if err != nil {
			errMsg := fmt.Sprintf("Insert user err: %v, QQ: %s", err, session.QQ.UUID)
			Logger.Error(errMsg)
		} else {
			return user, errResponse, false
		}
	}
	if session.Wechat.UUID.V() != "" {
		user.Wechat = session.Wechat.UUID
		user.Name = session.Name
		user.PlatformChannelID = getChanneID(session)
		user.RegisterIP = session.RemoteIP
		err := util.PG.Model(&user).Create(&user).Error

		if err != nil || user.ID.V() == 0 {
			errMsg := fmt.Sprintf("Insert user err: %v, Wechat: %s", err, session.Wechat.UUID)
			Logger.Error(errMsg)
		} else {
			return user, errResponse, false
		}
	}
	return
}

// def load_user
// # 获得现有用户
// self.user = User.find_by(phone: self.phone) if self.phone                           # => 手机号登录
// self.user = User.find_by(wechat: self.wechat[:uuid]) if (self.wechat && !self.user) # => 微信登陆
// self.user = User.find_by(weibo: self.weibo[:uuid]) if (self.weibo && !self.user)    # => 微博登陆
// self.user = User.find_by(qq: self.qq[:uuid]) if (self.qq && !self.user)             # => qq登录
// return self.user if self.user

// # 创建新用户
// recent_register_count = ::User.where(register_ip: remote_ip).where('created_at > ?', Time.now - 1.day).count
// return REGISTEROVERLIMIT unless 5 > recent_register_count
// return self.user = User.find_or_create_by!(dialing_code: valid_dialing_code, phone: self.phone, invite_code: invite_code, platform_channel: channel, register_ip: remote_ip, name: name) if (self.phone && Sms::Code.not_destroy_check(self.code, self.phone, valid_dialing_code))  # => 手机号
// return self.user = User.find_or_create_by!(qq: self.qq[:uuid], invite_code: invite_code, platform_channel: channel, register_ip: remote_ip, name: name) if (self.qq && self.qq[:uuid].present?)                                             # => qq
// return self.user = User.find_or_create_by!(wechat: self.wechat[:uuid], invite_code: invite_code, platform_channel: channel, register_ip: remote_ip, name: name) if (self.wechat && self.wechat[:uuid].present?)                             # => 微信
// end

func validDialingCode(session *Session) optional.String {
	if session.DialingCode.V() != "" {
		return session.DialingCode
	}
	return optional.OfString("86")
}

// def valid_dialing_code
// self.dialing_code.nil? ? '86' : self.dialing_code # app兼容, 老版app无输入区号选项
// end

func getChanneID(session *Session) optional.Int64 {
	var platformChannel models.PlatformChannel
	if session.Header["Platform"] != "" {
		err := util.PG.Model(&platformChannel).Where("platform = ? and channel = ? ", session.Header["Platform"], session.Header["Channel"]).Find(&platformChannel).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			errMsg := fmt.Sprintf("find PlatformChannel err: %v, platform: %s,channel: %s", err, session.Header["Platform"], session.Header["Channel"])
			Logger.Error(errMsg)
		}
		if platformChannel.ID.V() != 0 {
			return platformChannel.ID
		}
		platformChannel.Platform = optional.OfString(session.Header["Platform"])
		platformChannel.Channel = optional.OfString(session.Header["Channel"])
		err = util.PG.Model(&platformChannel).Create(&platformChannel).Error
		if err != nil || platformChannel.ID.V() == 0 {
			errMsg := fmt.Sprintf("Insert PlatformChannel err: %v, platform: %s, channel: %s", err, session.Header["Platform"], session.Header["Channel"])
			Logger.Error(errMsg)
		} else {
			return platformChannel.ID
		}
	}
	return optional.OfInt64(0)
}

// def channel
// raise UnknownRequestError unless request.headers['Platform'].presence
// Platform::Channel.find_or_create_by(
// 	platform: request.headers['Platform'].downcase,
// 	channel: request.headers['Channel'].to_s.downcase)
// end

func verifyUser(user *models.User, session *Session) (errResponse response.ErrResponse, hasError bool) {
	if session.QQ.Avatar.V() != "" {
		user.Avatar = session.QQ.Avatar
		util.PG.Model(&user).Save(&user)
		return
	}
	if session.Wechat.Avatar.V() != "" {
		user.Avatar = session.Wechat.Avatar
		util.PG.Model(&user).Save(&user)
		return
	}
	// if session.Code.V() != "" {
	// 	var code models.SmsCode
	// 	err := code.FindSmsCode(session.Code, session.Phone, session.DialingCode)
	// 	if err != nil { // super code: 5002, text: '验证码错误', status: 200
	// 		errResponse.Status = 200
	// 		errResponse.ErrMessage.Code = 5002
	// 		errResponse.ErrMessage.Text = "验证码错误"
	// 		return errResponse, true
	// 	}
	// }
	if user.ID.V() == 0 { // super code: 9001, text: '用户不存在', status: 200
		errResponse.Status = 200
		errResponse.ErrMessage.Code = 9001
		errResponse.ErrMessage.Text = "用户不存在"
		return errResponse, true
	}
	// if session.Code.V() != "" {
	// 	var code models.SmsCode
	// 	code.FindAndDestorySmsCode(session.Code, session.Phone, session.DialingCode)
	// }
	// if user.PasswordDigest.V() == "" { // super code: 2003, text: '用户没有密码', status: 200
	// 	errResponse.Status = 200
	// 	errResponse.ErrMessage.Code = 2003
	// 	errResponse.ErrMessage.Text = "用户没有密码"
	// 	return errResponse, true
	// }
	if session.Password.V() != "" {
		errResponse, hasError = authUser(user, session)
	}
	return
}

// def verify_user
// return self.user if ( (self.user.update_attribute(:avatar, self.qq[:touxiang]) if self.qq) or self.qq.present?)             # => qq校验
// return self.user if ( (self.user.update_attribute(:avatar, self.wechat[:touxiang]) if self.wechat) or self.wechat.present?) # => 微信校验
// return SMSCODENOTCORRECT if (self.code && !Sms::Code.not_destroy_check(self.code, self.phone, valid_dialing_code))          # => 验证码错误
// return USERNOTEXIST if self.user.blank?                                                                                     # => 用户不存在错误
// return (self.user.authenticate_code(self.code) ? self.user : SMSCODENOTCORRECT) unless self.code.blank?                     # => 验证码校验及删除
// return USERHASNOPASSWORD if self.user.password_digest.blank?                                                                # => 密码存在校验;  如果没有密码,返回 2
// return (auth_user ? self.user : USERPASSWORDERROR) unless self.password.blank?                                              # => 密码校验
// end

func authUser(user *models.User, session *Session) (errResponse response.ErrResponse, hasError bool) {
	if session.Password.V() != "" {
		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest.V()), []byte(session.Password.V()))
		if err != nil {
			//   super code: 3005, text: '密码错误', status: 200
			errResponse.Status = 200
			errResponse.ErrMessage.Code = 3005
			errResponse.ErrMessage.Text = "密码错误"
			return errResponse, true
		}
	}
	return
}

// def auth_user
// self.user.authenticate(self.password) && self.user.authenticate_dialing_code(valid_dialing_code)
// end
