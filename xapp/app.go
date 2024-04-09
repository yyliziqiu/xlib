package xapp

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

type ModuleWrapper struct {
	Module Module
	IsBoot bool
}

type App struct {
	// app 名称
	Name string

	// app 版本
	Version string

	// 配置文件路径
	ConfigFile string

	// 日志目录路径。如果不为空，将覆盖配置文件中的日志路径
	LogDir string

	// 应用关闭等待毫秒数
	WaitMS time.Duration

	// 应用模块
	Modules     []ModuleWrapper
	ModulesFunc func() []ModuleWrapper

	// 全局配置
	Config Config
}

func (app *App) Init() (err error) {
	err = app.InitConfigAndLogger()
	if err != nil {
		return err
	}
	return app.InitModules()
}

func (app *App) InitConfigAndLogger() (err error) {
	err = xconfig.Init(app.ConfigFile, app.Config)
	if err != nil {
		return fmt.Errorf("init config error [%v]", err)
	}

	err = app.Config.Check()
	if err != nil {
		return err
	}

	app.Config.Default()

	logC := app.Config.GetLog()
	if app.LogDir != "" {
		logC.Path = app.LogDir
	}
	err = xlog.Init(logC)
	if err != nil {
		return fmt.Errorf("init log error [%v]", err)
	}

	return nil
}

func (app *App) InitModules() (err error) {
	app.registerModules()

	xlog.Info("Prepare init modules.")
	err = InitModules()
	if err != nil {
		xlog.Errorf("Init modules failed, error: %v", err)
		return err
	}

	xlog.Info("Init modules successfully.")

	return nil
}

func (app *App) registerModules() {
	for _, wrapper := range app.Modules {
		RegisterModule(wrapper.Module, false)
	}
	if app.ModulesFunc == nil {
		return
	}
	wrappers := app.ModulesFunc()
	for _, wrapper := range wrappers {
		RegisterModule(wrapper.Module, false)
	}
}

func (app *App) Exec() (err error, f context.CancelFunc) {
	err = app.InitConfigAndLogger()
	if err != nil {
		return err, nil
	}
	return app.ExecModules()
}

func (app *App) ExecModules() (err error, f context.CancelFunc) {
	app.registerModules()

	ctx, cancel := context.WithCancel(context.Background())

	xlog.Info("Prepare exec modules.")
	err = ExecModules(ctx)
	if err != nil {
		xlog.Errorf("Exec modules failed, error: %v", err)
		cancel()
		return err, nil
	}

	xlog.Info("Exec modules successfully.")

	return nil, cancel
}

func (app *App) Exec2() (err error) {
	err = app.InitConfigAndLogger()
	if err != nil {
		return err
	}
	return app.ExecModules2()
}

func (app *App) ExecModules2() (err error) {
	app.registerModules()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	xlog.Info("Prepare exec modules.")
	err = ExecModules(ctx)
	if err != nil {
		xlog.Errorf("Exec modules failed, error: %v", err)
		return err
	}

	xlog.Info("Exec modules successfully.")

	exitCh := make(chan os.Signal)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-exitCh

	xlog.Info("App prepare exit.")

	cancel()

	if app.WaitMS > 0 {
		time.Sleep(app.WaitMS * time.Millisecond)
	}

	xlog.Info("App exit.")

	return nil
}
