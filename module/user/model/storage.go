package model

import (
	"context"
)

type SessionStorage interface {
	Create(ctx context.Context, userID string) (*UserSession, error)
	Get(ctx context.Context, sessionID string) (*UserSession, error)
	Delete(ctx context.Context, sessionID string) error
}

type UserStorage interface {
	Create(context.Context, *User, *UserAuth) (*User, error)
	Update(context.Context, *User) (*User, error)

	Get(context.Context, string) (*User, error)
	GetGroups(context.Context, string) ([]UserGroup, error)

	Authenticate(ctx context.Context, auth UserAuth) (*User, error)
}
