package logger

import (
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	// log "github.com/sirupsen/logrus"
)

var defaultLog *log.Logger

var (
	Logger  *zap.SugaredLogger
	ZLogger *zap.Logger
)

func InitLogger(logDir, env string) {
	var (
		err error
	)

	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s.log", logDir, env),
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	ZLogger := zap.New(core)
	// if deubg {
	// ZLogger, err = zap.NewDevelopment()
	// } else {
	// ZLogger, err = zap.NewProduction()
	// }
	if err != nil {
		panic(err)
	}

	Logger = ZLogger.Sugar()
}
func formatEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%d%02d%02d_%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
}
