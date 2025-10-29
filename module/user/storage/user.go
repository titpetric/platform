package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/module/user/model"
)

// UserStorage implements the model.Storage interface using MySQL via sqlx.
type UserStorage struct {
	// model.UnimplementedStorage
	db *sqlx.DB
}

// NewUserStorage returns a new UserStorage backed by the given sqlx.DB.
func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{db: db}
}

// Create inserts a new user and their authentication credentials.
// Returns an error if authentication information is missing.
func (s *UserStorage) Create(ctx context.Context, u *model.User, userAuth *model.UserAuth) (*model.User, error) {
	if userAuth == nil || userAuth.Email == "" || userAuth.Password == "" {
		return nil, errors.New("missing authentication info: email and password are required")
	}

	now := time.Now()
	u.SetCreatedAt(now)
	u.SetUpdatedAt(now)
	u.ID = internal.ULID()

	userQuery := `
		INSERT INTO user
		(id, first_name, last_name, deleted_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	if _, err := s.db.ExecContext(ctx, userQuery,
		u.ID, u.FirstName, u.LastName, u.DeletedAt, u.CreatedAt, u.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(userAuth.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	authQuery := `
		INSERT INTO user_auth
			(user_id, email, password, created_at, updated_at)
		VALUES
			(?, ?, ?, ?, ?)
	`
	if _, err := s.db.ExecContext(ctx, authQuery,
		u.ID, userAuth.Email, hashed, now, now,
	); err != nil {
		return nil, fmt.Errorf("create user_auth: %w", err)
	}

	return s.GetUser(ctx, u.ID)
}

// Update modifies an existing user and updates the updated_at timestamp.
func (s *UserStorage) Update(ctx context.Context, u *model.User) (*model.User, error) {
	u.SetUpdatedAt(time.Now())

	query := `UPDATE user SET first_name=?, last_name=?, deleted_at=?, updated_at=? WHERE id=?`

	_, err := s.db.ExecContext(ctx, query,
		u.FirstName, u.LastName, u.DeletedAt, u.UpdatedAt, u.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return u, nil
}

// GetUser retrieves a user by ULID.
func (s *UserStorage) GetUser(ctx context.Context, id string) (*model.User, error) {
	query := `SELECT * FROM user WHERE id=?`
	var u *model.User
	if err := s.db.GetContext(ctx, u, query, id); err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

// GetUserGroups returns all groups the user belongs to.
func (s *UserStorage) GetUserGroups(ctx context.Context, userID string) ([]model.UserGroup, error) {
	query := `
		SELECT g.id, g.title, g.created_at, g.updated_at
		FROM user_group g
		JOIN user_group_member m ON m.group_id = g.id
		WHERE m.user_id = ?
	`
	var groups []model.UserGroup
	if err := s.db.SelectContext(ctx, &groups, query, userID); err != nil {
		return nil, fmt.Errorf("get user groups: %w", err)
	}
	return groups, nil
}

// Authenticate verifies a user's credentials using bcrypt and returns the user.
func (s *UserStorage) Authenticate(ctx context.Context, auth model.UserAuth) (*model.User, error) {
	query := `
		SELECT a.user_id, a.password
		FROM user_auth a
		WHERE a.email = ?
		LIMIT 1
	`

	var dbAuth *model.UserAuth
	if err := s.db.GetContext(ctx, dbAuth, query, auth.Email); err != nil {
		return nil, fmt.Errorf("authenticate lookup: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbAuth.Password), []byte(auth.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("bcrypt compare: %w", err)
	}

	user, err := s.GetUser(ctx, dbAuth.UserID)
	if err != nil {
		return nil, fmt.Errorf("authenticate get user: %w", err)
	}

	return user, nil
}
