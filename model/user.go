package model

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
	CreatedAt    string       `json:"created_at"`
	UpdatedAt    string       `json:"updated_at"`
}
