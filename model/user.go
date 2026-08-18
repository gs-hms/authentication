package model

import "time"

type UserDialCode string

const (
	INDIA_DIAL_CODE UserDialCode = "91"

	USER_TABLE_NAME = "users"
)

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

type PaginatedUsers struct {
	Users      []*User `json:"users"`
	TotalCount int     `json:"total_count"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
}
