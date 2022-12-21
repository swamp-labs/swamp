package volume

import (
	"github.com/swamp-labs/swamp/engine/timedFunctions"
	"time"
)

type Volume interface {
	Inject(fct timedFunctions.TickedFunction)
}

type wait struct {
	Wait time.Duration `yaml:"wait"`
}

type constant struct {
	ConstantUserPerSec int           `yaml:"constantUserPerSec"`
	During             time.Duration `yaml:"during"`
}

type increase struct {
	IncreaseUserPerSec variation     `yaml:"increaseUserPerSec"`
	During             time.Duration `yaml:"during"`
}

type decrease struct {
	DecreaseUserPerSec variation     `yaml:"decreaseUserPerSec"`
	During             time.Duration `yaml:"during"`
}

type variation struct {
	from int `yaml:"from"`
	to   int `yaml:"to"`
}

type instant struct {
	instant int `yaml:"instant"`
}

func (w wait) Inject(_ timedFunctions.TickedFunction) {
	time.Sleep(w.Wait)
}

func (c constant) Inject(fct timedFunctions.TickedFunction) {
	timedFunctions.Constant(fct, c.ConstantUserPerSec, c.During)
}

func (i increase) Inject(fct timedFunctions.TickedFunction) {
	timedFunctions.Linear(fct, i.IncreaseUserPerSec.from, i.IncreaseUserPerSec.to, i.During)
}

func (d decrease) Inject(fct timedFunctions.TickedFunction) {
	timedFunctions.Linear(fct, d.DecreaseUserPerSec.from, d.DecreaseUserPerSec.to, d.During)
}

func (i instant) Inject(_ timedFunctions.TickedFunction) {

}

func NewWaitVolume(duration time.Duration) Volume {
	return wait{
		Wait: duration,
	}
}

func NewConstantVolume(cst int, duration time.Duration) Volume {
	return constant{
		ConstantUserPerSec: cst,
		During:             duration,
	}
}

func NewIncreaseVolume(from int, to int, duration time.Duration) Volume {
	return increase{
		IncreaseUserPerSec: variation{from: from, to: to},
		During:             duration,
	}
}

func NewDecreaseVolume(from int, to int, duration time.Duration) Volume {
	return decrease{
		DecreaseUserPerSec: variation{from: from, to: to},
		During:             duration,
	}
}
