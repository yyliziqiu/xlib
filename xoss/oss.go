package xoss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/yyliziqiu/xlib/xutil"
)

const DefaultId = "default"

var (
	configs map[string]Config
	clients map[string]*oss.Client
)

type Config struct {
	Id        string         `json:"id"`
	Endpoint  string         `json:"endpoint"`
	KeyId     string         `json:"key_id"`
	KeySecret string         `json:"key_secret"`
	Buckets   []BucketConfig `json:"buckets"`
}

func (c Config) WithDefault() Config {
	if c.Id == "" {
		c.Id = DefaultId
	}
	return c
}

type BucketConfig struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
}

func Initialize(configs ...Config) error {
	clients = make(map[string]*oss.Client, len(configs))
	for _, config := range configs {
		db, err := New(config)
		if err != nil {
			return err
		}
		clients[xutil.IES(config.Id, DefaultId)] = db
	}
	return nil
}

func New(config Config) (*oss.Client, error) {
	config = config.WithDefault()
	return oss.New(config.Endpoint, config.KeyId, config.KeySecret)
}

func GetClient(id string) *oss.Client {
	return clients[id]
}

func GetDefaultClient() *oss.Client {
	return GetClient(DefaultId)
}

func GetBucketConfig(clientId string, bucketId string) BucketConfig {
	def := BucketConfig{}

	config, ok := configs[clientId]
	if !ok {
		return def
	}

	for _, bucket := range config.Buckets {
		if bucket.Id == bucketId {
			return bucket
		}
	}

	return def
}

func GetDefaultClientBucketConfig(bucketId string) BucketConfig {
	return GetBucketConfig(DefaultId, bucketId)
}
