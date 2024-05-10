package xboot

import (
	"context"
	"fmt"

	"github.com/yyliziqiu/xlib/xdb"
	"github.com/yyliziqiu/xlib/xelastic"
	"github.com/yyliziqiu/xlib/xkafka"
	"github.com/yyliziqiu/xlib/xlog"
	"github.com/yyliziqiu/xlib/xredis"
)

func BaseInit(config any) InitFunc {
	return func() (err error) {
		val, ok := GetFieldValue(config, "DB")
		if ok {
			c, ok2 := val.([]xdb.Config)
			if ok2 && len(c) > 0 {
				xlog.Info("Init DB.")
				err = xdb.Init(c...)
				if err != nil {
					return fmt.Errorf("init DB error [%v]", err)
				}
			}
		}

		val, ok = GetFieldValue(config, "Redis")
		if ok {
			c, ok2 := val.([]xredis.Config)
			if ok2 && len(c) > 0 {
				xlog.Info("Init redis.")
				err = xredis.Init(c...)
				if err != nil {
					return fmt.Errorf("init redis error [%v]", err)
				}
			}
		}

		val, ok = GetFieldValue(config, "Kafka")
		if ok {
			c, ok2 := val.([]xkafka.Config)
			if ok2 && len(c) > 0 {
				xlog.Info("Init kafka.")
				err = xkafka.Init(c...)
				if err != nil {
					return fmt.Errorf("init kafka error [%v]", err)
				}
			}
		}

		val, ok = GetFieldValue(config, "Elastic")
		if ok {
			c, ok2 := val.([]xelastic.Config)
			if ok2 && len(c) > 0 {
				xlog.Info("Init elastic.")
				err = xelastic.Init(c...)
				if err != nil {
					return fmt.Errorf("init elastic error [%v]", err)
				}
			}
		}

		return nil
	}
}

func BaseBoot() BootFunc {
	return func(ctx context.Context) error {
		go func() {
			<-ctx.Done()
			xdb.Finally()
			xredis.Finally()
			xkafka.Finally()
			xelastic.Finally()
		}()
		return nil
	}
}
