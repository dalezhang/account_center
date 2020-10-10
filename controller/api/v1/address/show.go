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

// swagger:operation GET /api/address/{id} addresss ShowAddress
// ShowAddress.
// ---
//
// responses:
//   '200':
//     description: address
//     schema:
//       "$ref": "#/definitions/Address"
func ShowAddress(context echo.Context) error {
	var address models.Address
	var errResponse response.ErrResponse
	addressID := context.Param("id")

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

	return context.JSON(http.StatusOK, address)
}
