package xoss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const DefaultId = "default"

var (
	configs map[string]Config
	clients map[string]*oss.Client
)

type Config struct {
	Id        string
	Endpoint  string
	KeyId     string
	KeySecret string
}

func (c *Config) Default() {
	if c.Id == "" {
		c.Id = DefaultId
	}
}

type BucketConfig struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
}

func Init(cfs ...Config) error {
	configs = make(map[string]Config, 16)
	for _, config := range cfs {
		config.Default()
		configs[config.Id] = config
	}

	clients = make(map[string]*oss.Client, 16)
	for _, config := range configs {
		db, err := New(config)
		if err != nil {
			return err
		}
		clients[config.Id] = db
	}

	return nil
}

func New(config Config) (*oss.Client, error) {
	return oss.New(config.Endpoint, config.KeyId, config.KeySecret)
}

func GetClient(id string) *oss.Client {
	return clients[id]
}

func GetDefaultClient() *oss.Client {
	return GetClient(DefaultId)
}
