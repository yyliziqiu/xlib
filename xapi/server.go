package xapi

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/xlib/xapi/xresponse"
	"github.com/yyliziqiu/xlib/xlog"
)

type Config struct {
	Id          string `json:"id"`
	Addr        string `json:"addr"`
	LogName     string `json:"log_name"`
	LogDisabled bool   `json:"log_disabled"`
}

func (c Config) WithDefault() Config {
	if c.Id == "" {
		c.Id = "api"
	}
	if c.LogName == "" {
		c.LogName = "http-"
	}
	return c
}

type SetRoutes func(engine *gin.Engine)

func RunServer(config Config, routes ...SetRoutes) error {
	config = config.WithDefault()

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	setLogger(config)

	engine := createEngine()

	for _, v := range routes {
		v(engine)
	}

	return engine.Run(config.Addr)
}

var logger *logrus.Logger

func setLogger(config Config) {
	if !config.LogDisabled && logger == nil {
		logger = xlog.NewLoggerMust(config.LogName)
	}

	var writer io.Writer
	if logger == nil {
		writer = logger.Writer()
	} else {
		writer = io.Discard
	}
	gin.DefaultWriter = writer
	gin.DefaultErrorWriter = writer
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
	logger.Errorf("Panic, path: %s, error: %v", ctx.FullPath(), err)
	xresponse.InternalServerError(ctx)
}

func GetLogger() *logrus.Logger {
	return logger
}
