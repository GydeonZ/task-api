package repository

import (
	"context"
	"time"

	"github.com/GydeonZ/task-api/internal/database"
	"github.com/GydeonZ/task-api/internal/models"
	"github.com/GydeonZ/task-api/pkg/errors"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
	INSERT INTO users (username, email, password, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := database.DB.QueryRowContext(ctx, query, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, errors.NewWithErr("DB_CREATE_USER_ERROR",
			"Failed to create user in database: "+err.Error(), 500, err)
	}

	return user, nil
}

func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
	SELECT id, username, email
	FROM users
	WHERE username=$1
	`

	user := &models.User{}
	err := database.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
	SELECT id, username, email, password
	FROM users
	WHERE email=$1
	`

	user := &models.User{}
	err := database.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}
