package factory

import (
	"fmt"

	"gitee.com/dalezhang/account_center/helpers/secure_random"
	"gitee.com/dalezhang/account_center/models"
	"gitee.com/dalezhang/account_center/util"
	"github.com/imiskolee/optional"
)

func CreateSmsCode(phone string) (smsCode models.SmsCode) {
	smsCode = models.SmsCode{
		Code:        optional.OfString(fmt.Sprintln(secure_random.RandomNumber(5))),
		Phone:       optional.OfString(phone),
		DialingCode: optional.OfString("86"),
	}

	if err := util.PG.Create(&smsCode).Error; err != nil || smsCode.ID.V() == 0 {
		err = fmt.Errorf("\n smsCode创建失败：%+v", err)
		panic(err)
	}
	return
}
