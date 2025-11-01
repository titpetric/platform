package service

import (
	"net/http"

	"github.com/titpetric/platform/module/user/model"
	"github.com/titpetric/platform/telemetry"
)

// Register handles creating a new user and starting a session via HTML form submission.
func (h *Service) Register(w http.ResponseWriter, r *http.Request) {
	r, span := telemetry.StartRequest(r, "user.service.Register")
	defer span.End()

	ctx := r.Context()

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if firstName == "" || lastName == "" || email == "" || password == "" {
		h.Error(r, "All fields are required", nil)
		return
	}

	user := &model.User{
		FirstName: firstName,
		LastName:  lastName,
	}

	auth := &model.UserAuth{
		Email:    email,
		Password: password,
	}

	createdUser, err := h.UserStorage.Create(ctx, user, auth)
	if err != nil {
		h.Error(r, "Failed to create user", err)
		h.RegisterView(w, r)
		return
	}

	session, err := h.SessionStorage.Create(ctx, createdUser.ID)
	if err != nil {
		h.Error(r, "Failed to create session", err)
		h.RegisterView(w, r)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Expires:  *session.ExpiresAt,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
