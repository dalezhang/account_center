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

// swagger:operation GET /api/address addresss ListAddress
// ListAddress.
// ---
//
// parameters:
// - name: user_id
//   in: path
//   description: user id
//   required: true
// responses:
//   '200':
//     description: address
//     schema:
//       "$ref": "#/definitions/AddressList"
func ListAddress(context echo.Context) error {
	var address models.Address
	var addressList AddressList
	var errResponse response.ErrResponse
	userID := context.QueryParam("user_id")
	if userID == "" {
		errResponse.Status = 200
		errResponse.ErrMessage.Code = 8022
		errResponse.ErrMessage.Text = "user_id不能为空"
		return context.JSON(http.StatusUnauthorized, errResponse.ErrMessage)
	}
	if err := util.PG.Model(&address).Where("user_id = ?", userID).Find(&addressList.AddressList).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errResponse.Status = 200
			errResponse.ErrMessage.Code = 8021
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

	return context.JSON(http.StatusOK, addressList)
}

// swagger:model AddressList
type AddressList struct {
	AddressList []models.Address `json:"address_list"`
}
