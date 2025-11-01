package service

import (
	"net/http"

	"github.com/titpetric/platform/telemetry"
)

// LogoutView just wraps LoginView. The view changes based on if
// the user is logged in already, allowing them to log in or out.
func (h *Service) LogoutView(w http.ResponseWriter, r *http.Request) {
	r, span := telemetry.StartRequest(r, "user.service.LogoutView")
	defer span.End()

	h.LoginView(w, r)
}
