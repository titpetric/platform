package service

import (
	"net/http"

	"github.com/titpetric/platform/module/theme"
)

// RegisterView renders the registration page.
func (h *Service) RegisterView(w http.ResponseWriter, r *http.Request) {
	type templateData struct {
		Theme *theme.Options
	}
	data := templateData{
		Theme: theme.NewOptions(),
	}

	h.View(w, "register.tpl", data)
}
