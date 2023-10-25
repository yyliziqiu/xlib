package xoss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const DefaultId = "default"

var (
	cfs     map[string]Config
	clients map[string]*oss.Client
)

type Config struct {
	Id        string
	Endpoint  string
	KeyId     string
	KeySecret string
	Buckets   []BucketConfig
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

func Init(configs ...Config) error {
	cfs = make(map[string]Config, len(configs))
	for _, config := range configs {
		config = config.WithDefault()
		cfs[config.Id] = config
	}

	clients = make(map[string]*oss.Client, len(cfs))
	for _, cf := range cfs {
		db, err := New(cf)
		if err != nil {
			return err
		}
		clients[cf.Id] = db
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

func GetBucketConfig(clientId string, bucketId string) BucketConfig {
	def := BucketConfig{}

	config, ok := cfs[clientId]
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
