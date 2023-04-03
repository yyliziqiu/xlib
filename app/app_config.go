package app

import (
	"path/filepath"

	"github.com/yyliziqiu/xlib/xdb"
	"github.com/yyliziqiu/xlib/xelastic"
	"github.com/yyliziqiu/xlib/xenv"
	"github.com/yyliziqiu/xlib/xkafka"
	"github.com/yyliziqiu/xlib/xlog"
	"github.com/yyliziqiu/xlib/xredis"
	"github.com/yyliziqiu/xlib/xweb"
)

type Config interface {
	// Check 检查配置是否正确
	Check() error

	// Default 为配置项设置默认值
	Default()

	// GetLog 获取日志配置
	GetLog() xlog.Config
}

type BaseConfig struct {
	Env      string
	AppId    string
	SvcId    string
	BasePath string
	DataPath string

	Log     xlog.Config
	Web     xweb.Config
	DB      []xdb.Config
	Redis   []xredis.Config
	Kafka   []xkafka.Config
	Elastic []xelastic.Config

	Migration struct {
		EnableTables  bool
		EnableRecords bool
	}
}

func (c *BaseConfig) Check() error {
	return nil
}

func (c *BaseConfig) Default() {
	if c.Env == "" {
		c.Env = xenv.Prod
	}
	if c.AppId == "" {
		c.AppId = "app"
	}
	if c.SvcId == "" {
		c.SvcId = "1"
	}
	if c.BasePath == "" {
		c.BasePath = "."
	}
	if c.DataPath == "" {
		c.DataPath = filepath.Join(c.BasePath, "data")
	}
}

func (c *BaseConfig) GetLog() xlog.Config {
	return c.Log
}
