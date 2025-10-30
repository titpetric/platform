package service

import (
	"context"
	"net/http"

	"github.com/titpetric/platform/module/user/model"
)

// Login handles user authentication via HTML form submission.
func (h *Service) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		h.Error(r, "Email and Password are required", nil)
		h.LoginView(w, r)
		return
	}

	user, err := h.UserStorage.Authenticate(context.Background(), model.UserAuth{
		Email:    email,
		Password: password,
	})
	if err != nil || !user.IsActive() {
		h.Error(r, "Invalid credentials for login", err)
		h.LoginView(w, r)
		return
	}

	session, err := h.SessionStorage.Create(context.Background(), user.ID)
	if err != nil {
		h.Error(r, "Can't create session", err)
		h.LoginView(w, r)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // set false for local dev if needed
		Expires:  *session.ExpiresAt,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
