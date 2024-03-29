package xredis

import (
	"github.com/go-redis/redis/v8"
)

var (
	configs  map[string]Config
	clients  map[string]*redis.Client
	clusters map[string]*redis.ClusterClient
)

func Init(cfs ...Config) error {
	configs = make(map[string]Config, 16)
	for _, config := range cfs {
		config.Default()
		configs[config.Id] = config
	}

	clients = make(map[string]*redis.Client, 16)
	clusters = make(map[string]*redis.ClusterClient, 16)
	for _, config := range configs {
		cli, clu := New(config)
		if cli != nil {
			clients[config.Id] = cli
		}
		if clu != nil {
			clusters[config.Id] = clu
		}
	}

	return nil
}

func New(config Config) (*redis.Client, *redis.ClusterClient) {
	switch config.Mode {
	case ModeSentinel:
		return NewFailoverClient(config), nil
	case ModeCluster:
		return nil, NewClusterClient(config)
	case ModeSentinelCluster:
		return nil, NewFailoverClusterClient(config)
	default:
		return NewSingleClient(config), nil
	}
}

func NewSingleClient(config Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:               config.Addr,
		Username:           config.Username,
		Password:           config.Password,
		DB:                 config.DB,
		MaxRetries:         config.MaxRetries,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		MaxConnAge:         config.MaxConnAge,
		PoolTimeout:        config.PoolTimeout,
		IdleTimeout:        config.IdleTimeout,
		IdleCheckFrequency: config.IdleCheckFrequency,
	})
}

func NewClusterClient(config Config) *redis.ClusterClient {
	ops := &redis.ClusterOptions{
		Addrs:              config.Addrs,
		Username:           config.Username,
		Password:           config.Password,
		MaxRetries:         config.MaxRetries,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		MaxConnAge:         config.MaxConnAge,
		PoolTimeout:        config.PoolTimeout,
		IdleTimeout:        config.IdleTimeout,
		IdleCheckFrequency: config.IdleCheckFrequency,
	}

	switch config.ReadPreference {
	case "ReadOnly":
		ops.ReadOnly = true
	case "RouteByLatency":
		ops.RouteByLatency = true
	case "RouteRandomly":
		ops.RouteRandomly = true
	}

	return redis.NewClusterClient(ops)
}

func NewFailoverClient(config Config) *redis.Client {
	ops := &redis.FailoverOptions{
		MasterName:         config.MasterName,
		SentinelAddrs:      config.SentinelAddrs,
		SentinelPassword:   config.SentinelPassword,
		Username:           config.Username,
		Password:           config.Password,
		DB:                 config.DB,
		MaxRetries:         config.MaxRetries,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		MaxConnAge:         config.MaxConnAge,
		PoolTimeout:        config.PoolTimeout,
		IdleTimeout:        config.IdleTimeout,
		IdleCheckFrequency: config.IdleCheckFrequency,
	}

	switch config.ReadPreference {
	case "SlaveOnly":
		ops.SlaveOnly = true
	}

	return redis.NewFailoverClient(ops)
}

func NewFailoverClusterClient(config Config) *redis.ClusterClient {
	ops := &redis.FailoverOptions{
		MasterName:         config.MasterName,
		SentinelAddrs:      config.SentinelAddrs,
		SentinelPassword:   config.SentinelPassword,
		Username:           config.Username,
		Password:           config.Password,
		DB:                 config.DB,
		MaxRetries:         config.MaxRetries,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		MaxConnAge:         config.MaxConnAge,
		PoolTimeout:        config.PoolTimeout,
		IdleTimeout:        config.IdleTimeout,
		IdleCheckFrequency: config.IdleCheckFrequency,
	}

	switch config.ReadPreference {
	case "SlaveOnly":
		ops.SlaveOnly = true
	case "RouteByLatency":
		ops.RouteByLatency = true
	case "RouteRandomly":
		ops.RouteRandomly = true
	}

	return redis.NewFailoverClusterClient(ops)
}

func Finally() {
	for _, client := range clients {
		_ = client.Close()
	}
	for _, cluster := range clusters {
		_ = cluster.Close()
	}
}

func GetCli(id string) *redis.Client {
	return clients[id]
}

func GetDefaultCli() *redis.Client {
	return GetCli(DefaultId)
}

func GetClu(id string) *redis.ClusterClient {
	return clusters[id]
}

func GetDefaultClu() *redis.ClusterClient {
	return GetClu(DefaultId)
}
