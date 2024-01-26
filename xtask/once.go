package xtask

import (
	"context"
)

type OnceTask struct {
	Name string
	GON  int
	Cmd  func(ctx context.Context)
}

func StartOnceTasks(ctx context.Context, tasksFunc func() []OnceTask) {
	for _, task := range tasksFunc() {
		if task.GON <= 0 {
			continue
		}
		for i := 0; i < task.GON; i++ {
			go task.Cmd(ctx)
		}
		Logger.Infof("Add once task: %s (%d).", task.Name, task.GON)
	}
}

func StartOnceTasksWithConfig(ctx context.Context, tasksFunc func() []OnceTask, configs []OnceTask) {
	tasks := tasksFunc()

	ConfigOnceTasks(tasks, configs)

	StartOnceTasks(ctx, func() []OnceTask { return tasks })
}

func ConfigOnceTasks(tasks []OnceTask, configs []OnceTask) {
	index := make(map[string]OnceTask, len(configs))
	for _, config := range configs {
		index[config.Name] = config
	}
	for i := 0; i < len(tasks); i++ {
		if config, ok := index[tasks[i].Name]; ok {
			tasks[i].GON = config.GON
		}
	}
}
