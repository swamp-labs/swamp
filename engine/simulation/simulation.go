package simulation

import (
	"fmt"
	"github.com/swamp-labs/swamp/engine/httpreq"
	"github.com/swamp-labs/swamp/engine/task"
)

type Simulation interface {
	GetTasks() map[string]task.Task
	Run()
}

type simulation struct {
	tasks           map[string]task.Task
	executionTraces map[string][]httpreq.Sample
}

func (s simulation) GetTasks() map[string]task.Task {
	return s.tasks
}

func MakeSimulation(tasks map[string]task.Task) Simulation {
	return simulation{
		tasks: tasks,
	}
}

func (s simulation) Run() {

	for _, t := range s.tasks {
		ch := make(chan map[string][]httpreq.Sample)

		go t.Execute(ch)
		fmt.Println("Result: ", <-ch)
	}

}

func (s simulation) AddTraces(name string, sample httpreq.Sample) {
	s.executionTraces[name] = append(s.executionTraces[name], sample)
}
