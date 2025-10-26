package platform

var global = struct {
	registry *Registry
}{
	registry: NewRegistry(),
}

var (
	AddModule     = global.registry.AddModule
	AddMiddleware = global.registry.AddMiddleware
)
