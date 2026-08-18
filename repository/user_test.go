package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"

	"github.com/supermarios-hotel-management-system/authentication/database"
	"github.com/supermarios-hotel-management-system/authentication/model"
	"github.com/supermarios-hotel-management-system/authentication/repository"
)

func setupRepository(t *testing.T) (repository.UserRepository, pgxmock.PgxPoolIface) {
	t.Helper()
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)

	db := &database.Postgres{
		Pool: mock,
	}

	return repository.NewUserRepository(db), mock
}

func TestCreateUser(t *testing.T) {
	repo, mock := setupRepository(t)
	ctx := context.Background()
	createdAt := time.Now()

	user := &model.User{
		FirstName:    gofakeit.FirstName(),
		LastName:     gofakeit.LastName(),
		Email:        gofakeit.Email(),
		DialCode:     "91",
		Phone:        gofakeit.PhoneFormatted(),
		PasswordHash: gofakeit.Password(true, true, true, true, true, 10),
		IsActive:     true,
	}
	mock.ExpectQuery("INSERT INTO users").
		WithArgs(
			user.FirstName,
			user.LastName,
			user.Email,
			user.DialCode,
			user.Phone,
			user.PasswordHash,
			user.IsActive,
		).
		WillReturnRows(
			pgxmock.NewRows([]string{
				"id",
				"created_at",
			}).AddRow(
				uint64(1),
				createdAt,
			),
		)

	err := repo.CreateUser(ctx, user)
	require.NoError(t, err)
	require.Equal(t, uint64(1), user.ID)
	require.Equal(t, createdAt, user.CreatedAt)
	require.NoError(t, mock.ExpectationsWereMet())

}

