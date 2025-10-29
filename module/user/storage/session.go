package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/platform/internal"
	"github.com/titpetric/platform/module/user/model"
)

// SessionStorage implements session persistence using MySQL.
type SessionStorage struct {
	db *sqlx.DB
}

// NewSessionStorage creates a new SessionStorage.
func NewSessionStorage(db *sqlx.DB) *SessionStorage {
	return &SessionStorage{db: db}
}

// Create inserts a new session for the given userID.
func (s *SessionStorage) Create(ctx context.Context, userID string) (*model.UserSession, error) {
	now := time.Now()
	session := &model.UserSession{
		ID:     internal.ULID(),
		UserID: userID,
	}
	session.SetCreatedAt(now)
	session.SetExpiresAt(now.Add(24 * time.Hour)) // default 24h expiration

	query := `INSERT INTO user_session (id, user_id, expires_at, created_at) VALUES (?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, session.ID, session.UserID, session.ExpiresAt, session.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return session, nil
}

// Get retrieves a session by sessionID.
// Returns model.ErrSessionExpired if the session has expired.
func (s *SessionStorage) Get(ctx context.Context, sessionID string) (*model.UserSession, error) {
	query := `SELECT * FROM user_session WHERE id=?`
	session := &model.UserSession{}
	if err := s.db.GetContext(ctx, session, query, sessionID); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("get session: %w", err)
	}

	if time.Now().After(*session.ExpiresAt) {
		return nil, model.ErrSessionExpired
	}

	return session, nil
}

// Delete removes a session by sessionID.
func (s *SessionStorage) Delete(ctx context.Context, sessionID string) error {
	query := `DELETE FROM user_session WHERE id=?`
	_, err := s.db.ExecContext(ctx, query, sessionID)
	return fmt.Errorf("delete session: %w", err)
}

var _ model.SessionStorage = (*SessionStorage)(nil)
