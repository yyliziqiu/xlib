package xapp

import (
	"context"
	"fmt"
)

type Module interface {
	Init() error
	Boot(context.Context) error
}

var (
	_modules []Module
	_isBoots []bool
)

func RegisterModule(module Module, isBoot bool) {
	_modules = append(_modules, module)
	_isBoots = append(_isBoots, isBoot)
}

func ExecModules(ctx context.Context) (err error) {
	err = InitModules()
	if err != nil {
		return err
	}

	err = BootModules(ctx)
	if err != nil {
		return err
	}

	return nil
}

func InitModules() error {
	for _, module := range _modules {
		err := module.Init()
		if err != nil {
			return fmt.Errorf("init module error [%v]", err)
		}
	}
	return nil
}

func BootModules(ctx context.Context) error {
	for i, module := range _modules {
		if !_isBoots[i] {
			continue
		}
		err := module.Boot(ctx)
		if err != nil {
			return fmt.Errorf("boot module error [%v]", err)
		}
	}
	return nil
}
