package user

import (
	"context"
	"net/http"

	"github.com/titpetric/platform/module/theme"
	"github.com/titpetric/platform/module/user/model"
)

// LoginView renders login.tpl when no valid session exists,
// or logout.tpl with the full user model when a valid session is found.
func (h *Handler) LoginView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type templateData struct {
		Theme   *theme.Options
		User    *model.User
		Session *model.UserSession
	}

	var data templateData = templateData{
		Theme: theme.NewOptions(),
	}

	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		if session, err := h.SessionStorage.Get(ctx, cookie.Value); err == nil {
			if user, err := h.UserStorage.GetUser(ctx, session.UserID); err == nil {
				data.User = user
				data.Session = session

				h.View(w, "logout.tpl", data)
				return
			}
		}
	}

	h.View(w, "login.tpl", data)
}

// Login handles user authentication via HTML form submission.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form submission", http.StatusBadRequest)
		return
	}

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
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
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
