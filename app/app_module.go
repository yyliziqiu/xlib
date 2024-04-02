package app

import (
	"context"
	"fmt"

	"github.com/yyliziqiu/xlib/xlog"
)

type Module interface {
	Name() string
	Init() error
	Boot(context.Context) error
	Exit() error
}

var _modules []Module

func RegisterModule(modules ...Module) {
	if len(modules) == 0 {
		return
	}
	_modules = append(_modules, modules...)
}

func ExecModules(ctx context.Context, modules ...Module) (err error) {
	err = InitModules(modules...)
	if err != nil {
		return err
	}

	err = BootModules(ctx)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		ExitModules()
	}()

	return nil
}

func InitModules(modules ...Module) error {
	RegisterModule(modules...)
	for _, module := range _modules {
		err := module.Init()
		if err != nil {
			return fmt.Errorf("init module[%s] error [%v]", module.Name(), err)
		}
		xlog.Infof("Init module succeed, module: %s", module.Name())
	}
	return nil
}

func BootModules(ctx context.Context, modules ...Module) error {
	RegisterModule(modules...)
	for _, module := range _modules {
		err := module.Boot(ctx)
		if err != nil {
			return fmt.Errorf("boot module[%s] error [%v]", module.Name(), err)
		}
		xlog.Infof("Boot module succeed, module: %s", module.Name())
	}
	return nil
}

func ExitModules() {
	for i := len(_modules); i > 0; i-- {
		module := _modules[i-1]
		err := module.Exit()
		if err != nil {
			xlog.Errorf("Exit module failed, name: %s, error: %v.", module.Name(), err)
			continue
		}
		xlog.Infof("Exit module succeed, module: %s", module.Name())
	}
}
