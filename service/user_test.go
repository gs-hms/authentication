package service_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/supermarios-hotel-management-system/authentication/dto"
	mocks "github.com/supermarios-hotel-management-system/authentication/mocks/github.com/supermarios-hotel-management-system/authentication/repository"
	"github.com/supermarios-hotel-management-system/authentication/model"
	"github.com/supermarios-hotel-management-system/authentication/service"
)

func newSignupRequest(email string) *dto.SignupRequest {
	return &dto.SignupRequest{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     email,
		DialCode:  "91",
		Phone:     "9876543210",
		Password:  gofakeit.Password(true, true, true, true, true, 10),
	}
}

type signupTestCase struct {
	name        string
	request     *dto.SignupRequest
	setupMock   func(repo *mocks.MockUserRepository)
	expectedErr error
}

func TestSignUp(t *testing.T) {
	tests := []signupTestCase{
		{
			name:    "successful signup",
			request: newSignupRequest("jondoe@example.com"),

			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On(
					"GetByEmail",
					context.Background(),
					"jondoe@example.com",
				).Return(nil, nil)

				repo.On(
					"GetByPhone",
					context.Background(),
					model.UserDialCode("91"),
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
			request:     newSignupRequest("invalid-email"),
			expectedErr: service.ErrInvalidEmail,
		},
		{
			name:    "existing email",
			request: newSignupRequest("jondoe@mailinator.com"),

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
			name: "invalid phone",
			request: func() *dto.SignupRequest {
				req := newSignupRequest("jondoe3@example.com")
				req.Phone = "123" // Invalid length
				return req
			}(),
			expectedErr: service.ErrInvalidPhone,
		},
		{
			name:    "existing phone",
			request: newSignupRequest("jondoe4@example.com"),
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockUserRepository(t)

			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			svc := service.NewUserService(repo)

			_, err := svc.Signup(context.Background(), tt.request)

			assert.ErrorIs(t, err, tt.expectedErr)
			repo.AssertExpectations(t)
		})
	}
}
