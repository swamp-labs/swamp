package task

import (
	"github.com/swamp-labs/swamp/engine/httpreq"
	"github.com/swamp-labs/swamp/engine/session"
	"github.com/swamp-labs/swamp/engine/timedFunctions"
	"github.com/swamp-labs/swamp/engine/volume"
	"time"
)

type Task interface {
	GetRequest() []httpreq.Request
	GetVolume() []volume.Volume
	Execute(ch chan map[string][]httpreq.Sample)
}

type task struct {
	requests []httpreq.Request
	volume   []volume.Volume
}

func (t task) GetRequest() []httpreq.Request {
	return t.requests
}

func (t task) GetVolume() []volume.Volume {
	return t.volume
}

func MakeTask(requests []httpreq.Request, volume []volume.Volume) Task {
	return task{
		requests: requests,
		volume:   volume,
	}
}

func (t task) Execute(ch chan map[string][]httpreq.Sample) {
	sessionVar := make(map[string]string)
	taskReport := make(map[string][]httpreq.Sample)
	var test timedFunctions.TickedFunction = func(_ timedFunctions.TickFunctionParameter) interface{} {
		var s session.Session
		s.New()
		for _, r := range t.GetRequest() {
			r.Execute(sessionVar)
		}
		return true
	}
	duration, _ := time.ParseDuration("5s")
	timedFunctions.Constant(test, 2, duration)

	ch <- taskReport
}
