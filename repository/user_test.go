package repository_test

import (
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"

	"github.com/supermarios-hotel-management-system/authentication/database"
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
