package session

import (
	"encoding/json"
	"fmt"
	"testing"

	"gitee.com/dalezhang/account_center/helpers/secure_random"
	. "gitee.com/dalezhang/account_center/logger"
	"gitee.com/dalezhang/account_center/models"
	"gitee.com/dalezhang/account_center/presenter"
	"gitee.com/dalezhang/account_center/util"

	"github.com/imiskolee/optional"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {

	util.InitPG(Logger)

}

func TestCreateUser(t *testing.T) {

	Convey("create user by qq", t, func(c C) {
		var s Session
		uuid1 := uuid.NewV4()
		s.QQ.UUID = optional.OfString(uuid1.String())
		s.QQ.Avatar = optional.OfString("touxiang")

		// var user models.User
		// util.PG.Where("phone = ?", smsCode.Phone).Delete(&user)
		user, errResponse, hasError := s.Authenticate()
		if hasError {
			fmt.Printf("errResponse: %v \n", errResponse)
		}

		c.So(hasError, ShouldBeFalse)
		c.So(user.Avatar.V(), ShouldEqual, "touxiang")
	})
	Convey("create user by phone", t, func(c C) {
		var users []models.User
		var user models.User
		err := util.PG.Model(&user).Find(&users).Error
		if err != nil {
			panic(err)
		}
		for i, u := range users {
			fmt.Printf("user: %d, %d, phone: %s \n", u.ID, i, u.Phone)
		}

		if err := util.PG.Model(&user).Where("phone = '13800138000'").First(&user).Error; err == nil {
			if err := util.PG.Unscoped().Delete(&user).Error; err != nil {
				panic(err)
			}
		} else if err != gorm.ErrRecordNotFound {
			panic(err)
		}
		password := secure_random.NewString(8)
		var s Session
		s.Phone = optional.OfString("13800138000")
		s.DialingCode = optional.OfString("86")
		s.Password = optional.OfString(password)
		s.Header = map[string]string{"Platform": "android"}

		user, errResponse, hasError := s.Authenticate()
		if hasError {
			fmt.Printf("errResponse: %v \n", errResponse)
		}
		var p presenter.UserPresenter
		p.Present(&user)
		j, _ := json.Marshal(p)
		fmt.Printf("p: %+v \n", string(j))
		c.So(hasError, ShouldBeFalse)
		c.So(user.Account().ID.V(), ShouldNotEqual, 0)
		c.So(user.PlatformChannelID.V(), ShouldNotEqual, 0)
		c.So(p.ID.V(), ShouldNotEqual, 0)
		var s2 Session
		s2.Phone = optional.OfString("13800138000")
		s2.DialingCode = optional.OfString("86")
		s2.Password = optional.OfString(password)
		s2.Header = map[string]string{"Platform": "android"}

		user, errResponse, hasError = s2.Authenticate()
		if hasError {
			fmt.Printf("\n errResponse: %v \n", errResponse)
		}
		j, _ = json.Marshal(p)
		fmt.Printf("p: %+v \n", string(j))
		c.So(hasError, ShouldBeFalse)
		c.So(user.Account().ID.V(), ShouldNotEqual, 0)
		c.So(user.PlatformChannelID.V(), ShouldNotEqual, 0)
		c.So(p.ID.V(), ShouldNotEqual, 0)
	})
	Convey("create user by phone form zeus", t, func(c C) {
		var user models.User

		util.PG.Model(&user).Where("phone = '13800138000'").First(&user)
		util.PG.Unscoped().Delete(&user)
		var s Session
		s.Phone = optional.OfString("13800138000")
		s.Header = map[string]string{"Platform": "android"}

		user, errResponse, hasError := s.Authenticate()
		if hasError {
			fmt.Printf("errResponse: %v \n", errResponse)
		}
		var p presenter.UserPresenter
		p.Present(&user)
		j, _ := json.Marshal(p)
		fmt.Printf("p: %+v \n", string(j))
		c.So(hasError, ShouldBeFalse)
		c.So(user.Account().ID.V(), ShouldNotEqual, 0)
		c.So(user.PlatformChannelID.V(), ShouldNotEqual, 0)
		c.So(p.ID.V(), ShouldNotEqual, 0)
	})
	Convey("create user by qq and update phone", t, func(c C) {
		var user2 models.User
		if err := util.PG.Model(&user2).Where("phone = '13800138000'").First(&user2).Error; err == nil {
			if err := util.PG.Unscoped().Delete(&user2).Error; err != nil {
				panic(err)
			}
		} else if err != gorm.ErrRecordNotFound {
			panic(err)
		}
		var s Session
		uuid1 := uuid.NewV4()
		s.QQ.UUID = optional.OfString(uuid1.String())
		s.QQ.Avatar = optional.OfString("touxiang")

		// var user models.User
		// util.PG.Where("phone = ?", smsCode.Phone).Delete(&user)
		user, errResponse, hasError := s.Authenticate()
		if hasError {
			fmt.Printf("errResponse: %v \n", errResponse)
		}

		c.So(hasError, ShouldBeFalse)
		c.So(user.Avatar.V(), ShouldEqual, "touxiang")
		var profile Profile
		profile.Phone = optional.OfString("13800138000")
		profile.Avatar = optional.OfString("newtouxiang")
		profile.Wechat = optional.OfString("testwechat")
		profile.QQ = optional.OfString("testqq")
		profile.PersonalWords = optional.OfString("PersonalWords")
		profile.Description = optional.OfString("Description")
		pass := "123456"
		profile.Password = optional.OfString(pass)
		errResponse, hasError = profile.UpdateUser(&user)
		if hasError {
			fmt.Printf("errResponse: %+v", errResponse)
		}
		c.So(hasError, ShouldBeFalse)
		if err := util.PG.Model(&user).Find(&user, user.ID.V()).Error; err == nil {
		} else if err != gorm.ErrRecordNotFound {
			panic(err)
		}
		c.So(user.Phone.V(), ShouldEqual, profile.Phone.V())
		c.So(user.Avatar.V(), ShouldEqual, profile.Avatar.V())
		c.So(user.QQ.V(), ShouldEqual, profile.QQ.V())
		c.So(user.Wechat.V(), ShouldEqual, profile.Wechat.V())
		c.So(user.PersonalWords.V(), ShouldEqual, profile.PersonalWords.V())
		c.So(user.Description.V(), ShouldEqual, profile.Description.V())
		c.So(user.PasswordDigest.V(), ShouldNotBeBlank)
	})
}
