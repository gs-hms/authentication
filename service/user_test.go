package service_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/supermarios-hotel-management-system/authentication/dto"
	mocks "github.com/supermarios-hotel-management-system/authentication/mocks/github.com/supermarios-hotel-management-system/authentication/repository"
	"github.com/supermarios-hotel-management-system/authentication/service"
)

func TestSignUp(t *testing.T) {
	type test struct {
		name        string
		user        *dto.SignupRequest
		setupMock   func(repo *mocks.MockUserRepository)
		expectedErr error
	}

	tests := []test{
		{
			name: "successful signup",
			user: &dto.SignupRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     "jondoe@example.com",
				DialCode:  "91",
				Phone:     gofakeit.Phone(),
				Password:  gofakeit.Password(true, true, true, true, true, 10),
			},
			setupMock: func(repo *mocks.MockUserRepository) {
				repo.On("GetByEmail", context.Background(), "jondoe@example.com").Return(nil, nil)
				repo.On("CreateUser", context.Background(), mock.AnythingOfType("*model.User")).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "invalid email",
			user: &dto.SignupRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     "invalid-email",
				DialCode:  "91",
				Phone:     gofakeit.Phone(),
				Password:  gofakeit.Password(true, true, true, true, true, 10),
			},
			expectedErr: service.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockUserRepository(t)
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}
			svc := service.NewUserService(repo)

			_, err := svc.Signup(context.Background(), tt.user)
			assert.ErrorIs(t, err, tt.expectedErr)
			repo.AssertExpectations(t)
		})
	}
}
