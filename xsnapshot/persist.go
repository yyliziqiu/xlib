package xsnapshot

import (
	"context"
	"time"

	"github.com/yyliziqiu/xlib/xlog"
	"github.com/yyliziqiu/xlib/xutil"
)

type Persistence interface {
	Name() string
	Load() error
	Save() error
	Interval() time.Duration
}

func Persist(ctx context.Context, persistencesFunc func() []Persistence) error {
	persistences := persistencesFunc()

	err := load(persistences)
	if err != nil {
		return err
	}

	for _, persistence := range persistences {
		go runSave(ctx, persistence)
	}

	return nil
}

func load(persistences []Persistence) error {
	timer := xutil.NewTimer()
	for _, persistence := range persistences {
		err := persistence.Load()
		if err != nil {
			xlog.Errorf("Load snapshot failed, name: %s, error: %v.", persistence.Name(), err)
			return err
		}
		xlog.Infof("Load snapshot succeed, name: %s, cost: %s.", persistence.Name(), timer.Pauses())
	}
	xlog.Infof("Loaded all snapshots, cost: %s.", timer.Stops())
	return nil
}

func runSave(ctx context.Context, persistence Persistence) {
	interval := persistence.Interval()
	if persistence.Interval() <= 0 {
		interval = 10 * 365 * 24 * time.Hour
	}

	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			_ = save(persistence)
		case <-ctx.Done():
			_ = save(persistence)
			return
		}
	}
}

func save(persistence Persistence) error {
	timer := xutil.NewTimer()
	err := persistence.Save()
	if err != nil {
		xlog.Errorf("Save snapshot failed, name: %s, error: %v.", persistence.Name(), err)
	} else {
		xlog.Infof("Save snapshot succeed, name: %s, cost: %s.", persistence.Name(), timer.Stops())
	}
	return err
}
