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

// NewMonitor creates a new Monitor for telemetry tracking.
func NewMonitor() *Monitor {
	return &Monitor{
		state: make(map[string]*expvar.Int),
	}
}

// SetEnabled enables or disables telemetry monitoring.
func (m *Monitor) SetEnabled(enabled bool) {
	m.enabled = enabled
}

// Enabled reports whether telemetry monitoring is enabled.
func (m *Monitor) Enabled() bool {
	return m.enabled
}

// Touch increments the counter for a telemetry event name.
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
