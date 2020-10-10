package address

import (
	"fmt"
	"net/http"

	log "gitee.com/dalezhang/account_center/logger"
	"gitee.com/dalezhang/account_center/models"
	"gitee.com/dalezhang/account_center/modules/response"
	"gitee.com/dalezhang/account_center/util"
	"github.com/imiskolee/optional"
	"github.com/labstack/echo"
)

// swagger:operation POST /api/address addresss CreateAddress
// CreateAddress.
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
func CreateAddress(context echo.Context) error {
	var req addressParam
	var address models.Address
	if err := context.Bind(&req); err != nil {
		log.Logger.Error(fmt.Sprintf("create address err: %v", err))
		errMessage := map[string]string{"message": "请求格式错误"}
		return context.JSON(http.StatusUnauthorized, errMessage)
	}
	errResponse, hasError := req.CreateAddress(&address)
	if hasError {
		return context.JSON(errResponse.Status, errResponse.ErrMessage)
	}
	return context.JSON(http.StatusOK, address)
}

// swagger:model
type addressParam struct {
	UserID optional.Int64 `json:"user_id"`
	// 收件人
	Name optional.String `json:"name"`
	// 邮编
	Postcode optional.String `json:"postcode"`
	// 电话
	Phone optional.String `json:"phone"`
	// 省
	Province optional.String `json:"province"`
	// 市
	City optional.String `json:"city"`
	// 区
	Zone optional.String `json:"zone"`
	// 详细地址
	House optional.String `json:"house"`
	// 标记
	Remark optional.String `json:"remark"`
	// 默认
	IsDefault optional.Bool `json:"is_default"`
}

func (p *addressParam) CreateAddress(address *models.Address) (errResponse response.ErrResponse, hasError bool) {
	var user models.User
	if p.UserID.IsBlank() {
		errResponse.Status = 200
		errResponse.ErrMessage.Code = 8001
		errResponse.ErrMessage.Text = "user_id 不能为空"
		hasError = true
		return
	}
	if err := util.PG.Model(&user).Find(&user, p.UserID.V()).Error; err != nil {
		errResponse.Status = 200
		errResponse.ErrMessage.Code = 8001
		errResponse.ErrMessage.Text = "user_id 无效"
		hasError = true
		return
	}
	address.UserID = p.UserID
	address.Name = p.Name
	address.Postcode = p.Postcode
	address.Phone = p.Phone
	address.Province = p.Province
	address.City = p.City
	address.Zone = p.Zone
	address.House = p.House
	address.Remark = p.Remark
	address.IsDefault = p.IsDefault
	if err := util.PG.Model(&address).Create(&address).Error; err != nil {
		errMsg := fmt.Sprintf("update address err: %v, user: %+v", err, address)
		log.Logger.Error(errMsg)
		errResponse.Status = 200
		errResponse.ErrMessage.Code = 8001
		errResponse.ErrMessage.Text = errMsg
		hasError = true
		return
	}
	return
}

func (p *addressParam) UpdateAddress(address *models.Address) (errResponse response.ErrResponse, hasError bool) {
	address.Name = p.Name
	address.Postcode = p.Postcode
	address.Phone = p.Phone
	address.Province = p.Province
	address.City = p.City
	address.Zone = p.Zone
	address.House = p.House
	address.Remark = p.Remark
	address.IsDefault = p.IsDefault
	if err := util.PG.Model(&address).Updates(&address).Error; err != nil {
		errMsg := fmt.Sprintf("update address err: %v, user: %+v", err, address)
		log.Logger.Error(errMsg)
		errResponse.Status = 200
		errResponse.ErrMessage.Code = 8001
		errResponse.ErrMessage.Text = errMsg
		hasError = true
		return
	}
	return
}
