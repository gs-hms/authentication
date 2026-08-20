// Package service provides the business logic for the authentication service.
package service

import (
	"context"
	"fmt"

	"github.com/supermarios-hotel-management-system/authentication/dto"
	"github.com/supermarios-hotel-management-system/authentication/model"
	"github.com/supermarios-hotel-management-system/authentication/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines the interface for user service.
type UserService interface {
	Signup(ctx context.Context, req *dto.SignupRequest) (*model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new instance of the UserService.
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Signup(ctx context.Context, req *dto.SignupRequest) (*model.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := model.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		DialCode:     req.DialCode,
		Phone:        req.Phone,
		PasswordHash: string(passwordHash),
		IsActive:     true,
	}

	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, ErrUserWithEmailExists
	}

	if err := s.userRepo.CreateUser(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
