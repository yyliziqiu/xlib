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

type App struct {
	// app 名称
	Name string

	// app 版本
	Version string

	// 配置文件路径
	ConfigPath string

	// 日志目录路径
	LogPath string

	// 全局配置
	Config any

	// 模块
	Modules     func() []Module
	ModuleWraps func() []ModuleWrap
}

// Init app
func (app *App) Init() (err error) {
	err = app.InitConfig()
	if err != nil {
		return err
	}
	return app.InitModules()
}

func (app *App) InitConfig() (err error) {
	// 加载配置文件
	err = xconfig.Init(app.ConfigPath, app.Config)
	if err != nil {
		return fmt.Errorf("init config error [%v]", err)
	}

	// 检查配置是否正确
	icheck, ok := app.Config.(ICheck)
	if ok {
		err = icheck.Check()
		if err != nil {
			return err
		}
	}

	// 为配置项设置默认值
	idefault, ok := app.Config.(IDefault)
	if ok {
		idefault.Default()
	}

	// 初始化日志
	logc := xlog.Config{}
	logv, ok := GetFieldValue(app.Config, "Log")
	if ok {
		logc = logv.(xlog.Config)
	}
	if app.LogPath != "" {
		logc.Path = app.LogPath
	}
	err = xlog.Init(logc)
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
	for _, module := range app.Modules() {
		RegisterModule(module)
	}
	for _, wrap := range app.ModuleWraps() {
		RegisterModuleWrap(wrap)
	}
}

// Start app
func (app *App) Start() (err error, f context.CancelFunc) {
	err = app.InitConfig()
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

// Run app
func (app *App) Run() (err error) {
	err = app.InitConfig()
	if err != nil {
		return err
	}
	return app.ExecModulesBlocked()
}

func (app *App) ExecModulesBlocked() (err error) {
	err, cancel := app.ExecModules()
	if err != nil {
		return err
	}

	xlog.Info("App run successfully.")

	exitCh := make(chan os.Signal)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-exitCh

	xlog.Info("App prepare exit.")

	cancel()

	waittimev, ok := GetFieldValue(app.Config, "WaitTime")
	if ok {
		time.Sleep(waittimev.(time.Duration))
	}

	xlog.Info("App exit.")

	return nil
}
