package users

import (
	"fmt"
	"net/http"

	. "gitee.com/dalezhang/account_center/logger"
	"gitee.com/dalezhang/account_center/models"
	"gitee.com/dalezhang/account_center/modules/response"
	"gitee.com/dalezhang/account_center/presenter"
	"gitee.com/dalezhang/account_center/util"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// swagger:operation GET /api/users/{id} users ShowUser
// ShowUser.
// ---
//
// responses:
//   '200':
//     description: current_user
//     schema:
//       "$ref": "#/definitions/UserPresenter"
func ShowUser(context echo.Context) error {
	var errResponse response.ErrResponse
	var user models.User
	userID := context.Param("id")
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
	var p presenter.UserPresenter
	p.Present(&user)
	return context.JSON(http.StatusOK, p)
}
