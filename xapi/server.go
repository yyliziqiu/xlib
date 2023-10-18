package xapi

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/xlib/xapi/xresponse"
	"github.com/yyliziqiu/xlib/xlog"
)

type Config struct {
	Id               string
	Addr             string
	ErrorLogName     string
	AccessLogName    string
	DisableAccessLog bool
}

func (c Config) WithDefault() Config {
	if c.Id == "" {
		c.Id = "api"
	}
	if c.ErrorLogName == "" {
		c.ErrorLogName = "http-error"
	}
	if c.AccessLogName == "" {
		c.AccessLogName = "http-access"
	}
	return c
}

type RoutesFunc func(engine *gin.Engine)

func Run(config Config, routes ...RoutesFunc) error {
	config = config.WithDefault()

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
		errorLogger = xlog.MustNewLoggerByName(config.ErrorLogName)
	}
	if accessLogger == nil && !config.DisableAccessLog {
		accessLogger = xlog.MustNewLoggerByName(config.AccessLogName)
	}

	gin.DefaultErrorWriter = errorLogger.Writer()
	if accessLogger == nil {
		gin.DefaultWriter = accessLogger.Writer()
	} else {
		gin.DefaultWriter = io.Discard
	}
}

func createEngine() *gin.Engine {
	engine := gin.New()
	engine.NoRoute(xresponse.NotFound)
	engine.NoMethod(xresponse.MethodNotAllowed)
	engine.Use(gin.Logger())
	engine.Use(gin.CustomRecovery(recovery))
	return engine
}

func recovery(ctx *gin.Context, err interface{}) {
	errorLogger.Errorf("Panic, path: %s, error: %v", ctx.FullPath(), err)
	xresponse.InternalServerError(ctx)
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
