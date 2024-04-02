package xelastic

import (
	"github.com/olivere/elastic/v7"
)

var (
	_configs map[string]Config
	_clients map[string]*elastic.Client
)

func Init(configs ...Config) error {
	_configs = make(map[string]Config, 8)
	for _, config := range configs {
		_configs[config.Id] = config.Default()
	}

	_clients = make(map[string]*elastic.Client, 8)
	for _, config := range _configs {
		client, err := NewClient(config)
		if err != nil {
			Finally()
			return err
		}
		_clients[config.Id] = client
	}

	return nil
}

func NewClient(config Config) (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetURL(config.Hosts...),
		elastic.SetBasicAuth(config.Username, config.Password),
		elastic.SetHttpClient(config.Client),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetInfoLog(config.Logger),
		elastic.SetErrorLog(config.Logger),
		elastic.SetTraceLog(config.Logger),
	)
}

func Finally() {
	for _, client := range _clients {
		client.Stop()
	}
}

func GetCli(id string) *elastic.Client {
	return _clients[id]
}

func DefaultCli() *elastic.Client {
	return GetCli(DefaultId)
}
