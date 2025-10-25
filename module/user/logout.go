package user

import (
	"net/http"
)

// Logout deletes the session cookie and optionally the session in storage.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("session_id")

	if err == nil && cookie.Value != "" {
		// Delete the session from the database
		_ = h.SessionStorage.Delete(ctx, cookie.Value)

		// Clear cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
