package xelastic

import (
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/xlib/xlog"
	"github.com/yyliziqiu/xlib/xutil"
)

var globalLogger *logrus.Logger

func SetGlobalLogger(logger *logrus.Logger) {
	globalLogger = logger
}

var clients map[string]*elastic.Client

func Init(configs ...Config) error {
	clients = make(map[string]*elastic.Client)
	for _, config := range configs {
		client, err := NewClient(config)
		if err != nil {
			Finally()
			return err
		}
		clients[xutil.IES(config.Id, DefaultId)] = client
	}
	return nil
}

func NewClient(config Config) (*elastic.Client, error) {
	config = config.WithDefault()

	var logger elastic.Logger
	if config.LogEnabled {
		if globalLogger != nil {
			logger = globalLogger
		} else {
			logger = xlog.MustNewLoggerByName(config.LogName)
		}
	}

	var traceLogger elastic.Logger
	if config.LogTrace {
		traceLogger = logger
	}

	return elastic.NewClient(
		elastic.SetURL(config.Hosts...),
		elastic.SetBasicAuth(config.User, config.Password),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetInfoLog(logger),
		elastic.SetErrorLog(logger),
		elastic.SetTraceLog(traceLogger),
	)
}

func Finally() {
	for _, client := range clients {
		client.Stop()
	}
}

func GetClient(id string) *elastic.Client {
	return clients[id]
}

func GetDefaultClient() *elastic.Client {
	return GetClient(DefaultId)
}
