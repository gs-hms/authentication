package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/gs-hms/authentication/dto"
	mocks "github.com/gs-hms/authentication/mocks/github.com/gs-hms/authentication/repository"
	"github.com/gs-hms/authentication/model"
	"github.com/gs-hms/authentication/service"
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

			client, _ := redismock.NewClientMock()
			svc := service.NewUserService(userRepo, authSessionRepo, client, nil)

			_, err := svc.Signup(context.Background(), tt.request)

			assert.ErrorIs(t, err, tt.expectedErr)
			userRepo.AssertExpectations(t)
			authSessionRepo.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	t.Setenv("JWT_SECRET_STRING", "secret")

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
		{
			name: "empty email",
			request: &dto.LoginRequest{
				Email:    "   ",
				Password: "password",
			},
			expectedErr: service.ErrInvalidCredentials,
		},
		{
			name: "empty password",
			request: &dto.LoginRequest{
				Email:    "jondoe@example.com",
				Password: "   ",
			},
			expectedErr: service.ErrInvalidCredentials,
		},
		{
			name: "invalid email (user not found)",
			request: &dto.LoginRequest{
				Email:    "notfound@example.com",
				Password: "password",
			},
			setupMock: func(repo *mocks.MockUserRepository, _ *mocks.MockAuthenticationSessionRepository) {
				repo.On(
					"GetByEmail",
					mock.Anything,
					"notfound@example.com",
				).Return(nil, nil)
			},
			expectedErr: service.ErrUserNotFound,
		},
		{
			name: "incorrect credentials",
			request: &dto.LoginRequest{
				Email:    "jondoe@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(repo *mocks.MockUserRepository, _ *mocks.MockAuthenticationSessionRepository) {
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
			},
			expectedErr: service.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		userRepo := mocks.NewMockUserRepository(t)
		authSessionRepo := mocks.NewMockAuthenticationSessionRepository(t)

		if tt.setupMock != nil {
			tt.setupMock(userRepo, authSessionRepo)
		}

		client, _ := redismock.NewClientMock()
		svc := service.NewUserService(userRepo, authSessionRepo, client, nil)

		_, err := svc.Login(context.Background(), tt.request)

		assert.ErrorIs(t, err, tt.expectedErr)
		userRepo.AssertExpectations(t)
		authSessionRepo.AssertExpectations(t)
	}
}

func TestLogout(t *testing.T) {
	tests := []struct {
		name        string
		req         *dto.LogoutRequest
		userID      uint64
		jti         string
		exp         time.Time
		setupMock   func(authSessionRepo *mocks.MockAuthenticationSessionRepository, redisMock redismock.ClientMock)
		expectedErr error
	}{
		{
			name:   "successful logout",
			req:    &dto.LogoutRequest{RefreshToken: "valid-refresh-token"},
			userID: 1,
			jti:    "test-jti",
			exp:    time.Now().Add(1 * time.Hour),
			setupMock: func(repo *mocks.MockAuthenticationSessionRepository, redisMock redismock.ClientMock) {
				repo.On("GetActiveSessionByRefreshToken", context.Background(), "valid-refresh-token").
					Return(&model.AuthenticationSession{ID: 10, UserID: 1}, nil)
				repo.On("RevokeSession", context.Background(), uint64(10)).Return(nil)
				redisMock.ExpectSet("test-jti", "blacklisted", 1*time.Hour).SetVal("OK")
			},
			expectedErr: nil,
		},
		{
			name:   "unauthorized session revocation (IDOR)",
			req:    &dto.LogoutRequest{RefreshToken: "other-user-token"},
			userID: 1,
			jti:    "test-jti",
			exp:    time.Now().Add(1 * time.Hour),
			setupMock: func(repo *mocks.MockAuthenticationSessionRepository, _ redismock.ClientMock) {
				repo.On("GetActiveSessionByRefreshToken", context.Background(), "other-user-token").
					Return(&model.AuthenticationSession{ID: 10, UserID: 2}, nil)
			},
			expectedErr: errors.New("unauthorized session revocation"),
		},
		{
			name:   "get session error",
			req:    &dto.LogoutRequest{RefreshToken: "error-token"},
			userID: 1,
			jti:    "test-jti",
			exp:    time.Now().Add(1 * time.Hour),
			setupMock: func(repo *mocks.MockAuthenticationSessionRepository, _ redismock.ClientMock) {
				repo.On("GetActiveSessionByRefreshToken", context.Background(), "error-token").
					Return(nil, errors.New("db error"))
			},
			expectedErr: errors.New("failed to get active session: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mocks.NewMockUserRepository(t)
			authSessionRepo := mocks.NewMockAuthenticationSessionRepository(t)
			client, redisMock := redismock.NewClientMock()

			if tt.setupMock != nil {
				tt.setupMock(authSessionRepo, redisMock)
			}

			svc := service.NewUserService(userRepo, authSessionRepo, client, nil)

			err := svc.Logout(context.Background(), tt.userID, tt.jti, tt.exp, tt.req)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			userRepo.AssertExpectations(t)
			authSessionRepo.AssertExpectations(t)
			assert.NoError(t, redisMock.ExpectationsWereMet())
		})
	}
}
