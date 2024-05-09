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

func NewBaseModule(config any) BaseModule {
	return BaseModule{Config: config}
}

func (m BaseModule) Init() (err error) {
	dbi, ok := GetFieldValue(m.Config, "DB")
	if ok {
		dbc, ok2 := dbi.([]xdb.Config)
		if ok2 && len(dbc) > 0 {
			xlog.Info("Init DB.")
			err = xdb.Init(dbc...)
			if err != nil {
				return fmt.Errorf("init DB error [%v]", err)
			}
		}
	}

	redisi, ok := GetFieldValue(m.Config, "Redis")
	if ok {
		redisc, ok2 := redisi.([]xredis.Config)
		if ok2 && len(redisc) > 0 {
			xlog.Info("Init redis.")
			err = xredis.Init(redisc...)
			if err != nil {
				return fmt.Errorf("init redis error [%v]", err)
			}
		}
	}

	kafkai, ok := GetFieldValue(m.Config, "Kafka")
	if ok {
		kafkac, ok2 := kafkai.([]xkafka.Config)
		if ok2 && len(kafkac) > 0 {
			xlog.Info("Init kafka.")
			err = xkafka.Init(kafkac...)
			if err != nil {
				return fmt.Errorf("init kafka error [%v]", err)
			}
		}
	}

	elastici, ok := GetFieldValue(m.Config, "Elastic")
	if ok {
		elasticc, ok2 := elastici.([]xelastic.Config)
		if ok2 && len(elasticc) > 0 {
			xlog.Info("Init elastic.")
			err = xelastic.Init(elasticc...)
			if err != nil {
				return fmt.Errorf("init elastic error [%v]", err)
			}
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
