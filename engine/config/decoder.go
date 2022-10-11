package config

import (
	"fmt"
	"github.com/swamp-labs/swamp/engine/httpreq"
	"github.com/swamp-labs/swamp/engine/simulation"
	volume "github.com/swamp-labs/swamp/engine/volume"
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

func (v volumeTemplate) decode() (volume.Volume, error) {
	return volume.MakeVolume(v.RequestGroup), nil
}

func parseGroups(groupTemplates map[string][]requestTemplate) (map[string][]httpreq.Request, error) {
	groups := make(map[string][]httpreq.Request)

	for groupId, requestsTemplates := range groupTemplates {
		requests := make([]httpreq.Request, cap(requestsTemplates), len(requestsTemplates))
		for i, requestTemplate := range requestsTemplates {
			request, err := requestTemplate.decode()
			if err != nil {
				return nil, err
			}
			requests[i] = *request
		}
		groups[groupId] = requests
	}
	return groups, nil
}

func (s simulationTemplate) decode() (simulation.Simulation, error) {
	volumes := make([]volume.Volume, cap(s.Volumes), len(s.Volumes))
	for i, volumeTpl := range s.Volumes {
		v, err := volumeTpl.decode()
		if err != nil {
			return nil, err
		}
		volumes[i] = v
	}

	groups, err := parseGroups(s.Groups)
	if err != nil {
		return nil, fmt.Errorf("fail to parse groups : %v", err)
	}

	return simulation.MakeSimulation(volumes, groups), nil
}
