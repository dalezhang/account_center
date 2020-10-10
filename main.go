//go:generate swagger generate spec
package main

import (
	"gitee.com/dalezhang/account_center/controller"
	log "gitee.com/dalezhang/account_center/logger"
	"gitee.com/dalezhang/account_center/models"
	"gitee.com/dalezhang/account_center/util"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Debug = util.Config.Debug

	e.Use(middleware.Logger())  // 日志
	e.Use(middleware.Recover()) // 异常恢复

	// 跨域处理
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},                                        // 允许所有来源
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE}, // 允许跨域请求的方法
	}))

	controller.InitRouter(e)

	util.InitPG(log.Logger)
	defer util.ClosePG()
	// util.InitMysql(Logger)
	// defer util.CloseMysql()
	// util.InitRedis()
	// defer util.CloseRedis()

	// 同步数据结构 可选
	err := util.PG.AutoMigrate(&models.User{}, &models.Account{}, &models.Address{}, &models.PlatformChannel{}).Error
	if err != nil {
		log.Logger.Error("同步数据结构出现异常!", err)
		panic(err)
	}
	log.Logger.Info("test logger")
	e.Logger.Fatal(e.Start(":8080"))
}
