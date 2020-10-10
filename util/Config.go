package util

import (
	"fmt"
	"os"

	"gitee.com/dalezhang/account_center/logger"
	"github.com/spf13/viper"
)

// Config 配置信息实例

var Config MyConfig

func init() {
	var config = viper.New()
	// GOPATH = os.Getenv("GOPATH")
	config.SetConfigName("config")
	config.AddConfigPath("./conf")
	config.AddConfigPath("../conf")
	config.AddConfigPath("../../conf")

	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	Config.PGDNS = os.Getenv("PG_DNS")
	fmt.Println("\n PG_DNS============ from env:  ", Config.PGDNS)

	if Config.PGDNS == "" {
		Config.PGDNS = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.GetString("pg.host"),
			config.GetInt("pg.port"),
			config.GetString("pg.user"),
			config.GetString("pg.password"),
			config.GetString("pg.dbname"),
		)
	}
	debug := os.Getenv("DEBUG")
	if debug == "" {
		Config.Debug = config.GetBool("server.debug")
	} else {
		Config.Debug = true
	}
	Config.ZeusDNS = os.Getenv("ZEUS_DNS")
	if Config.ZeusDNS == "" {
		Config.ZeusDNS = config.GetString("server.zeus_dns")
	}
	Config.ZeusAdminToken = os.Getenv("ZEUS_ADMIN_TOKEN")
	if Config.ZeusAdminToken == "" {
		Config.ZeusAdminToken = config.GetString("server.zeus_admin_token")
	}
	Config.ENV = os.Getenv("ENV")
	if Config.ENV == "" {
		Config.ENV = config.GetString("server.env")
	}
	Config.LogDir = os.Getenv("LOG_DIR")
	if Config.LogDir == "" {
		Config.LogDir = config.GetString("server.log_dir")
	}
	logger.InitLogger(Config.LogDir, Config.ENV)
}

type MyConfig struct {
	PGDNS          string
	Debug          bool
	ZeusDNS        string
	ZeusAdminToken string
	ENV            string
	LogDir         string
}
