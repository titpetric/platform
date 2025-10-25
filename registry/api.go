package registry

var global = struct {
	registry *Registry
}{
	registry: NewRegistry(),
}

var (
	Add           = global.registry.Add
	AddModule     = global.registry.AddModule
	AddMiddleware = global.registry.AddMiddleware
	Mount         = global.registry.Mount
	Close         = global.registry.Close
	Clone         = global.registry.Clone
)
