package repository

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/gs-hms/authentication/database"
	"github.com/gs-hms/authentication/model"
)

// AuthenticationSessionRepository is an interface for interacting with the authentication sessions table.
type AuthenticationSessionRepository interface {

	//CreateSession creates a new authentication session for a user.
	CreateSession(ctx context.Context, userID uint64, refreshToken string, expiresAt time.Time) error

	//GetActiveSessionByUserID gets all active authentication sessions for a user.
	GetActiveSessionByUserID(ctx context.Context, userID uint64) ([]*model.AuthenticationSession, error)

	//GetActiveSessionByRefreshToken gets an authentication session by its refresh token.
	GetActiveSessionByRefreshToken(ctx context.Context, refreshToken string) (*model.AuthenticationSession, error)

	//RevokeSession revokes an authentication session.
	RevokeSession(ctx context.Context, sessionID uint64) error
}

type authenticationSessionRepository struct {
	db *database.Postgres
}

// NewAuthenticationSessionRepository creates a new instance of the AuthenticationSessionRepository.
func NewAuthenticationSessionRepository(db *database.Postgres) AuthenticationSessionRepository {
	return &authenticationSessionRepository{
		db: db,
	}
}

// CreateSession creates a new authentication session for a user.
func (r *authenticationSessionRepository) CreateSession(ctx context.Context, userID uint64, refreshToken string, expiresAt time.Time) error {
	qry, args, err := sq.Insert(model.AuthenticationSessionTableName).
		Columns("user_id", "refresh_token", "expired_at").
		Values(userID, refreshToken, expiresAt).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("build create session query: %w", err)
	}

	_, err = r.db.Pool.Exec(ctx, qry, args...)
	if err != nil {
		return fmt.Errorf("execute create session query: %w", err)
	}
	return nil
}

// GetActiveSessionByUserID gets all active authentication sessions for a user.
func (r *authenticationSessionRepository) GetActiveSessionByUserID(ctx context.Context, userID uint64) ([]*model.AuthenticationSession, error) {
	qry, args, err := sq.Select("id", "user_id", "refresh_token", "expired_at").
		From(model.AuthenticationSessionTableName).
		Where(sq.Eq{"user_id": userID}).
		Where(sq.Gt{"expired_at": time.Now()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build get active session by user id query: %w", err)
	}

	rows, err := r.db.Pool.Query(ctx, qry, args...)
	if err != nil {
		return nil, fmt.Errorf("execute get active session by user id query: %w", err)
	}
	defer rows.Close()

	var sessions []*model.AuthenticationSession
	for rows.Next() {
		var session model.AuthenticationSession
		err := rows.Scan(&session.ID, &session.UserID, &session.RefreshToken, &session.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("scan active session by user id query: %w", err)
		}
		sessions = append(sessions, &session)
	}
	return sessions, nil
}

// GetActiveSessionByRefreshToken gets an authentication session by its refresh token.
func (r *authenticationSessionRepository) GetActiveSessionByRefreshToken(ctx context.Context, refreshToken string) (*model.AuthenticationSession, error) {
	qry, args, err := sq.Select("id", "user_id", "refresh_token", "expired_at").
		From(model.AuthenticationSessionTableName).
		Where(sq.Eq{"refresh_token": refreshToken}).
		Where(sq.Gt{"expired_at": time.Now()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build get active session by refresh token query: %w", err)
	}

	var session model.AuthenticationSession
	err = r.db.Pool.QueryRow(ctx, qry, args...).Scan(&session.ID, &session.UserID, &session.RefreshToken, &session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("execute get active session by refresh token query: %w", err)
	}
	return &session, nil
}

// RevokeSession revokes an authentication session.
func (r *authenticationSessionRepository) RevokeSession(ctx context.Context, sessionID uint64) error {
	qry, args, err := sq.Update(model.AuthenticationSessionTableName).
		Set("revoked_at", time.Now()).
		Where(sq.Eq{"id": sessionID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("build revoke session query: %w", err)
	}

	_, err = r.db.Pool.Exec(ctx, qry, args...)
	if err != nil {
		return fmt.Errorf("execute revoke session query: %w", err)
	}
	return nil
}
