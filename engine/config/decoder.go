package config

import (
	"fmt"
	"github.com/swamp-labs/swamp/engine/httpreq"
	"github.com/swamp-labs/swamp/engine/simulation"
	"github.com/swamp-labs/swamp/engine/task"
	"github.com/swamp-labs/swamp/engine/volume"
)

func (requestTemplate requestTemplate) decode() (*httpreq.Request, error) {
	assertions, err := requestTemplate.Assertions.decode()
	if err != nil {
		return nil, err
	}

	return &httpreq.Request{
		Name:            requestTemplate.Name,
		Method:          requestTemplate.Method,
		Protocol:        requestTemplate.Protocol,
		Headers:         requestTemplate.Headers,
		URL:             requestTemplate.URL,
		Body:            requestTemplate.Body,
		QueryParameters: requestTemplate.QueryParameters,
		Assertions:      assertions,
	}, nil
}

func (t taskTemplate) decode() (task.Task, error) {
	volumes := make([]volume.Volume, cap(t.Volume), len(t.Volume))

	for i, node := range t.Volume {
		v, err := yamlNodeToVolume(node)
		if err != nil {
			return nil, fmt.Errorf("fail to parse volume : %v", err)
		}
		volumes[i] = v
	}

	requests, err := parseRequests(t.Requests)
	if err != nil {
		return nil, err
	}
	return task.MakeTask(requests, volumes), nil
}

//func parseVolume(volumeBlock []volumeTemplate) ([]volume.Volume, error) {
//	volumes := make([]map[string]int, cap(volumeTemplate), len(volumeTemplate))
//
//	for volumeId, v := range volumeBlock {
//		volumes, err := v.decode()
//		if err != nil {
//			return nil, err
//		}
//	}
//	return volumes, nil
//}

func parseRequests(requestsTemplate []requestTemplate) ([]httpreq.Request, error) {
	requests := make([]httpreq.Request, cap(requestsTemplate), len(requestsTemplate))
	for requestId, r := range requestsTemplate {
		request, err := r.decode()
		if err != nil {
			return nil, err
		}
		requests[requestId] = *request
	}
	return requests, nil
}

func parseTasks(taskTemplates map[string]taskTemplate) (map[string]task.Task, error) {
	tasks := make(map[string]task.Task)

	for taskId, taskTemplate := range taskTemplates {
		var err error
		tasks[taskId], err = taskTemplate.decode()
		if err != nil {
			return nil, err
		}
	}

	return tasks, nil
}

func (s simulationTemplate) decode() (simulation.Simulation, error) {

	tasks, err := parseTasks(s.Tasks)
	if err != nil {
		return nil, fmt.Errorf("fail to parse tasks : %v", err)
	}

	return simulation.MakeSimulation(tasks), nil
}
