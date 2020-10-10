package controller

import (
	"gitee.com/dalezhang/account_center/controller/api/v1/address"
	"gitee.com/dalezhang/account_center/controller/api/v1/users"

	"github.com/labstack/echo"
)

// InitRouter 配置路由
func InitRouter(e *echo.Echo) {
	e.POST("/api/users", users.AuthUser)
	e.PUT("/api/users/:id", users.UpdateUser)
	e.GET("/api/users/:id", users.ShowUser)
	e.POST("/api/address", address.CreateAddress)
	e.PUT("/api/address/:id", address.UpdateAddress)
	e.GET("/api/address/:id", address.ShowAddress)
	e.GET("/api/address", address.ListAddress)
	e.DELETE("/api/address/:id", address.DeleteAddress)

}
