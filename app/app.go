package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yyliziqiu/xlib/xconfig"
	"github.com/yyliziqiu/xlib/xlog"
)

func Exec2(config Config, cpath string, register func()) (err error) {
	err = InitSystem(config, cpath)
	if err != nil {
		return err
	}
	return Exec(config, register)
}

func InitSystem(config Config, cpath string) (err error) {
	err = xconfig.Init(cpath, config)
	if err != nil {
		return fmt.Errorf("init config error [%v]", err)
	}

	config.Default()

	err = config.Check()
	if err != nil {
		return err
	}

	err = xlog.Init(config.GetLog())
	if err != nil {
		return fmt.Errorf("init log error [%v]", err)
	}

	return nil
}

func Exec(config Config, register func()) (err error) {
	register()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	xlog.Info("Prepare exec modules.")
	err = ExecModules(ctx)
	if err != nil {
		xlog.Errorf("Exec modules failed, error: %v", err)
		return err
	}

	xlog.Info("Exec app successfully.")

	exitCh := make(chan os.Signal)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-exitCh

	xlog.Info("App prepare exit.")

	cancel()

	if config.GetWaitTime() > 0 {
		time.Sleep(config.GetWaitTime())
	}

	xlog.Info("App exit.")

	return nil
}
