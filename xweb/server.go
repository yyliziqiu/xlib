package xweb

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/xlib/xlog"
	"github.com/yyliziqiu/xlib/xweb/xresponse"
)

type Config struct {
	Addr             string
	ErrorLogName     string
	DisableAccessLog bool
	AccessLogName    string
}

func (c *Config) Default() {
	if c.Addr == "" {
		c.Addr = ":80"
	}
	if c.ErrorLogName == "" {
		c.ErrorLogName = "web-error"
	}
	if c.AccessLogName == "" {
		c.AccessLogName = "web-access"
	}
}

func Run(config Config, routes ...func(engine *gin.Engine)) error {
	config.Default()

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	setGinWriter(config)

	engine := createEngine()
	for _, v := range routes {
		v(engine)
	}

	return engine.Run(config.Addr)
}

var (
	errorLogger  *logrus.Logger
	accessLogger *logrus.Logger
)

func setGinWriter(config Config) {
	if errorLogger == nil {
		errorLogger = xlog.NewWithNameMust(config.ErrorLogName)
	}
	if accessLogger == nil && !config.DisableAccessLog {
		accessLogger = xlog.NewWithNameMust(config.AccessLogName)
	}

	gin.DefaultErrorWriter = errorLogger.Writer()
	if accessLogger != nil {
		gin.DefaultWriter = accessLogger.Writer()
	} else {
		gin.DefaultWriter = io.Discard
	}
}

func createEngine() *gin.Engine {
	engine := gin.New()
	engine.NoRoute(xresponse.AbortNotFound)
	engine.NoMethod(xresponse.AbortMethodNotAllowed)
	engine.Use(gin.LoggerWithFormatter(logFormatter))
	engine.Use(gin.CustomRecovery(recovery))
	return engine
}

func logFormatter(param gin.LogFormatterParams) string {
	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("%3d | %13v | %15s |%-7s %#v\n%s",
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Path,
		param.ErrorMessage,
	)
}

func recovery(ctx *gin.Context, err interface{}) {
	errorLogger.Warnf("Panic, path: %s, error: %v", ctx.FullPath(), err)
	xresponse.AbortInternalServerError(ctx)
}

func GetErrorLogger() *logrus.Logger {
	return errorLogger
}

func SetErrorLogger(logger *logrus.Logger) {
	errorLogger = logger
}

func GetAccessLogger() *logrus.Logger {
	return accessLogger
}

func SetAccessLogger(logger *logrus.Logger) {
	accessLogger = logger
}
