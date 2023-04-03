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

type App struct {
	// app 名称
	Name string

	// app 版本
	Version string

	// 配置文件路径
	ConfigPath string

	// 日志目录路径。如果不为空，将覆盖配置文件中的日志路径
	LogPath string

	// 应用关闭等待时间
	ExitDuration time.Duration

	// 全局配置
	Config Config

	// 应用模块
	Modules func() []Module
}

func (app *App) Exec() (err error) {
	err = app.Init()
	if err != nil {
		return err
	}
	return app.Boot()
}

func (app *App) Init() (err error) {
	err = xconfig.Init(app.ConfigPath, app.Config)
	if err != nil {
		return fmt.Errorf("init config error [%v]", err)
	}

	err = app.Config.Check()
	if err != nil {
		return err
	}

	app.Config.Default()

	logC := app.Config.GetLog()
	if app.LogPath != "" {
		logC.Path = app.LogPath
	}
	err = xlog.Init(logC)
	if err != nil {
		return fmt.Errorf("init log error [%v]", err)
	}

	return nil
}

func (app *App) Boot() (err error) {
	modules := app.Modules()
	if len(modules) > 0 {
		RegisterModule(modules...)
	}

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

	if app.ExitDuration > 0 {
		time.Sleep(app.ExitDuration)
	}

	xlog.Info("App exit.")

	return nil
}
