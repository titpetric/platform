package user

import (
	"github.com/titpetric/platform"
	"github.com/titpetric/platform/module/theme"
	"github.com/titpetric/platform/module/user/service"
	"github.com/titpetric/platform/module/user/storage"
)

// Handler implements a module contract.
type Handler struct {
	Service *service.Service
}

// Verify contract.
var _ platform.Module = (*Handler)(nil)

// NewHandler sets up dependencies and produces a handler.
func NewHandler() *Handler {
	return &Handler{}
}

// Start will initialize the service to handle requests.
func (h *Handler) Start() error {
	db, err := DB()
	if err != nil {
		return err
	}

	themeFS := theme.TemplateFS
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
		return err
	}

	h.Service = svc
	return nil
}

// Name returns the name of the containing package.
func (h *Handler) Name() string {
	return "user"
}

// Mount registers login, logout, and register routes.
func (h *Handler) Mount(r platform.Router) {
	h.Service.Mount(r)
}

// Stop implements a closer for graceful shutdown.
func (h *Handler) Stop() error {
	h.Service.Close()
	return nil
}
