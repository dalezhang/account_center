package factory

import (
	"fmt"
	"testing"

	. "gitee.com/dalezhang/account_center/logger"
	"gitee.com/dalezhang/account_center/models"
	"gitee.com/dalezhang/account_center/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/imiskolee/optional"
	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {

	util.InitPG(Logger)

}

func TestCreateBetOrder(t *testing.T) {

	user1 := CreateUser()
	user2 := CreateUser()
	user3 := CreateUser()

	Convey("create bet order", t, func(c C) {
		c.So(user1.ID, ShouldNotEqual, 0)
		c.So(user2.ID, ShouldNotEqual, 0)
		c.So(user3.ID, ShouldNotEqual, 0)
		user1.IsBan = optional.OfBool(true)
		if err := util.PG.Update(&user1).Error; err != nil || user1.ID.V() == 0 {
			err = fmt.Errorf("user修改失败：%+v", err)
			panic(err)
		}
		c.So(user1.IsBan, ShouldBeTrue)
		user1.IsBan = optional.OfBool(false)
		if err := util.PG.Update(&user1).Error; err != nil || user1.ID.V() == 0 {
			err = fmt.Errorf("user修改失败：%+v", err)
			panic(err)
		}
		c.So(user1.IsBan, ShouldBeFalse)
		var user4 models.User
		if err := util.PG.Find(&user4, user1.ID.V()).Error; err != nil || user4.ID.V() == 0 {
			err = fmt.Errorf("user修改失败：%+v", err)
			panic(err)
		}
		c.So(user4.IsBan, ShouldBeFalse)
	})

}
