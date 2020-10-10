package api

import (
	"gitee.com/dalezhang/account_center/helpers"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type SuperController struct {
	helpers.SuperController
	UserID int64
}

func (c *SuperController) Prepare(contest echo.Context) {
	c.SuperController.Prepare(contest)
	jwtToken := contest.Get("user").(*jwt.Token) // 默认使用user字符串来获取jwt中间件过滤后解析出来的token
	claims := jwtToken.Claims.(jwt.MapClaims)
	c.UserID = int64(claims["userId"].(float64)) // map转化后的所有数字都会变成float64
}
