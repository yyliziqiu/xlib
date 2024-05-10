package xboot

import (
	"context"
)

type InitFunc func() error

type BootFunc func(context.Context) error

type InitFuncs []InitFunc

func (list InitFuncs) Init() error {
	for _, fun := range list {
		err := fun()
		if err != nil {
			return err
		}
	}
	return nil
}

type BootFuncs []BootFunc

func (list BootFuncs) Boot(ctx context.Context) error {
	for _, fun := range list {
		err := fun(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
