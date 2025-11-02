package telemetry

import (
	"expvar"
	"sync"
)

var monitor = NewMonitor()

// Monitor adds expvar instrumentation over the telemetry API.
type Monitor struct {
	mu      sync.Mutex
	enabled bool
	state   map[string]*expvar.Int
}

func NewMonitor() *Monitor {
	return &Monitor{
		state: make(map[string]*expvar.Int),
	}
}

func (m *Monitor) SetEnabled(enabled bool) {
	m.enabled = enabled
}

func (m *Monitor) Enabled() bool {
	return m.enabled
}

func (m *Monitor) Touch(name string) {
	if !m.enabled {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if v, ok := m.state[name]; ok {
		v.Add(1)
		return
	}
	m.state[name] = expvar.NewInt(name)
}
