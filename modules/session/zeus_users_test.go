package session

import (
	"fmt"
	"strings"
	"testing"

	_ "gitee.com/dalezhang/account_center/util"

	_ "gitee.com/dalezhang/account_center/logger"
	"github.com/imiskolee/optional"
	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadOneUserFromZeus(t *testing.T) {
	fmt.Println("has port", strings.LastIndex("0.0.0.0:3001", ":"))
	// 测试api需要启动zeus
	Convey("find user by phone", t, func(c C) {
		var s Session
		s.Phone = optional.OfString("18700000000")
		err, noFound, zeusUser := s.LoadOneUserFromZeus()
		c.So(err, ShouldBeNil)
		c.So(noFound, ShouldBeFalse)
		c.So(zeusUser.Phone.V(), ShouldEqual, "18700000000")
	})
}
