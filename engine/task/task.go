package task

import (
	"github.com/swamp-labs/swamp/engine/httpreq"
	"github.com/swamp-labs/swamp/engine/volume"
)

type Task interface {
	GetRequest() []httpreq.Request
	GetVolume() volume.Volume
}

type task struct {
	requests []httpreq.Request
	volume   volume.Volume
}

func (t task) GetRequest() []httpreq.Request {
	return t.requests
}

func (t task) GetVolume() volume.Volume {
	return t.volume
}

func MakeTask(requests []httpreq.Request, volume volume.Volume) Task {
	return task{
		requests: requests,
		volume:   volume,
	}
}
