package xredis

import (
	"errors"
	"time"
)

const (
	DefaultId = "default"

	ModeSingle          = "single"
	ModeCluster         = "cluster"
	ModeSentinel        = "sentinel"
	ModeSentinelCluster = "sentinel-cluster"
)

var (
	ErrNotSupportMode = errors.New("not support redis mode")
)

type Config struct {
	Id   string `json:"id"`
	Mode string `json:"mode"`

	// 单机模式
	Addr string `json:"addr"`
	DB   int    `json:"DB"`

	// 集群模式
	Addrs          []string `json:"addrs"`
	ReadPreference string   `json:"readPreference"`

	// 哨兵模式
	MasterName       string   `json:"masterName"`
	SentinelAddrs    []string `json:"sentinelAddrs"`
	SentinelPassword string   `json:"sentinelPassword"`
	// DB int
	// ReadPreference string

	Username           string        `json:"username"`
	Password           string        `json:"password"`
	MaxRetries         int           `json:"maxRetries"`
	DialTimeout        time.Duration `json:"dialTimeout"`
	ReadTimeout        time.Duration `json:"readTimeout"`
	WriteTimeout       time.Duration `json:"writeTimeout"`
	PoolSize           int           `json:"poolSize"`
	MinIdleConns       int           `json:"minIdleConns"`
	MaxConnAge         time.Duration `json:"maxConnAge"`
	PoolTimeout        time.Duration `json:"poolTimeout"`
	IdleTimeout        time.Duration `json:"idleTimeout"`
	IdleCheckFrequency time.Duration `json:"idleCheckFrequency"`
}

func (c Config) WithDefault() Config {
	if c.Id == "" {
		c.Id = DefaultId
	}
	if c.MaxRetries == 0 {
		c.MaxRetries = 3
	}
	if c.DialTimeout == 0 {
		c.DialTimeout = 30 * time.Second
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 10 * time.Second
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = 10 * time.Second
	}
	if c.PoolSize == 0 {
		c.PoolSize = 50
	}
	if c.MinIdleConns == 0 {
		c.MinIdleConns = 10
	}
	if c.MaxConnAge == 0 {
		c.MaxConnAge = time.Hour
	}
	if c.PoolTimeout == 0 {
		c.PoolTimeout = 10 * time.Second
	}
	if c.IdleTimeout == 0 {
		c.IdleTimeout = 10 * time.Minute
	}
	if c.IdleCheckFrequency == 0 {
		c.IdleCheckFrequency = 30 * time.Second
	}
	return c
}
