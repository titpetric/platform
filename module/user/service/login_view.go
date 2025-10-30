package service

import (
	"log"
	"net/http"

	"github.com/titpetric/platform/module/theme"
	"github.com/titpetric/platform/module/user/model"
)

// LoginView renders login.tpl when no valid session exists,
// or logout.tpl with the full user model when a valid session is found.
func (h *Service) LoginView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type templateData struct {
		Theme   *theme.Options
		User    *model.User
		Session *model.UserSession

		ErrorMessage string
		Form         map[string]string
	}

	var data templateData = templateData{
		Theme:        theme.NewOptions(),
		ErrorMessage: h.GetError(r),
		Form: map[string]string{
			"email": r.FormValue("email"),
		},
	}

	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		if session, err := h.SessionStorage.Get(ctx, cookie.Value); err == nil {
			if user, err := h.UserStorage.Get(ctx, session.UserID); err == nil {
				data.User = user
				data.Session = session

				h.View(w, "logout.tpl", data)
				return
			}
			log.Println(err)
		} else {
			log.Println(err)
		}
	}

	h.View(w, "login.tpl", data)
}
