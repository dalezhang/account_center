// Package classification User API.
//
// The purpose of this service is to provide an application
// that is using plain go code to define an API
//
//      Host: localhost
//      Version: 0.0.1
//
// swagger:meta
package users

import (
	"fmt"
	"net/http"

	. "gitee.com/dalezhang/account_center/logger"
	"gitee.com/dalezhang/account_center/modules/session"
	"gitee.com/dalezhang/account_center/presenter"
	"github.com/labstack/echo"
)

// swagger:operation POST /api/users users AuthUserRequest
// AuthUser.
// ---
//
// parameters:
// - name: Platform
//   in: header
//   description: ios, android
//   required: true
// - name: Channel
//   in: header
//   description: apple, huawei, postman, vivo, test, ios_official, ios_pugongyi, nx1, nx2, nx3, risewinter, undefined, zl2, zl3
//   required: true
// - name: session
//   in: body
//   description: the full description of the session
//   required: true
//   schema:
//     "$ref": "#/definitions/Session"
// responses:
//   '200':
//     description: user
//     schema:
//       "$ref": "#/definitions/UserPresenter"
func AuthUser(context echo.Context) error {
	var req session.Session
	platform := context.Request().Header.Get("Platform")
	channel := context.Request().Header.Get("Channel")

	if platform == "" {
		errMessage := map[string]interface{}{"message": "未知的渠道来源", "error_code": 2099}
		return context.JSON(400, errMessage)
	}
	if err := context.Bind(&req); err != nil {
		Logger.Error(fmt.Sprintf("find user err: %v", err))
		errMessage := map[string]string{"message": "请求格式错误"}
		return context.JSON(http.StatusUnauthorized, errMessage)
	}
	req.Header = map[string]string{"Platform": platform, "Channel": channel}

	user, errMessage, hasError := req.Authenticate()
	if hasError {
		return context.JSON(errMessage.Status, errMessage.ErrMessage)
	}
	var p presenter.UserPresenter
	p.Present(&user)
	// j, _ := json.Marshal(p)
	return context.JSON(http.StatusOK, p)
}

// {
// 	"session": {
// 		"dialing_code": "8å6",
// 		"password": "123456",
// 		"phone": "19876543211"
// 	}
// }

// code = self.find_by(code: code, phone: phone, dialing_code: dialing_code)
// return true if code
// return false
