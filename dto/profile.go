package dto

import "github.com/supermarios-hotel-management-system/authentication/model"

// UpdateProfileRequest represents the request object for updating a user profile.
type UpdateProfileRequest struct {
	FirstName string             `json:"first_name" validate:"required"`
	LastName  string             `json:"last_name" validate:"required"`
	DialCode  model.UserDialCode `json:"dial_code" validate:"required"`
	Phone     string             `json:"phone" validate:"required"`
}

// ChangePasswordRequest represents the request object for changing the password.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
}

// ProfileResponse represents the response object for user profile.
type ProfileResponse struct {
	ID        uint64             `json:"id"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Email     string             `json:"email"`
	DialCode  model.UserDialCode `json:"dial_code"`
	Phone     string             `json:"phone"`
}
