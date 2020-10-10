package helpers

import (
	"time"

	"gitee.com/ouresports/hex/logger"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type SuperController struct {
	Body        []byte
	RequestID   string
	RequestTime time.Time
	Log         *logrus.Entry
	code        int
	Context     echo.Context
}

func (c *SuperController) Prepare(contest echo.Context) {
	c.Context = contest
	c.code = 200
	contest.Bind(c.Body)
	uuid, _ := uuid.NewV4()
	c.RequestID = uuid.String()
	c.RequestTime = time.Now()
	c.Log = logger.WithFields(logrus.Fields{
		"request_id":     c.RequestID,
		"request_time":   c.RequestTime,
		"request_uri":    c.Context.Request().RequestURI,
		"request_query":  c.Context.Request().URL.Query(),
		"request_param":  c.Context.QueryParams(),
		"request_method": c.Context.Request().Method,
	})
	header := c.Context.Response().Header()

	header.Set("X-Request-ID", c.RequestID)
}