func TestCreateuserError(t *testing.T) {
	repo, mock := setupRepository(t)
	ctx := context.Background()

	user := &model.User{
		FirstName:    gofakeit.FirstName(),
		LastName:     gofakeit.LastName(),
		Email:        gofakeit.Email(),
		DialCode:     "91",
		Phone:        gofakeit.PhoneFormatted(),
		PasswordHash: gofakeit.Password(true, true, true, true, true, 10),
		IsActive:     gofakeit.Bool(),
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(
			user.FirstName,
			user.LastName,
			user.Email,
			user.DialCode,
			user.Phone,
			user.PasswordHash,
			user.IsActive,
		).
		WillReturnError(errors.New("database error"))

	err := repo.CreateUser(ctx, user)
	require.Error(t, err)
	require.Contains(t, err.Error(), "execute create user query")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestListUsers(t *testing.T) {
	repo, mock := setupRepository(t)

	ctx := context.Background()
	createdAt := time.Now()
	updatedAt := time.Now()

	// mock expectations
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users WHERE deleted_at IS NULL`).
		WillReturnRows(
			pgxmock.NewRows([]string{
				"count",
			}).AddRow(5),
		)

	mock.ExpectQuery(`SELECT id, first_name, last_name, email, dial_code, phone, is_active, created_at, updated_at FROM users WHERE deleted_at IS NULL ORDER BY id ASC LIMIT 2 OFFSET 0`).
		WillReturnRows(
			pgxmock.NewRows([]string{
				"id",
				"first_name",
				"last_name",
				"email",
				"dial_code",
				"phone",
				"is_active",
				"created_at",
				"updated_at",
			}).AddRow(
				uint64(1),
				gofakeit.FirstName(),
				gofakeit.LastName(),
				gofakeit.Email(),
				"91",
				gofakeit.PhoneFormatted(),
				gofakeit.Bool(),
				createdAt,
				updatedAt,
			).AddRow(
				uint64(2),
				gofakeit.FirstName(),
				gofakeit.LastName(),
				gofakeit.Email(),
				"91",
				gofakeit.PhoneFormatted(),
				gofakeit.Bool(),
				createdAt,
				updatedAt,
			),
		)

	users, err := repo.List(ctx, 1, 2)
	require.NoError(t, err)
	require.NotNil(t, users)
	require.Equal(t, 5, users.TotalCount)
	require.Equal(t, 2, len(users.Users))
	require.Equal(t, uint64(1), users.Users[0].ID)
	require.Equal(t, createdAt, users.Users[0].CreatedAt)
	require.Equal(t, updatedAt, users.Users[0].UpdatedAt)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestListUsersError(t *testing.T) {
	repo, mock := setupRepository(t)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users WHERE deleted_at IS NULL`).
		WillReturnRows(
			pgxmock.NewRows([]string{"count"}).
				AddRow(5),
		)

	mock.ExpectQuery(`SELECT id, first_name, last_name, email, dial_code, phone, is_active, created_at, updated_at FROM users WHERE deleted_at IS NULL ORDER BY id ASC LIMIT 10 OFFSET 0`).
		WillReturnError(errors.New("database error"))

	result, err := repo.List(context.Background(), 1, 10)

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "execute list users query")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByEmail(t *testing.T) {
	repo, mock := setupRepository(t)
	ctx := context.Background()

	createdAt := time.Now()
	updatedAt := time.Now()
	email := gofakeit.Email()

	mock.ExpectQuery(`SELECT id, first_name, last_name, email, dial_code, phone, password_hash, is_active, created_at, updated_at FROM users WHERE deleted_at IS NULL AND email = \$1`).
		WithArgs(email).
		WillReturnRows(
			pgxmock.NewRows([]string{
				"id",
				"first_name",
				"last_name",
				"email",
				"dial_code",
				"phone",
				"password_hash",
				"is_active",
				"created_at",
				"updated_at",
			}).AddRow(
				uint64(1),
				gofakeit.FirstName(),
				gofakeit.LastName(),
				email,
				"91",
				gofakeit.PhoneFormatted(),
				gofakeit.Password(true, true, true, true, true, 10),
				gofakeit.Bool(),
				createdAt,
				updatedAt,
			),
		)

	user, err := repo.GetByEmail(ctx, email)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, email, user.Email)
	require.Equal(t, uint64(1), user.ID)
	require.Equal(t, createdAt, user.CreatedAt)
	require.Equal(t, updatedAt, user.UpdatedAt)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByEmailNotFound(t *testing.T) {
	repo, mock := setupRepository(t)
	ctx := context.Background()
	email := gofakeit.Email()

	mock.ExpectQuery(`SELECT id, first_name, last_name, email, dial_code, phone, password_hash, is_active, created_at, updated_at FROM users WHERE deleted_at IS NULL AND email = \$1`).
		WithArgs(email).
		WillReturnError(pgx.ErrNoRows)

	user, err := repo.GetByEmail(ctx, email)
	require.Error(t, err)
	require.Nil(t, user)
	require.ErrorContains(t, err, "no rows in result set")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID(t *testing.T) {
	repo, mock := setupRepository(t)
	ctx := context.Background()

	createdAt := time.Now()
	updatedAt := time.Now()
	id := uint64(1)

	mock.ExpectQuery(`SELECT id, first_name, last_name, email, dial_code, phone, password_hash, is_active, created_at, updated_at FROM users WHERE deleted_at IS NULL AND id = \$1`).
		WithArgs(id).
		WillReturnRows(
			pgxmock.NewRows([]string{
				"id",
				"first_name",
				"last_name",
				"email",
				"dial_code",
				"phone",
				"password_hash",
				"is_active",
				"created_at",
				"updated_at",
			}).AddRow(
				id,
				gofakeit.FirstName(),
				gofakeit.LastName(),
				gofakeit.Email(),
				"91",
				gofakeit.PhoneFormatted(),
				gofakeit.Password(true, true, true, true, true, 10),
				gofakeit.Bool(),
				createdAt,
				updatedAt,
			),
		)

	user, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, id, user.ID)
	require.Equal(t, createdAt, user.CreatedAt)
	require.Equal(t, updatedAt, user.UpdatedAt)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByIDNotFound(t *testing.T) {
	repo, mock := setupRepository(t)
	ctx := context.Background()
	id := uint64(1)

	mock.ExpectQuery(`SELECT id, first_name, last_name, email, dial_code, phone, password_hash, is_active, created_at, updated_at FROM users WHERE deleted_at IS NULL AND id = \$1`).
		WithArgs(id).
		WillReturnError(pgx.ErrNoRows)

	user, err := repo.GetByID(ctx, id)
	require.Error(t, err)
	require.Nil(t, user)
	require.ErrorContains(t, err, "no rows in result set")

	require.NoError(t, mock.ExpectationsWereMet())
}