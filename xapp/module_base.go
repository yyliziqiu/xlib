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
	Config any
}

func (m BaseModule) Init() (err error) {
	dbc, ok := GetFieldValue(m.Config, "DB")
	if ok {
		xlog.Info("Init DB.")
		err = xdb.Init(dbc.([]xdb.Config)...)
		if err != nil {
			return fmt.Errorf("init DB error [%v]", err)
		}
	}

	redisc, ok := GetFieldValue(m.Config, "Redis")
	if ok {
		xlog.Info("Init redis.")
		err = xredis.Init(redisc.([]xredis.Config)...)
		if err != nil {
			return fmt.Errorf("init redis error [%v]", err)
		}
	}

	kafkac, ok := GetFieldValue(m.Config, "Kafka")
	if ok {
		xlog.Info("Init kafka.")
		err = xkafka.Init(kafkac.([]xkafka.Config)...)
		if err != nil {
			return fmt.Errorf("init kafka error [%v]", err)
		}
	}

	elasticc, ok := GetFieldValue(m.Config, "Elastic")
	if ok {
		xlog.Info("Init elastic.")
		err = xelastic.Init(elasticc.([]xelastic.Config)...)
		if err != nil {
			return fmt.Errorf("init elastic error [%v]", err)
		}
	}

	return nil
}

func (m BaseModule) Boot(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		xdb.Finally()
		xredis.Finally()
		xkafka.Finally()
		xelastic.Finally()
	}()
	return nil
}
