package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/model"
	"github.com/khonE3/chat-backend/pkg/database"
)

// GormFileRepository implements file repository with GORM
type GormFileRepository struct {
	db *database.GormDB
}

// NewGormFileRepository creates a new GORM file repository
func NewGormFileRepository(db *database.GormDB) *GormFileRepository {
	return &GormFileRepository{db: db}
}

// Create creates a new file record
func (r *GormFileRepository) Create(ctx context.Context, file *model.File) error {
	return r.db.DB.WithContext(ctx).Create(file).Error
}

// GetByID retrieves a file by ID
func (r *GormFileRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.File, error) {
	var file model.File
	if err := r.db.DB.WithContext(ctx).First(&file, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// GetByMessageID retrieves files for a message
func (r *GormFileRepository) GetByMessageID(ctx context.Context, messageID uuid.UUID) ([]model.File, error) {
	var files []model.File
	err := r.db.DB.WithContext(ctx).Where("message_id = ?", messageID).Find(&files).Error
	return files, err
}

// GetByUserID retrieves files uploaded by a user
func (r *GormFileRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.File, error) {
	var files []model.File
	err := r.db.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&files).Error
	return files, err
}

// LinkToMessage links a file to a message
func (r *GormFileRepository) LinkToMessage(ctx context.Context, fileID, messageID uuid.UUID) error {
	return r.db.DB.WithContext(ctx).
		Model(&model.File{}).
		Where("id = ?", fileID).
		Update("message_id", messageID).Error
}

// Delete soft-deletes a file
func (r *GormFileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.DB.WithContext(ctx).Delete(&model.File{}, "id = ?", id).Error
}

// GetByFilename retrieves a file by filename
func (r *GormFileRepository) GetByFilename(ctx context.Context, filename string) (*model.File, error) {
	var file model.File
	if err := r.db.DB.WithContext(ctx).First(&file, "filename = ?", filename).Error; err != nil {
		return nil, err
	}
	return &file, nil
}
