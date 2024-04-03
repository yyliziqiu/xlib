package xapp

import (
	"context"
	"fmt"

	"github.com/yyliziqiu/xlib/xdb"
	"github.com/yyliziqiu/xlib/xelastic"
	"github.com/yyliziqiu/xlib/xkafka"
	"github.com/yyliziqiu/xlib/xlog"
	"github.com/yyliziqiu/xlib/xredis"
)

type BaseModule struct {
	Config BaseConfig
}

func (m BaseModule) Name() string {
	return "Base"
}

func (m BaseModule) Init() (err error) {
	c := m.Config
	if len(c.DB) > 0 {
		xlog.Info("Init DB.")
		err = xdb.Init(c.DB...)
		if err != nil {
			return fmt.Errorf("init DB error [%v]", err)
		}
	}

	if len(c.Redis) > 0 {
		xlog.Info("Init redis.")
		err = xredis.Init(c.Redis...)
		if err != nil {
			return fmt.Errorf("init redis error [%v]", err)
		}
	}

	if len(c.Kafka) > 0 {
		xlog.Info("Init kafka.")
		err = xkafka.Init(c.Kafka...)
		if err != nil {
			return fmt.Errorf("init kafka error [%v]", err)
		}
	}

	if len(c.Elastic) > 0 {
		xlog.Info("Init elastic.")
		err = xelastic.Init(c.Elastic...)
		if err != nil {
			return fmt.Errorf("init elastic error [%v]", err)
		}
	}

	return nil
}

func (m BaseModule) Boot(ctx context.Context) error {
	return nil
}

func (m BaseModule) Exit() error {
	xdb.Finally()
	xredis.Finally()
	xkafka.Finally()
	xelastic.Finally()
	return nil
}
