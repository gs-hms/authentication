// Package dto provides data transfer objects for the authentication service.
package dto

import "github.com/supermarios-hotel-management-system/authentication/model"

// SignupRequest represents the request object for user signup.
type SignupRequest struct {
	FirstName string             `json:"first_name" validate:"required"`
	LastName  string             `json:"last_name" validate:"required"`
	Email     string             `json:"email" validate:"required"`
	DialCode  model.UserDialCode `json:"dial_code" validate:"required"`
	Phone     string             `json:"phone" validate:"required"`
	Password  string             `json:"password" validate:"required"`
}
