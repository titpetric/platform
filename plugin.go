package platform

type Plugin struct {
	name     string
	hook     func(Router)
	shutdown func()
}

func NewPlugin(name string, hook func(Router), shutdown func()) Plugin {
	return Plugin{
		name:     name,
		hook:     hook,
		shutdown: shutdown,
	}
}

func (p *Plugin) Name() string {
	return p.name
}

func (p *Plugin) Mount(r Router) {
	if p.hook != nil {
		p.hook(r)
	}
}

func (p *Plugin) Close() {
	if p.shutdown != nil {
		p.shutdown()
	}
}
