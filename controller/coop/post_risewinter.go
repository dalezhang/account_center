package coop

import (
	"net/http"

	"gitee.com/dalezhang/account_center/modules/risewinter"
	"github.com/labstack/echo"
)

func PostRisewinter(context echo.Context) error {
	var c *Controller
	var requestBody []byte
	var data []byte
	var encrypted []byte
	var err error
	var errMessage map[string]string
	c.Prepare(context)
	data = c.Context.Get("data").([]byte)
	encrypted = c.Context.Get("encrypted").([]byte)
	err = risewinter.AnalysisPost(&data, &encrypted)
	if err = c.Context.Bind(&requestBody); err != nil {
		errMessage = map[string]string{"message": "获取请求内容错误"}
		return c.Context.JSON(http.StatusForbidden, errMessage)
	}
	return c.Context.JSON(http.StatusOK, "ok")
}

type params struct {
	Encrypted string `json:"encrypted"`
	Data      string `json:"data"`
}
