package users

import (
	"fmt"
	"net/http"
	"strconv"

	. "gitee.com/dalezhang/account_center/logger"
	"gitee.com/dalezhang/account_center/models"
	"gitee.com/dalezhang/account_center/modules/response"
	"gitee.com/dalezhang/account_center/modules/session"
	"gitee.com/dalezhang/account_center/presenter"
	"gitee.com/dalezhang/account_center/util"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// swagger:operation PUT /api/users/{id} users UpdateUser
// UpdateUser.
// ---
//
// parameters:
// - name: session
//   in: body
//   description: the full description of the UpdateUserRequest
//   required: true
//   schema:
//     "$ref": "#/definitions/Profile"
// responses:
//   '200':
//     description: user
//     schema:
//       "$ref": "#/definitions/UserPresenter"
func UpdateUser(context echo.Context) error {
	var user models.User
	var errResponse response.ErrResponse
	var req session.Profile
	userID, _ := strconv.ParseInt(context.Param("id"), 10, 64)
	if err := util.PG.Find(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			//   super code: 9003, text: '用户不存在', status: 200
			errResponse.Status = 200
			errResponse.ErrMessage.Code = 9003
			errResponse.ErrMessage.Text = "用户不存在"
			return context.JSON(errResponse.Status, errResponse.ErrMessage)
		}
		errString := fmt.Sprintf("find user err: %v", err)
		Logger.Error(errString)
		errResponse.Status = 200
		errResponse.ErrMessage.Code = 9003
		errResponse.ErrMessage.Text = errString
		return context.JSON(errResponse.Status, errResponse.ErrMessage)
	}
	if err := context.Bind(&req); err != nil {
		Logger.Error(fmt.Sprintf("Bind req err: %v", err))
		errMessage := map[string]string{"message": "请求格式错误"}
		return context.JSON(http.StatusUnauthorized, errMessage)
	}
	errResponse, hasError := req.UpdateUser(&user)
	if hasError {
		return context.JSON(errResponse.Status, errResponse.ErrMessage)
	}
	var p presenter.UserPresenter
	p.Present(&user)
	// j, _ := json.Marshal(p)
	return context.JSON(http.StatusOK, p)
}
