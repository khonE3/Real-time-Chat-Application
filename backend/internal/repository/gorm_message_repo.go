package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/model"
	"github.com/khonE3/chat-backend/pkg/database"
)

// GormMessageRepository implements message repository with GORM
type GormMessageRepository struct {
	db *database.GormDB
}

// NewGormMessageRepository creates a new GORM message repository
func NewGormMessageRepository(db *database.GormDB) *GormMessageRepository {
	return &GormMessageRepository{db: db}
}

// Create creates a new message
func (r *GormMessageRepository) Create(ctx context.Context, msg *model.Message) error {
	return r.db.DB.WithContext(ctx).Create(msg).Error
}

// GetByID retrieves a message by ID
func (r *GormMessageRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Message, error) {
	var msg model.Message
	if err := r.db.DB.WithContext(ctx).Preload("Files").First(&msg, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &msg, nil
}

// GetByRoom retrieves messages for a room with pagination
func (r *GormMessageRepository) GetByRoom(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]model.MessageWithUser, error) {
	var messages []model.Message

	query := r.db.DB.WithContext(ctx).
		Preload("Files").
		Where("room_id = ?", roomID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset)

	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}

	// Get user information
	result := make([]model.MessageWithUser, len(messages))
	for i, msg := range messages {
		result[i] = model.MessageWithUser{Message: msg}

		if msg.UserID != nil {
			var user model.User
			if err := r.db.DB.First(&user, "id = ?", *msg.UserID).Error; err == nil {
				result[i].Username = user.Username
				result[i].DisplayName = user.DisplayName
				result[i].AvatarURL = user.AvatarURL
			}
		}
	}

	return result, nil
}

// GetByRoomBefore retrieves messages before a specific message ID (for pagination)
func (r *GormMessageRepository) GetByRoomBefore(ctx context.Context, roomID, beforeID uuid.UUID, limit int) ([]model.MessageWithUser, error) {
	var beforeMsg model.Message
	if err := r.db.DB.First(&beforeMsg, "id = ?", beforeID).Error; err != nil {
		return nil, err
	}

	var messages []model.Message

	query := r.db.DB.WithContext(ctx).
		Preload("Files").
		Where("room_id = ? AND created_at < ?", roomID, beforeMsg.CreatedAt).
		Order("created_at DESC").
		Limit(limit)

	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}

	// Get user information
	result := make([]model.MessageWithUser, len(messages))
	for i, msg := range messages {
		result[i] = model.MessageWithUser{Message: msg}

		if msg.UserID != nil {
			var user model.User
			if err := r.db.DB.First(&user, "id = ?", *msg.UserID).Error; err == nil {
				result[i].Username = user.Username
				result[i].DisplayName = user.DisplayName
				result[i].AvatarURL = user.AvatarURL
			}
		}
	}

	return result, nil
}

// Search searches messages by content
func (r *GormMessageRepository) Search(ctx context.Context, roomID uuid.UUID, query string, limit int) ([]model.MessageWithUser, error) {
	var messages []model.Message
	search := "%" + query + "%"

	dbQuery := r.db.DB.WithContext(ctx).
		Preload("Files").
		Where("room_id = ? AND content ILIKE ?", roomID, search).
		Order("created_at DESC").
		Limit(limit)

	if err := dbQuery.Find(&messages).Error; err != nil {
		return nil, err
	}

	// Get user information
	result := make([]model.MessageWithUser, len(messages))
	for i, msg := range messages {
		result[i] = model.MessageWithUser{Message: msg}

		if msg.UserID != nil {
			var user model.User
			if err := r.db.DB.First(&user, "id = ?", *msg.UserID).Error; err == nil {
				result[i].Username = user.Username
				result[i].DisplayName = user.DisplayName
				result[i].AvatarURL = user.AvatarURL
			}
		}
	}

	return result, nil
}

// Delete soft-deletes a message
func (r *GormMessageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.DB.WithContext(ctx).Delete(&model.Message{}, "id = ?", id).Error
}
