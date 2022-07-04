package simulation

import (
	"github.com/swamp-labs/swamp/engine/httpreq"
	"github.com/swamp-labs/swamp/engine/volume"
)

type Simulation interface {
	GetVolumes() []volume.Volume
	GetGroups() map[string][]httpreq.Request
}

type simulation struct {
	volumes []volume.Volume
	groups  map[string][]httpreq.Request
}

func (s simulation) GetVolumes() []volume.Volume {
	return s.volumes
}

func (s simulation) GetGroups() map[string][]httpreq.Request {
	return s.groups
}

func MakeSimulation(volumes []volume.Volume, groups map[string][]httpreq.Request) Simulation {
	return simulation{
		volumes: volumes,
		groups:  groups,
	}
}
