package util

import (
	"log"
	"os"

	"go.uber.org/zap"

	"github.com/jinzhu/gorm"
)

var (
	// Engine xorm引擎
	PG          *gorm.DB
	PGEngine *gorm.DB
	// Redis redis引擎
	// Redis         *redis.Client
)

type Logger struct {
	*zap.SugaredLogger
}

func InitPG(sl *zap.SugaredLogger) {
	var err error
	if PG, err = gorm.Open("postgres", Config.PGDNS); err != nil {
		log.Fatal("数据库postgres初始化失败!", err)
	}
	if Config.Debug {
		PG.Debug()
		PG.LogMode(true)
		PG.SetLogger(log.New(os.Stdout, "\r\n", 0))
	} else {
		PG.SetLogger(gorm.Logger{
			LogWriter: &Logger{
				SugaredLogger: sl,
			},
		})
	}

}
func (l *Logger) Println(values ...interface{}) {
	l.SugaredLogger.Info(values...)
}

// func InitMysql(sl *zap.SugaredLogger) {
// 	var err error

// 	mysqlURL := os.Getenv("MYSQL_DNS")
// 	fmt.Println("\n MYSQL_DNS============ from env:  ", mysqlURL)

// 	if mysqlURL == "" {
// 		mysqlURL = "root:root@tcp(mysql:33061)/arc-warden_test?parseTime=True&loc=Local"
// 		fmt.Println("\n MYSQL_DNS ============ :  ", mysqlURL)
// 	}

// 	if MysqlEngine, err = gorm.Open("mysql", mysqlURL); err != nil {
// 		log.Fatal("数据库Mysql初始化失败!", err)
// 	}
// 	MysqlEngine.SetLogger(gorm.Logger{
// 		LogWriter: &Logger{
// 			SugaredLogger: sl,
// 		},
// 	})
// }

// func InitRedis() {
// 	Redis = redis.NewClient(&redis.Options{
// 		Addr: fmt.Sprintf("%s:%d",
// 			Config.GetString("redis.host"),
// 			Config.GetInt("redis.port"),
// 		),
// 	})

// 	if _, err := Redis.Ping().Result(); err != nil {
// 		log.Fatal("redis连接失败!", err)
// 	}
// }

func ClosePG() {
	PG.Close()
}

// func CloseMysql() {
// 	MysqlEngine.Close()
// }

// func CloseRedis() {
// 	Redis.Close()
// }
