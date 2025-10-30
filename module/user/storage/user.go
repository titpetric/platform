package storage

import (
	"context"
	"database/sql"
	"errors"
	"expvar"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/module/user/model"
)

// UserStorage implements the model.Storage interface using MySQL via sqlx.
type UserStorage struct {
	db *sqlx.DB

	monitor UserStorageMonitor
}

type UserStorageMonitor struct {
	Create       *expvar.Int
	Update       *expvar.Int
	Get          *expvar.Int
	GetGroups    *expvar.Int
	Authenticate *expvar.Int
}

func NewUserStorageMonitor() UserStorageMonitor {
	return UserStorageMonitor{
		Create:       expvar.NewInt("user.storage.user.create"),
		Update:       expvar.NewInt("user.storage.user.update"),
		Get:          expvar.NewInt("user.storage.user.get"),
		GetGroups:    expvar.NewInt("user.storage.user.getgroups"),
		Authenticate: expvar.NewInt("user.storage.user.authenticate"),
	}
}

// NewUserStorage returns a new UserStorage backed by the given sqlx.DB.
func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{
		db:      db,
		monitor: NewUserStorageMonitor(),
	}
}

// Create inserts a new user and their authentication credentials.
// Returns an error if authentication information is missing.
func (s *UserStorage) Create(ctx context.Context, u *model.User, userAuth *model.UserAuth) (*model.User, error) {
	defer s.monitor.Create.Add(1)

	if userAuth == nil || userAuth.Email == "" || userAuth.Password == "" {
		return nil, errors.New("missing authentication info: email and password are required")
	}

	now := time.Now()
	u.SetCreatedAt(now)
	u.SetUpdatedAt(now)
	u.ID = internal.ULID()

	hashed, err := bcrypt.GenerateFromPassword([]byte(userAuth.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	err = internal.Transaction(s.db, func(tx *sqlx.Tx) error {
		var err error
		userQuery := `
			INSERT INTO user
			(id, first_name, last_name, deleted_at, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`
		_, err = s.db.ExecContext(ctx, userQuery, u.ID, u.FirstName, u.LastName, u.DeletedAt, u.CreatedAt, u.UpdatedAt)
		if err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		authQuery := `
			INSERT INTO user_auth
				(user_id, email, password, created_at, updated_at)
			VALUES
				(?, ?, ?, ?, ?)
		`
		_, err = s.db.ExecContext(ctx, authQuery, u.ID, userAuth.Email, hashed, now, now)
		if err != nil {
			return fmt.Errorf("create user_auth: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.Get(ctx, u.ID)
}

// Update modifies an existing user and updates the updated_at timestamp.
func (s *UserStorage) Update(ctx context.Context, u *model.User) (*model.User, error) {
	defer s.monitor.Update.Add(1)

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

// Get retrieves a user by ULID.
func (s *UserStorage) Get(ctx context.Context, id string) (*model.User, error) {
	defer s.monitor.Get.Add(1)

	u := &model.User{}
	query := `SELECT * FROM user WHERE id=?`
	if err := s.db.GetContext(ctx, u, query, id); err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

// GetGroups returns all groups the user belongs to.
func (s *UserStorage) GetGroups(ctx context.Context, userID string) ([]model.UserGroup, error) {
	defer s.monitor.GetGroups.Add(1)

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
	defer s.monitor.Authenticate.Add(1)

	query := `SELECT user_id, password FROM user_auth WHERE email=? LIMIT 1`

	dbAuth := &model.UserAuth{}
	if err := s.db.GetContext(ctx, dbAuth, query, auth.Email); err != nil {
		return nil, fmt.Errorf("authenticate lookup: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbAuth.Password), []byte(auth.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("bcrypt compare: %w", err)
	}

	user, err := s.Get(ctx, dbAuth.UserID)
	if err != nil {
		return nil, fmt.Errorf("authenticate get user: %w", err)
	}

	return user, nil
}

var _ model.UserStorage = (*UserStorage)(nil)
