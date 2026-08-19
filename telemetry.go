package platform

import (
	"github.com/titpetric/oida"
)

// Telemetry configures the recorder and the debug dashboard that New registers.
// It is the oida options, named in the platform namespace so the field type is
// ours and later additions do not change the shape of Options.
//
// Enabled is oida's own field, read here as the single switch: a disabled
// Telemetry registers no module, so there is no dashboard route and no
// middleware, rather than a mounted dashboard with nothing behind it. Set it in
// a host that registers its own recorder, otherwise both mount the same path
// and the router panics on the duplicate route.
//
// Retention defaults to an in-memory ring buffer sized by RingBufferSize.
// Assign Storage to keep traces across a restart:
//
//	store, err := oida.NewStorageDisk(500, "/var/lib/app/traces")
//	if err != nil {
//		return err
//	}
//	options.Telemetry.Storage = store
type Telemetry struct {
	oida.Options `yaml:",inline"`
}

// NewTelemetry returns the default telemetry configuration.
func NewTelemetry() Telemetry {
	return Telemetry{Options: oida.NewOptions()}
}
