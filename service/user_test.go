package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/supermarios-hotel-management-system/authentication/dto"
	mocks "github.com/supermarios-hotel-management-system/authentication/mocks/github.com/supermarios-hotel-management-system/authentication/repository"
	"github.com/supermarios-hotel-management-system/authentication/model"
	"github.com/supermarios-hotel-management-system/authentication/service"
)

func newSignupRequest(email string, dialCode model.UserDialCode, phone string, password string) *dto.SignupRequest {
	return &dto.SignupRequest{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     email,
		DialCode:  dialCode,
		Phone:     phone,
		Password:  password,
	}
}

type signupTestCase struct {
	name        string
	request     *dto.SignupRequest
	setupMock   func(repo *mocks.MockUserRepository)
	expectedErr error
}

type loginTestCase struct {
	name         string
	request      *dto.LoginRequest
	setupMock    func(repo *mocks.MockUserRepository, authSessionRepo *mocks.MockAuthenticationSessionRepository)
	expectedErr  error
	expectedResp *dto.LoginResponse
}

func TestSignUp(t *testing.T) {
	tests := []signupTestCase{
		{
			name:    "successful signup",
			request: newSignupRequest("jondoe@example.com", model.IndiaDialCode, "9876543210", "password"),

			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On(
					"GetByEmail",
					context.Background(),
					"jondoe@example.com",
				).Return(nil, nil)

				repo.On(
					"GetByPhone",
					context.Background(),
					model.IndiaDialCode,
					"9876543210",
				).Return(nil, nil)

				repo.On(
					"CreateUser",
					context.Background(),
					mock.AnythingOfType("*model.User"),
				).Return(nil)
			},

			expectedErr: nil,
		},
		{
			name:        "invalid email",
			request:     newSignupRequest("invalid-email", "91", "9876543210", "password"),
			expectedErr: service.ErrInvalidEmail,
		},
		{
			name:    "existing email",
			request: newSignupRequest("jondoe@mailinator.com", "91", "9876543210", "password"),

			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On(
					"GetByEmail",
					mock.Anything,
					"jondoe@mailinator.com",
				).Return(&model.User{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "jondoe@mailinator.com",
					Phone:     "1234567890",
				}, nil)
			},

			expectedErr: service.ErrUserWithEmailExists,
		},
		{
			name:        "invalid phone",
			request:     newSignupRequest("jondoe3@example.com", "91", "123", "password"),
			expectedErr: service.ErrInvalidPhone,
		},
		{
			name:    "existing phone",
			request: newSignupRequest("jondoe4@example.com", "91", "9876543210", "password"),
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On(
					"GetByEmail",
					mock.Anything,
					"jondoe4@example.com",
				).Return(nil, nil)

				repo.On(
					"GetByPhone",
					mock.Anything,
					model.UserDialCode("91"),
					"9876543210",
				).Return(&model.User{
					FirstName: "John",
					LastName:  "Doe",
					Email:     "existing@example.com",
					Phone:     "9876543210",
				}, nil)
			},
			expectedErr: service.ErrUserWithPhoneExists,
		},
		{
			name:        "invalid phone number length",
			request:     newSignupRequest("jondoe5@example.com", "91", "123", "password"),
			expectedErr: service.ErrInvalidPhone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mocks.NewMockUserRepository(t)
			authSessionRepo := mocks.NewMockAuthenticationSessionRepository(t)

			if tt.setupMock != nil {
				tt.setupMock(userRepo)
			}

			svc := service.NewUserService(userRepo, authSessionRepo)

			_, err := svc.Signup(context.Background(), tt.request)

			assert.ErrorIs(t, err, tt.expectedErr)
			userRepo.AssertExpectations(t)
			authSessionRepo.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	os.Setenv("JWT_SECRET_STRING", "secret")
	defer os.Unsetenv("JWT_SECRET_STRING")

	tests := []loginTestCase{
		{
			name: "successful login",
			request: &dto.LoginRequest{
				Email:    "jondoe@example.com",
				Password: "password",
			},
			setupMock: func(repo *mocks.MockUserRepository, authSessionRepo *mocks.MockAuthenticationSessionRepository) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
				repo.On(
					"GetByEmail",
					mock.Anything,
					"jondoe@example.com",
				).Return(&model.User{
					ID:           1,
					FirstName:    "John",
					LastName:     "Doe",
					Email:        "jondoe@example.com",
					Phone:        "1234567890",
					PasswordHash: string(hash),
					IsActive:     true,
				}, nil)
				authSessionRepo.On(
					"CreateSession",
					mock.Anything,
					uint64(1),
					mock.AnythingOfType("string"),
					mock.Anything,
				).Return(nil)
			},
			expectedErr: nil,
			expectedResp: &dto.LoginResponse{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				UserID:       1,
				Username:     "Jon Doe",
				Email:        "jondoe@example.com",
				Phone:        "1234567890",
			},
		},
	}

	for _, tt := range tests {
		userRepo := mocks.NewMockUserRepository(t)
		authSessionRepo := mocks.NewMockAuthenticationSessionRepository(t)

		if tt.setupMock != nil {
			tt.setupMock(userRepo, authSessionRepo)
		}

		svc := service.NewUserService(userRepo, authSessionRepo)

		_, err := svc.Login(context.Background(), tt.request)

		assert.ErrorIs(t, err, tt.expectedErr)
		userRepo.AssertExpectations(t)
		authSessionRepo.AssertExpectations(t)
	}
}
