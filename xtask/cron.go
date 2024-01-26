package xtask

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
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
			Logger.Errorf("Add cron task failed, error: %v.", err)
			return
		}
		Logger.Infof("Add cron task: %s.", task.Name)
	}

	cronRunner.Start()
	Logger.Info("Cron task started.")
	<-ctx.Done()
	cronRunner.Stop()
	Logger.Info("Cron task exit.")
}

func location(loc *time.Location) *time.Location {
	if loc != nil {
		return loc
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		Logger.Errorf("Load locatioin failed, error: %v.", err)
		return time.UTC
	}
	return loc
}

func RunCronTasksWithConfig(ctx context.Context, loc *time.Location, tasksFunc func() []CronTask, configs []CronTask) {
	tasks := tasksFunc()

	ConfigCronTasks(tasks, configs)

	RunCronTasks(ctx, loc, func() []CronTask { return tasks })
}

func ConfigCronTasks(tasks []CronTask, configs []CronTask) {
	index := make(map[string]CronTask, len(configs))
	for _, config := range configs {
		index[config.Name] = config
	}
	for i := 0; i < len(tasks); i++ {
		if config, ok := index[tasks[i].Name]; ok {
			tasks[i].Spec = config.Spec
		}
	}
}
