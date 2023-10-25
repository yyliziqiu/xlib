package xelastic

import (
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/xlib/xlog"
)

var (
	globalLogger *logrus.Logger

	cfs     map[string]Config
	clients map[string]*elastic.Client
)

func SetGlobalLogger(logger *logrus.Logger) {
	globalLogger = logger
}

func Init(configs ...Config) error {
	cfs = make(map[string]Config, len(configs))
	for _, config := range configs {
		config = config.WithDefault()
		cfs[config.Id] = config
	}

	clients = make(map[string]*elastic.Client, len(cfs))
	for _, cf := range cfs {
		client, err := NewClient(cf)
		if err != nil {
			Finally()
			return err
		}
		clients[cf.Id] = client
	}

	return nil
}

func NewClient(config Config) (*elastic.Client, error) {
	var logger elastic.Logger
	if config.EnableLog {
		if config.LogName != "" {
			logger = xlog.MustNewLoggerByName(config.LogName)
		} else {
			if globalLogger == nil {
				globalLogger = xlog.MustNewLoggerByName("es")
			}
			logger = globalLogger
		}
	}

	var traceLogger elastic.Logger
	if config.EnableLogTrace {
		traceLogger = logger
	}

	return elastic.NewClient(
		elastic.SetURL(config.Hosts...),
		elastic.SetBasicAuth(config.Username, config.Password),
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
