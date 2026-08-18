// Package model provides the data models for the authentication service.
package model

import "time"

// UserDialCode represents the dial code for a user.
type UserDialCode string

const (
	// IndiaDialCode represents the dial code for India.
	IndiaDialCode UserDialCode = "91"

	// UserTableName represents the name of the users table.
	UserTableName = "users"
)

// User represents a user in the authentication service.
type User struct {
	ID           uint64       `json:"id" `
	FirstName    string       `json:"first_name"`
	LastName     string       `json:"last_name"`
	Email        string       `json:"email"`
	DialCode     UserDialCode `json:"dial_code"`
	Phone        string       `json:"phone"`
	PasswordHash string       `json:"-"`
	IsActive     bool         `json:"is_active"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// PaginatedUsers represents a paginated list of users.
type PaginatedUsers struct {
	Users      []*User `json:"users"`
	TotalCount uint64  `json:"total_count"`
	Page       uint64  `json:"page"`
	PageSize   uint64  `json:"page_size"`
}
