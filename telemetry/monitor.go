package telemetry

import (
	"expvar"
	"sync"
)

var monitor = struct {
	sync.Mutex
	enabled bool
	state   map[string]*expvar.Int
}{
	state: make(map[string]*expvar.Int),
}

func monitorTouch(name string) {
	if !monitor.enabled {
		return
	}

	monitor.Lock()
	defer monitor.Unlock()

	if v, ok := monitor.state[name]; ok {
		v.Add(1)
		return
	}
	monitor.state[name] = expvar.NewInt(name)
}
