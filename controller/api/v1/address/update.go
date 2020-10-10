package address

import (
	"fmt"
	"net/http"

	log "gitee.com/dalezhang/account_center/logger"
	"gitee.com/dalezhang/account_center/models"
	"gitee.com/dalezhang/account_center/modules/response"
	"gitee.com/dalezhang/account_center/util"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// swagger:operation PUT /api/address/{id} addresss UpdateAddress
// UpdateAddress.
// ---
//
// parameters:
// - name: address
//   in: body
//   description: the full description of the UpdateUserRequest
//   required: true
//   schema:
//     "$ref": "#/definitions/addressParam"
// responses:
//   '200':
//     description: address
//     schema:
//       "$ref": "#/definitions/Address"
func UpdateAddress(context echo.Context) error {
	var req addressParam
	var address models.Address
	var errResponse response.ErrResponse
	addressID := context.Param("id")
	if err := context.Bind(&req); err != nil {
		log.Logger.Error(fmt.Sprintf("update address err: %v", err))
		errMessage := map[string]string{"message": "请求格式错误"}
		return context.JSON(http.StatusUnauthorized, errMessage)
	}
	if err := util.PG.Model(&address).Find(&address, addressID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errResponse.Status = 200
			errResponse.ErrMessage.Code = 8011
			errResponse.ErrMessage.Text = "地址不存在"
			return context.JSON(errResponse.Status, errResponse.ErrMessage)
		}
		errString := fmt.Sprintf("find address err: %v", err)
		log.Logger.Error(errString)
		errResponse.Status = 200
		errResponse.ErrMessage.Code = 8011
		errResponse.ErrMessage.Text = errString
		return context.JSON(errResponse.Status, errResponse.ErrMessage)
	}
	errResponse, hasError := req.UpdateAddress(&address)
	if hasError {
		return context.JSON(errResponse.Status, errResponse.ErrMessage)
	}
	return context.JSON(http.StatusOK, address)
}
