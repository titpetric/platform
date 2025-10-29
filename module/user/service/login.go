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
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.UserStorage.Authenticate(context.Background(), model.UserAuth{
		Email:    email,
		Password: password,
	})
	if err != nil {
		h.Error(r, "Invalid credentials for login", err)
		h.LoginView(w, r)
		return
	}

	session, err := h.SessionStorage.Create(context.Background(), user.ID)
	if err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
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
