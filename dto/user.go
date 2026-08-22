// Package dto provides data transfer objects for the authentication service.
package dto

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/supermarios-hotel-management-system/authentication/model"
)

// SignupRequest represents the request object for user signup.
type SignupRequest struct {
	FirstName string             `json:"first_name" validate:"required"`
	LastName  string             `json:"last_name" validate:"required"`
	Email     string             `json:"email" validate:"required"`
	DialCode  model.UserDialCode `json:"dial_code" validate:"required"`
	Phone     string             `json:"phone" validate:"required"`
	Password  string             `json:"password" validate:"required"`
}

// LoginRequest represents the request object for user login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the response object for user login.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       uint64 `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
}

// Claims represents the claims for the JWT token.
type Claims struct {
	UserID uint64 `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}
