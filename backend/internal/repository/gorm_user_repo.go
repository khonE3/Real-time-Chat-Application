package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/model"
	"github.com/khonE3/chat-backend/pkg/database"
	"gorm.io/gorm"
)

// GormUserRepository implements user repository with GORM
type GormUserRepository struct {
	db *database.GormDB
}

// NewGormUserRepository creates a new GORM user repository
func NewGormUserRepository(db *database.GormDB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// Create creates a new user
func (r *GormUserRepository) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	user := &model.User{
		Username:    req.Username,
		DisplayName: req.DisplayName,
	}

	if req.AvatarURL != "" {
		user.AvatarURL = &req.AvatarURL
	}

	if err := r.db.DB.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (r *GormUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.DB.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *GormUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.DB.WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *GormUserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.DB.WithContext(ctx).Save(user).Error
}

// GetOrCreate gets an existing user or creates a new one
func (r *GormUserRepository) GetOrCreate(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	var user model.User
	err := r.db.DB.WithContext(ctx).Where("username = ?", req.Username).First(&user).Error

	if err == nil {
		return &user, nil
	}

	if err == gorm.ErrRecordNotFound {
		return r.Create(ctx, req)
	}

	return nil, err
}

// List retrieves all users
func (r *GormUserRepository) List(ctx context.Context, limit, offset int) ([]model.User, error) {
	var users []model.User
	err := r.db.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// Search searches users by username or display name
func (r *GormUserRepository) Search(ctx context.Context, query string, limit int) ([]model.User, error) {
	var users []model.User
	search := "%" + query + "%"
	err := r.db.DB.WithContext(ctx).
		Where("username ILIKE ? OR display_name ILIKE ?", search, search).
		Limit(limit).
		Find(&users).Error
	return users, err
}
