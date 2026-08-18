package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/supermarios-hotel-management-system/authentication/database"
	"github.com/supermarios-hotel-management-system/authentication/model"
)

type UserRepository interface {
	// CreateUser creates a new user in the database.
	CreateUser(ctx context.Context, user *model.User) error

	// List retrieves all users from the database.
	List(ctx context.Context) ([]*model.User, error)

	// GetByEmail retrieves a user by their email address.
	GetByEmail(ctx context.Context, email string) (*model.User, error)

	// GetByID retrieves a user by their ID.
	GetByID(ctx context.Context, id uint64) (*model.User, error)

	// UpdateUser updates an existing user's information in the database.
	UpdateUser(ctx context.Context, user *model.User) error

	// DeleteUser removes a user from the database by their ID.
	DeleteUser(ctx context.Context, id uint64) error
}

type userRepository struct {
	db *database.Postgres
}

func NewUserRepository(db *database.Postgres) UserRepository {
	// Return an implementation of UserRepository, e.g., a struct that interacts with the database.
	// This is a placeholder; you would replace it with your actual implementation.
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	query, args, err := sq.Insert(model.USER_TABLE_NAME).
		Columns("first_name", "last_name", "email", "dial_code", "phone", "password_hash", "is_active").
		Values(user.FirstName, user.LastName, user.Email, user.DialCode, user.Phone, user.PasswordHash, user.IsActive).
		PlaceholderFormat(sq.Dollar).ToSql()
	
	if err != nil {
		return fmt.Errorf("build create user query : %w", err)
	}

	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("execute create user query : %w", err)	
	}
	return nil
}

func (r *userRepository) List(ctx context.Context) ([]*model.User, error) {
	// Implement the logic to list all users from the database.
	return nil, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	// Implement the logic to retrieve a user by email from the database.
	return nil, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	// Implement the logic to retrieve a user by ID from the database.
	return nil, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *model.User) error {
	// Implement the logic to update a user in the database.
	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id uint64) error {
	// Implement the logic to delete a user from the database.
	return nil
}
