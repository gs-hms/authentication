// Package service provides the business logic for the authentication service.
package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/supermarios-hotel-management-system/authentication/dto"
	"github.com/supermarios-hotel-management-system/authentication/model"
	"github.com/supermarios-hotel-management-system/authentication/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines the interface for user service.
type UserService interface {
	// Signup creates a new user.
	Signup(ctx context.Context, req *dto.SignupRequest) (*model.User, error)
}

type userService struct {
	userRepo        repository.UserRepository
	authSessionRepo repository.AuthenticationSessionRepository
}

// NewUserService creates a new instance of the UserService.
func NewUserService(userRepo repository.UserRepository, authSessionRepo repository.AuthenticationSessionRepository) UserService {
	return &userService{
		userRepo:        userRepo,
		authSessionRepo: authSessionRepo,
	}
}

func (s *userService) Signup(ctx context.Context, req *dto.SignupRequest) (*model.User, error) {
	if err := s.validateContactDetails(req.Email, req.DialCode, req.Phone); err != nil {
		return nil, err
	}

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

	// Check if user with the same phone number exists
	existingUser, err = s.userRepo.GetByPhone(ctx, req.DialCode, req.Phone)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, ErrUserWithPhoneExists
	}

	if err := s.userRepo.CreateUser(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokens(ctx, user)
}

// Private method
func (s *userService) validateContactDetails(email string, dialCode model.UserDialCode, phone string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidEmail
	}

	if !IsValidPhone(dialCode, phone) {
		return ErrInvalidPhone
	}

	return nil
}

// Private method
func (s *userService) generateTokens(ctx context.Context, user *model.User) (*dto.LoginResponse, error) {
	if user == nil && user.ID == 0 && user.Email == "" {
		return nil, errors.New("generateTokens: invalid request")
	}

	secretString := os.Getenv("JWT_SECRET_STRING")
	if secretString == "" {
		return nil, errors.New("generateTokens: JWT_SECRET_STRING not set")
	}

	now := time.Now()
	claims := dto.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    "authentication-service",
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(secretString))
	if err != nil {
		return nil, fmt.Errorf("generateTokens: failed to sign access token: %w", err)
	}

	authSession, err := model.NewAuthenticationSession(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generateTokens: failed to create authentication session: %w", err)
	}

	if err := s.authSessionRepo.CreateSession(ctx, user.ID, authSession.RefreshToken, authSession.ExpiresAt); err != nil {
		return nil, fmt.Errorf("generateTokens: failed to create authentication session: %w", err)
	}

	resp := &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: authSession.RefreshToken,
		UserID:       user.ID,
		Username:     fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Email:        user.Email,
		Phone:        fmt.Sprintf("%s%s", user.DialCode, user.Phone),
	}
	return resp, nil
}
