package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/model"
	"github.com/khonE3/chat-backend/pkg/database"
)

type UserRepository struct {
	db *database.Postgres
}

func NewUserRepository(db *database.Postgres) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	user := &model.User{
		ID:          uuid.New(),
		Username:    req.Username,
		DisplayName: req.DisplayName,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if req.AvatarURL != "" {
		user.AvatarURL = &req.AvatarURL
	}

	query := `
		INSERT INTO users (id, username, display_name, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, username, display_name, avatar_url, created_at, updated_at
	`

	err := r.db.Pool.QueryRow(ctx, query,
		user.ID, user.Username, user.DisplayName, user.AvatarURL, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID, &user.Username, &user.DisplayName, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user := &model.User{}

	query := `
		SELECT id, username, display_name, avatar_url, created_at, updated_at
		FROM users WHERE id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.DisplayName, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}

	query := `
		SELECT id, username, display_name, avatar_url, created_at, updated_at
		FROM users WHERE username = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.DisplayName, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users 
		SET display_name = $2, avatar_url = $3, updated_at = $4
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, user.ID, user.DisplayName, user.AvatarURL, time.Now())
	return err
}

func (r *UserRepository) GetOrCreate(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	// Try to get existing user
	user, err := r.GetByUsername(ctx, req.Username)
	if err == nil {
		return user, nil
	}

	// Create new user if not found
	return r.Create(ctx, req)
}
