package model

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

const (
	// AuthenticationSessionTableName is the name of the authentication sessions table.
	AuthenticationSessionTableName = "authentication_sessions"
)

// AuthenticationSession represents an authentication session.
type AuthenticationSession struct {
	ID           uint64    `json:"id"`
	UserID       uint64    `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	RevokedAt    time.Time `json:"revoked_at"`
}

// NewAuthenticationSession creates a new authentication session for a user.
func NewAuthenticationSession(userID uint64) (*AuthenticationSession, error) {
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("NewAuthenticationSession: %w", err)
	}

	if userID == 0 {
		return nil, fmt.Errorf("NewAuthenticationSession: invalid user ID")
	}

	if strings.TrimSpace(refreshToken) == "" {
		return nil, fmt.Errorf("NewAuthenticationSession: invalid refresh token")
	}

	return &AuthenticationSession{
		UserID:       userID,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().AddDate(0, 1, 0), // 1 month
	}, nil
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generateRefreshToken: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
