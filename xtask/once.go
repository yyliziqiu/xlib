package xtask

import (
	"context"
	"fmt"
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
	index := make(map[string]OnceTask, len(configs))
	for _, config := range configs {
		index[config.Name] = config
	}

	tasks := tasksFunc()
	for i := 0; i < len(tasks); i++ {
		if config, ok := index[tasks[i].Name]; ok {
			tasks[i].GON = config.GON
		}
	}

	StartOnceTasks(ctx, func() []OnceTask { return tasks })
}

func RunOnceTasksByConfig(ctx context.Context, tasksFunc func() []OnceTask, configs []OnceTask) error {
	index := make(map[string]OnceTask)
	for _, task := range tasksFunc() {
		index[task.Name] = task
	}

	runnable := make([]OnceTask, 0)
	for _, config := range configs {
		task, ok := index[config.Name]
		if !ok {
			return fmt.Errorf("not found once task[%s]", config.Name)
		}
		task.GON = config.GON
		runnable = append(runnable, task)
	}

	StartOnceTasks(ctx, func() []OnceTask { return runnable })

	return nil
}
