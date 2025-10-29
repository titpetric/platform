package service

import (
	"net/http"

	"github.com/titpetric/platform/module/theme"
)

// RegisterView renders the registration page.
func (h *Service) RegisterView(w http.ResponseWriter, r *http.Request) {
	type templateData struct {
		Theme *theme.Options

		ErrorMessage string
		Form         map[string]string
	}

	data := templateData{
		Theme:        theme.NewOptions(),
		ErrorMessage: h.GetError(r),
		Form: map[string]string{
			"first_name": r.FormValue("first_name"),
			"last_name":  r.FormValue("last_name"),
			"email":      r.FormValue("email"),
		},
	}

	h.View(w, "register.tpl", data)
}
