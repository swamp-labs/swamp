package simulation

import (
	"github.com/swamp-labs/swamp/engine/task"
)

type Simulation interface {
	GetTasks() map[string]task.Task
}

type simulation struct {
	tasks map[string]task.Task
}

func (s simulation) GetTasks() map[string]task.Task {
	return s.tasks
}

func MakeSimulation(tasks map[string]task.Task) Simulation {
	return simulation{
		tasks: tasks,
	}
}
