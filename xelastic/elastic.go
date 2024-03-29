package xelastic

import (
	"sync"

	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/xlib/xlog"
)

var (
	logger     *logrus.Logger
	loggerOnce sync.Once
)

func SetLogger(lg *logrus.Logger) {
	logger = lg
}

func GetLogger() *logrus.Logger {
	if logger != nil {
		return logger
	}
	loggerOnce.Do(func() {
		if logger == nil {
			logger = xlog.NewWithNameMust("elastic")
		}
	})
	return logger
}

var (
	configs map[string]Config
	clients map[string]*elastic.Client
)

func Init(cfs ...Config) error {
	configs = make(map[string]Config, 16)
	for _, config := range cfs {
		config.Default()
		configs[config.Id] = config
	}

	clients = make(map[string]*elastic.Client, 16)
	for _, config := range configs {
		client, err := NewClient(config)
		if err != nil {
			Finally()
			return err
		}
		clients[config.Id] = client
	}

	return nil
}

func NewClient(config Config) (*elastic.Client, error) {
	var lgg elastic.Logger
	if config.EnableLog {
		lgg = GetLogger()
	}

	return elastic.NewClient(
		elastic.SetURL(config.Hosts...),
		elastic.SetBasicAuth(config.Username, config.Password),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetInfoLog(lgg),
		elastic.SetErrorLog(lgg),
		elastic.SetTraceLog(lgg),
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
