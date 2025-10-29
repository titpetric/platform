package user

import (
	"net/http"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/module/user/model"
	"github.com/titpetric/platform/module/user/storage"
)

type userKey struct{}

var userContext = internal.NewContextValue[*model.User](userKey{})

// IsLoggedIn returns true or false. Any errors are swallowed, returning false.
func IsLoggedIn(r *http.Request) bool {
	if user := userContext.Get(r); user != nil {
		return user.IsActive()
	}

	user, err := User(r)
	if user == nil || err != nil {
		return false
	}

	return user.IsActive()
}

// User will return the user bound to the session. If no user is bound to
// the session or there is no session, the function will return nil, nil.
func User(r *http.Request) (*model.User, error) {
	if user := userContext.Get(r); user != nil {
		return user, nil
	}

	ctx := r.Context()

	db, err := DB()
	if err != nil {
		return nil, err
	}

	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		userStorage := storage.NewUserStorage(db)
		sessionStorage := storage.NewSessionStorage(db)

		session, err := sessionStorage.Get(ctx, cookie.Value)
		if err != nil {
			return nil, err
		}

		user, err := userStorage.GetUser(ctx, session.UserID)
		if err != nil {
			return nil, err
		}

		userContext.Set(r, user)

		return user, nil
	}

	return nil, nil
}
