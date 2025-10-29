package service

import (
	"net/http"
)

// LogoutView just wraps LoginView. The view changes based on if
// the user is logged in already, allowing them to log in or out.
func (h *Service) LogoutView(w http.ResponseWriter, r *http.Request) {
	h.LoginView(w, r)
}
