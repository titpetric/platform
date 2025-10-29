package user

import (
	"io/fs"

	"github.com/titpetric/platform"
	"github.com/titpetric/platform/module/user/service"
	"github.com/titpetric/platform/module/user/storage"
)

type Handler struct {
	Service *service.Service
}

// NewHandler sets up dependencies and produces a handler.
func NewHandler(themeFS fs.FS) (*Handler, error) {
	db, err := platform.Database.Connect()
	if err != nil {
		return nil, err
	}

	userStorage := storage.NewUserStorage(db)
	sessionStorage := storage.NewSessionStorage(db)

	options := &service.Options{
		UserStorage:    userStorage,
		SessionStorage: sessionStorage,
		ThemeFS:        themeFS,
		ModuleFS:       TemplateFS,
	}

	svc, err := service.NewService(options)
	if err != nil {
		return nil, err
	}

	return &Handler{
		Service: svc,
	}, nil
}

// Name returns the name of the containing package.
func (h *Handler) Name() string {
	return "user"
}

// Mount registers login, logout, and register routes.
func (h *Handler) Mount(r platform.Router) {
	h.Service.Mount(r)
}

// Close implements a closer for graceful shutdown.
func (h *Handler) Close() {
	h.Service.Close()
}
