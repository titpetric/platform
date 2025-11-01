package service

import (
	"net/http"

	"github.com/titpetric/platform/module/theme"
	"github.com/titpetric/platform/telemetry"
)

// RegisterView renders the registration page.
func (h *Service) RegisterView(w http.ResponseWriter, r *http.Request) {
	r, span := telemetry.StartRequest(r, "user.service.RegisterView")
	defer span.End()

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

	h.View(w, r, "register.tpl", data)
}
