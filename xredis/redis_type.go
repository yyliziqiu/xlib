package xredis

import (
	"errors"
	"time"
)

const (
	DefaultID = "default"

	ModeSingle          = "single"
	ModeCluster         = "cluster"
	ModeSentinel        = "sentinel"
	ModeSentinelCluster = "sentinel-cluster"
)

var (
	ErrNotSupportMode = errors.New("not support redis mode")
)

type Config struct {
	ID   string // optional
	Mode string // must

	// 单机模式
	Addr string // must
	DB   int    // optional

	// 集群模式
	Addrs          []string // must
	ReadPreference string   // must

	// 哨兵模式
	MasterName       string   // must
	SentinelAddrs    []string // must
	SentinelPassword string   // optional
	// DB int // optional
	// ReadPreference string // optional

	Username           string        // optional
	Password           string        // optional
	MaxRetries         int           // optional
	DialTimeout        time.Duration // optional
	ReadTimeout        time.Duration // optional
	WriteTimeout       time.Duration // optional
	PoolSize           int           // optional
	MinIdleConns       int           // optional
	MaxConnAge         time.Duration // optional
	PoolTimeout        time.Duration // optional
	IdleTimeout        time.Duration // optional
	IdleCheckFrequency time.Duration // optional
}

func (c *Config) Default() {
	if c.ID == "" {
		c.ID = DefaultID
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
}
