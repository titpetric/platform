package expvar

import (
	"expvar"
	"sync"
	"time"

	"github.com/titpetric/platform"
)

var publishMu sync.Mutex

type Handler struct {
	platform.UnimplementedModule
}

func NewHandler() *Handler {
	return &Handler{}
}

func (m *Handler) Start() error {
	publishMu.Lock()
	defer publishMu.Unlock()

	start := time.Now()
	if expvar.Get("uptime") == nil {
		expvar.Publish("uptime", expvar.Func(func() interface{} {
			return time.Since(start).Seconds()
		}))
	}
	return nil
}

func (m *Handler) Mount(r platform.Router) error {
	r.Mount("/debug/vars", expvar.Handler())
	return nil
}
