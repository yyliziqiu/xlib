package xapp

import (
	"path/filepath"
	"time"

	"github.com/yyliziqiu/xlib/xdb"
	"github.com/yyliziqiu/xlib/xelastic"
	"github.com/yyliziqiu/xlib/xenv"
	"github.com/yyliziqiu/xlib/xkafka"
	"github.com/yyliziqiu/xlib/xlog"
	"github.com/yyliziqiu/xlib/xredis"
	"github.com/yyliziqiu/xlib/xtask"
	"github.com/yyliziqiu/xlib/xweb"
)

// ICheck 检查配置是否正确
type ICheck interface {
	Check() error
}

// IDefault 为配置项设置默认值
type IDefault interface {
	Default()
}

// IGetLog 为配置项设置默认值
type IGetLog interface {
	GetLog() xlog.Config
}

type Config struct {
	Env      string
	AppId    string
	InsId    string
	BasePath string
	DataPath string
	WaitTime time.Duration

	Log xlog.Config
	Web xweb.Config

	DB      []xdb.Config
	Redis   []xredis.Config
	Kafka   []xkafka.Config
	Elastic []xelastic.Config

	Migration struct {
		EnableTables  bool
		EnableRecords bool
	}

	CronTask []xtask.CronTask
	OnceTask []xtask.OnceTask

	Values map[string]string
}

func (c *Config) Check() error {
	return nil
}

func (c *Config) Default() {
	if c.Env == "" {
		c.Env = xenv.Prod
	}
	if c.AppId == "" {
		c.AppId = "app"
	}
	if c.InsId == "" {
		c.InsId = "1"
	}
	if c.BasePath == "" {
		c.BasePath = "."
	}
	if c.DataPath == "" {
		c.DataPath = filepath.Join(c.BasePath, "data")
	}
}

func (c *Config) GetLog() xlog.Config {
	return c.Log
}
