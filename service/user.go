// Package service provides the business logic for the authentication service.
package service

import (
	"context"

	"github.com/supermarios-hotel-management-system/authentication/dto"
	"github.com/supermarios-hotel-management-system/authentication/model"
)

// UserService defines the interface for user service.
type UserService interface {
	Signup(ctx context.Context, req *dto.SignupRequest) (*model.User, error)
}
