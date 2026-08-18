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
	List(ctx context.Context, page, pageSize int) (*model.PaginatedUsers, error)

	// GetByEmail retrieves a user by their email address.
	GetByEmail(ctx context.Context, email string) (*model.User, error)

	// GetByID retrieves a user by their ID.
	GetByID(ctx context.Context, id uint64) (*model.User, error)

	// UpdateById updates an existing user's information in the database.
	UpdateById(ctx context.Context, id uint64, user *model.User) error

	// DeleteUser removes a user from the database by their ID.
	DeleteByID(ctx context.Context, id uint64) error
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
		Suffix("RETURNING id, created_at").
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

func (r *userRepository) List(ctx context.Context, page, pageSize int) (*model.PaginatedUsers, error) {
	// Get Total Count of Users
	countQuery, countArgs, err := sq.Select("COUNT(*)").
		From(model.USER_TABLE_NAME).
		Where(sq.Eq{"deleted_at": nil}). // Exclude soft-deleted users
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build count users query : %w", err)
	}

	var totalCount int
	err = r.db.Pool.QueryRow(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("execute count users query : %w", err)
	}

	offset := (page - 1) * pageSize
	query, args, err := sq.Select(
		"id",
		"first_name",
		"last_name",
		"email",
		"dial_code",
		"phone",
		"is_active",
		"created_at",
		"updated_at").
		From(model.USER_TABLE_NAME).
		Where(sq.Eq{"deleted_at": nil}). // Exclude soft-deleted users
		OrderBy("id ASC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build list users query : %w", err)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("execute list users query : %w", err)
	}
	defer rows.Close()

	var paginatedUsers = &model.PaginatedUsers{
		Users:      []*model.User{},
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}

	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.DialCode,
			&user.Phone,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("scan user row : %w", err)
		}

		paginatedUsers.Users = append(paginatedUsers.Users, &user)
	}

	return paginatedUsers, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	qry, args, err := sq.Select(
		"id",
		"first_name",
		"last_name",
		"email",
		"dial_code",
		"phone",
		"password_hash",
		"is_active",
		"created_at",
		"updated_at").
		From(model.USER_TABLE_NAME).
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build get user by email query : %w", err)
	}

	var user model.User

	row := r.db.Pool.QueryRow(ctx, qry, args...)
	err = row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.DialCode,
		&user.Phone,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("execute get user by email query : %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	qry, args, err := sq.Select(
		"id",
		"first_name",
		"last_name",
		"email",
		"dial_code",
		"phone",
		"password_hash",
		"is_active",
		"created_at",
		"updated_at").
		From(model.USER_TABLE_NAME).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build get user by id query : %w", err)
	}

	var user model.User

	row := r.db.Pool.QueryRow(ctx, qry, args...)
	err = row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.DialCode,
		&user.Phone,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("execute get user by email query : %w", err)
	}

	return &user, nil
}

func (r *userRepository) UpdateById(ctx context.Context, id uint64, user *model.User) error {
	userUpdateQry := sq.Update(model.USER_TABLE_NAME)
	if user.FirstName != "" {
		userUpdateQry = userUpdateQry.Set("first_name", user.FirstName)
	}
	if user.LastName != "" {
		userUpdateQry = userUpdateQry.Set("last_name", user.LastName)
	}
	if user.Email != "" {
		userUpdateQry = userUpdateQry.Set("email", user.Email)
	}
	if user.DialCode != "" {
		userUpdateQry = userUpdateQry.Set("dial_code", user.DialCode)
	}
	if user.Phone != "" {
		userUpdateQry = userUpdateQry.Set("phone", user.Phone)
	}
	if user.PasswordHash != "" {
		userUpdateQry = userUpdateQry.Set("password_hash", user.PasswordHash)
	}
	if user.IsActive != false {
		userUpdateQry = userUpdateQry.Set("is_active", user.IsActive)
	}

	userUpdateQry = userUpdateQry.Where(sq.Eq{"id": id}).Suffix("RETURNING updated_at").PlaceholderFormat(sq.Dollar)
	qry, args, err := userUpdateQry.ToSql()
	if err != nil {
		return fmt.Errorf("build update user query : %w", err)
	}

	err = r.db.Pool.QueryRow(ctx, qry, args...).Scan(&user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("execute update user query : %w", err)
	}
	return nil
}

func (r *userRepository) DeleteByID(ctx context.Context, id uint64) error {
	qry, args, err := sq.Delete(model.USER_TABLE_NAME).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("build delete user query : %w", err)
	}

	err = r.db.Pool.QueryRow(ctx, qry, args...).Scan()
	if err != nil {
		return fmt.Errorf("execute delete user query : %w", err)
	}
	return nil
}
