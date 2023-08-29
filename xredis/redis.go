package xredis

import (
	"github.com/go-redis/redis/v8"

	"github.com/yyliziqiu/xlib/xutil"
)

var (
	clients  map[string]*redis.Client
	clusters map[string]*redis.ClusterClient
)

func Initialize(configs ...Config) error {
	clients = make(map[string]*redis.Client, len(configs))
	clusters = make(map[string]*redis.ClusterClient, len(configs))
	for _, config := range configs {
		cli, clu, err := New(config)
		if err != nil {
			Finally()
			return err
		}
		id := xutil.IES(config.Id, DefaultId)
		if cli != nil {
			clients[id] = cli
		}
		if clu != nil {
			clusters[id] = clu
		}
	}
	return nil
}

func New(config Config) (*redis.Client, *redis.ClusterClient, error) {
	switch config.Mode {
	case ModeSingle:
		return NewClient(config), nil, nil
	case ModeCluster:
		return nil, NewClusterClient(config), nil
	case ModeSentinel:
		return NewFailoverClient(config), nil, nil
	case ModeSentinelCluster:
		return nil, NewFailoverClusterClient(config), nil
	default:
		return nil, nil, ErrNotSupportMode
	}
}

func NewClient(config Config) *redis.Client {
	config = config.WithDefault()

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
	config = config.WithDefault()

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
	config = config.WithDefault()

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
	config = config.WithDefault()

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
