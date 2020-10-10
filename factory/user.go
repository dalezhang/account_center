package factory

import (
	"fmt"

	"gitee.com/dalezhang/account_center/helpers/secure_random"
	"gitee.com/dalezhang/account_center/models"
	"gitee.com/dalezhang/account_center/util"
	"github.com/imiskolee/optional"
)

func CreateUser() (user models.User) {
	// phone := fmt.Sprintf("138%d", secure_random.RandomNumber(8))
	fmt.Printf("===============phone: %+v \n", secure_random.RandomNumber(8))
	user = models.User{
		Name: optional.OfString(secure_random.NewString(5)),
		// Phone: optional.OfString(phone),
	}
	// errors := user.Validate()
	// if len(errors) > 0 {
	// 	panic(fmt.Errorf("%+v", errors))
	// }

	if err := util.PG.Create(&user).Error; err != nil || user.ID.V() == 0 {
		err = fmt.Errorf("user创建失败：%+v", err)
		panic(err)
	}
	return
}
