package platform

type App interface {
	Start() error
	Stop() error
}
