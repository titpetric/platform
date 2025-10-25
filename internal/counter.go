package internal

import (
	"expvar"
)

// Counter encapsulates an int32 counter and exposes it via expvar.
type Counter struct {
	*expvar.Int
}

// NewCounter constructs a Counter, registers it to expvar with the given name, and returns a pointer.
func NewCounter(name string) *Counter {
	c := &Counter{
		Int: expvar.NewInt(name),
	}
	return c
}

// Inc increments the counter atomically.
func (c *Counter) Inc() {
	c.Add(1)
}
