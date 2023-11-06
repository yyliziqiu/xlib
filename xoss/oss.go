package xoss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const DefaultID = "default"

var (
	configs map[string]Config
	clients map[string]*oss.Client
)

type Config struct {
	ID        string
	Endpoint  string
	KeyID     string
	KeySecret string
}

func (c Config) Default() Config {
	if c.ID == "" {
		c.ID = DefaultID
	}
	return c
}

type BucketConfig struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
}

func Init(cfs ...Config) error {
	configs = make(map[string]Config, len(cfs))
	for _, config := range cfs {
		config.Default()
		configs[config.ID] = config
	}

	clients = make(map[string]*oss.Client, len(configs))
	for _, config := range configs {
		db, err := New(config)
		if err != nil {
			return err
		}
		clients[config.ID] = db
	}

	return nil
}

func New(config Config) (*oss.Client, error) {
	return oss.New(config.Endpoint, config.KeyID, config.KeySecret)
}

func GetClient(id string) *oss.Client {
	return clients[id]
}

func GetDefaultClient() *oss.Client {
	return GetClient(DefaultID)
}
