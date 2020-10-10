package session

import (
	"fmt"

	"gitee.com/dalezhang/account_center/models"
	"gitee.com/dalezhang/account_center/modules/response"
	"gitee.com/dalezhang/account_center/util"
	"github.com/imiskolee/optional"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	log "gitee.com/dalezhang/account_center/logger"
)

// swagger:model Profile
type Profile struct {
	Name optional.String `json:"name"`
	// 用户个人简介
	PersonalWords optional.String `json:"personal_words"`
	// 头像链接
	Avatar optional.String `json:"avatar"`
	// 分析师简介
	Description optional.String `json:"description"`
	// 微信uuid
	Wechat optional.String `json:"wechat"`
	// qq uuid
	QQ optional.String `json:"qq"`
	// 登陆手机
	Phone optional.String `json:"phone"`
	// 区号
	DialingCode optional.String `json:"dialing_code"`
	// 密码
	Password optional.String `json:"password"`
	// 地址
	Address optional.String `json:"address"`
	//  邮编
	Postcode optional.String `json:"postcode"`
}

func (p *Profile) UpdateUser(user *models.User) (errResponse response.ErrResponse, hasError bool) {
	errResponse, hasError = p.updatePhone(user)
	if hasError {
		return
	}
	errResponse, hasError = p.updateWechat(user)
	if hasError {
		return
	}
	errResponse, hasError = p.updateQQ(user)
	if hasError {
		return
	}
	errResponse, hasError = p.updateUser(user)
	if hasError {
		return
	}
	return
}

// def update
// _profile_params = profile_params

// if params[:profile][:code].present? && params[:profile][:phone].present?
// 	if SmsCode.not_destroy_check_code(params[:profile][:code],  params[:profile][:phone], valid_dialing_code)
// 		raise ApiError.new(code: 9011, text: "该电话号码已经被使用") if  User.find_by(phone: params[:profile][:phone]).present?
// 		_profile_params = _profile_params.merge(phone: params[:profile][:phone], dialing_code: params[:profile][:dialing_code])
// 	else
// 		raise ApiError.new(code: 5002,  text: '验证码错误', status: 401)
// 	end
// end

// if params[:profile][:wechat].present?
// 	raise  ApiError.new(code: 9012, text: "该微信已经被使用") if User.find_by(wechat:  params[:profile][:wechat]).present? # => 微信
// 	_profile_params = _profile_params.merge(wechat: params[:profile][:wechat])
// end

// if params[:profile][:qq].present?
// 	raise  ApiError.new(code: 9013, text: "该QQ已经被使用") if User.find_by(qq:  params[:profile][:qq]).present? # => qq
// 	_profile_params = _profile_params.merge(qq: params[:profile][:qq])
// end

// if current_user.update!(_profile_params)
// 	render_success
// end
// end

func (p *Profile) updatePhone(user *models.User) (errResponse response.ErrResponse, hasError bool) {

	if p.Phone.IsPresent() {
		if p.DialingCode.IsBlank() {
			p.DialingCode = optional.OfString("86")
		}
		// var code models.SmsCode
		// err := code.FindSmsCode(p.Code, p.Phone, p.DialingCode)
		var user2 models.User
		if err := util.PG.Where("phone = ? and dialing_code = ?", p.Phone, p.DialingCode).Find(&user2).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				errMsg := fmt.Sprintf("find user err: %v, phone: %s, dailing_code: %s", err, p.Phone, p.DialingCode)
				log.Logger.Error(errMsg)
				errResponse.Status = 200
				errResponse.ErrMessage.Code = 9011
				errResponse.ErrMessage.Text = errMsg
				hasError = true
				return
			}
			user.Phone = p.Phone
			user.DialingCode = p.DialingCode
		} else {
			errResponse.Status = 200
			errResponse.ErrMessage.Code = 9011
			errResponse.ErrMessage.Text = "该电话号码已经被使用"
			hasError = true
			return
		}
	}
	return
}
func (p *Profile) updateWechat(user *models.User) (errResponse response.ErrResponse, hasError bool) {
	if p.Wechat.IsPresent() {
		var user2 models.User
		if err := util.PG.Where("wechat = ?", p.Wechat).Find(&user2).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				errMsg := fmt.Sprintf("find user err: %v, wechat: %s", err, p.Wechat)
				log.Logger.Error(errMsg)
				errResponse.Status = 200
				errResponse.ErrMessage.Code = 9012
				errResponse.ErrMessage.Text = errMsg
				hasError = true
				return
			}
			user.Wechat = p.Wechat
		} else {
			errResponse.Status = 200
			errResponse.ErrMessage.Code = 9012
			errResponse.ErrMessage.Text = "该微信已经被使用"
			hasError = true
			return
		}
	}
	return
}

func (p *Profile) updateQQ(user *models.User) (errResponse response.ErrResponse, hasError bool) {
	if p.QQ.IsPresent() {
		var user2 models.User
		if err := util.PG.Where("qq = ?", p.QQ).Find(&user2).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				errMsg := fmt.Sprintf("find user err: %v, qq: %s", err, p.QQ)
				log.Logger.Error(errMsg)
				errResponse.Status = 200
				errResponse.ErrMessage.Code = 9013
				errResponse.ErrMessage.Text = errMsg
				hasError = true
				return
			}
			user.QQ = p.QQ
		} else {
			errResponse.Status = 200
			errResponse.ErrMessage.Code = 9013
			errResponse.ErrMessage.Text = "该QQ已经被使用"
			hasError = true
			return
		}
	}
	return
}

func (p *Profile) updateUser(user *models.User) (errResponse response.ErrResponse, hasError bool) {
	user.PersonalWords = p.PersonalWords
	user.Description = p.Description
	user.Avatar = p.Avatar
	user.Name = p.Name
	hash, err := bcrypt.GenerateFromPassword([]byte(p.Password.V()), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	encodePW := string(hash) // 保存在数据库的密码，虽然每次生成都不同，只需保存一份即
	user.PasswordDigest = optional.OfString(encodePW)
	if err := util.PG.Model(&user).Updates(&user).Error; err != nil {
		errMsg := fmt.Sprintf("update user err: %v, user: %+v", err, user)
		log.Logger.Error(errMsg)
		errResponse.Status = 200
		errResponse.ErrMessage.Code = 9000
		errResponse.ErrMessage.Text = errMsg
		hasError = true
		return
	}
	return
}
