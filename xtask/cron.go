package xtask

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/yyliziqiu/xlib/xlog"
)

type CronTask struct {
	Name string
	Spec string
	Cmd  func()
}

func RunCronTasks(ctx context.Context, loc *time.Location, tasksFunc func() []CronTask) {
	cronRunner := cron.New(
		cron.WithSeconds(),
		cron.WithLocation(location(loc)),
	)

	for _, task := range tasksFunc() {
		if task.Spec == "" {
			continue
		}
		_, err := cronRunner.AddFunc(task.Spec, task.Cmd)
		if err != nil {
			xlog.Errorf("Add cron task failed, error: %v.", err)
			return
		}
		xlog.Infof("Add cron task: %s.", task.Name)
	}

	cronRunner.Start()
	xlog.Info("Cron task started.")
	<-ctx.Done()
	cronRunner.Stop()
	xlog.Info("Cron task exit.")
}

func location(loc *time.Location) *time.Location {
	if loc != nil {
		return loc
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		xlog.Errorf("Load locatioin failed, error: %v.", err)
		return time.UTC
	}
	return loc
}

func RunCronTasksWithConfig(ctx context.Context, loc *time.Location, tasksFunc func() []CronTask, configs []CronTask) {
	index := make(map[string]CronTask, len(configs))
	for _, config := range configs {
		index[config.Name] = config
	}

	tasks := tasksFunc()
	for i := 0; i < len(tasks); i++ {
		if config, ok := index[tasks[i].Name]; ok {
			tasks[i].Spec = config.Spec
		}
	}

	RunCronTasks(ctx, loc, func() []CronTask { return tasks })
}

func RunCronTasksByConfig(ctx context.Context, loc *time.Location, tasksFunc func() []CronTask, configs []CronTask) error {
	index := make(map[string]CronTask)
	for _, task := range tasksFunc() {
		index[task.Name] = task
	}

	runnable := make([]CronTask, 0)
	for _, config := range configs {
		task, ok := index[config.Name]
		if !ok {
			return fmt.Errorf("not found cron task[%s]", config.Name)
		}
		task.Spec = config.Spec
		runnable = append(runnable, task)
	}

	RunCronTasks(ctx, loc, func() []CronTask { return runnable })

	return nil
}
