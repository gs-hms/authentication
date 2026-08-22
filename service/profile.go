package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/supermarios-hotel-management-system/authentication/dto"
	"github.com/supermarios-hotel-management-system/authentication/model"
	"github.com/supermarios-hotel-management-system/authentication/repository"
	"golang.org/x/crypto/bcrypt"
)

// ProfileService defines the interface for profile operations.
type ProfileService interface {
	GetProfile(ctx context.Context, userID uint64) (*dto.ProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uint64, req *dto.UpdateProfileRequest) error
	ChangePassword(ctx context.Context, userID uint64, req *dto.ChangePasswordRequest) error
}

type profileService struct {
	userRepo repository.UserRepository
}

// NewProfileService creates a new instance of ProfileService.
func NewProfileService(userRepo repository.UserRepository) ProfileService {
	return &profileService{
		userRepo: userRepo,
	}
}

func (s *profileService) GetProfile(ctx context.Context, userID uint64) (*dto.ProfileResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &dto.ProfileResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		DialCode:  user.DialCode,
		Phone:     user.Phone,
	}, nil
}

func (s *profileService) UpdateProfile(ctx context.Context, userID uint64, req *dto.UpdateProfileRequest) error {
	if !IsValidPhone(req.DialCode, req.Phone) {
		return ErrInvalidPhone
	}

	user := &model.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		DialCode:  req.DialCode,
		Phone:     req.Phone,
	}

	return s.userRepo.UpdateByID(ctx, userID, user)
}

func (s *profileService) ChangePassword(ctx context.Context, userID uint64, req *dto.ChangePasswordRequest) error {
	if strings.TrimSpace(req.CurrentPassword) == "" || strings.TrimSpace(req.NewPassword) == "" {
		return fmt.Errorf("passwords cannot be empty")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return ErrIncorrectPassword
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash new password: %w", err)
	}

	updateUser := &model.User{PasswordHash: string(newPasswordHash)}
	return s.userRepo.UpdateByID(ctx, userID, updateUser)
}
