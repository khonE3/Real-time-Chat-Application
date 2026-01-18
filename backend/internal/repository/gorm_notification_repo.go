package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/model"
	"github.com/khonE3/chat-backend/pkg/database"
)

// GormNotificationRepository implements notification repository with GORM
type GormNotificationRepository struct {
	db *database.GormDB
}

// NewGormNotificationRepository creates a new GORM notification repository
func NewGormNotificationRepository(db *database.GormDB) *GormNotificationRepository {
	return &GormNotificationRepository{db: db}
}

// Create creates a new notification
func (r *GormNotificationRepository) Create(ctx context.Context, notification *model.Notification) error {
	return r.db.DB.WithContext(ctx).Create(notification).Error
}

// GetByID retrieves a notification by ID
func (r *GormNotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Notification, error) {
	var notification model.Notification
	if err := r.db.DB.WithContext(ctx).First(&notification, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

// GetByUserID retrieves notifications for a user
func (r *GormNotificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.NotificationWithUser, error) {
	var notifications []model.Notification
	err := r.db.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error
	if err != nil {
		return nil, err
	}

	// Enrich with from user info
	result := make([]model.NotificationWithUser, len(notifications))
	for i, n := range notifications {
		result[i] = model.NotificationWithUser{Notification: n}
		if n.FromUser != nil {
			var user model.User
			if err := r.db.DB.First(&user, "id = ?", *n.FromUser).Error; err == nil {
				result[i].FromUsername = user.Username
				result[i].FromDisplayName = user.DisplayName
				result[i].FromAvatarURL = user.AvatarURL
			}
		}
	}

	return result, nil
}

// GetUnreadByUserID retrieves unread notifications for a user
func (r *GormNotificationRepository) GetUnreadByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]model.NotificationWithUser, error) {
	var notifications []model.Notification
	err := r.db.DB.WithContext(ctx).
		Where("user_id = ? AND is_read = ?", userID, false).
		Order("created_at DESC").
		Limit(limit).
		Find(&notifications).Error
	if err != nil {
		return nil, err
	}

	// Enrich with from user info
	result := make([]model.NotificationWithUser, len(notifications))
	for i, n := range notifications {
		result[i] = model.NotificationWithUser{Notification: n}
		if n.FromUser != nil {
			var user model.User
			if err := r.db.DB.First(&user, "id = ?", *n.FromUser).Error; err == nil {
				result[i].FromUsername = user.Username
				result[i].FromDisplayName = user.DisplayName
				result[i].FromAvatarURL = user.AvatarURL
			}
		}
	}

	return result, nil
}

// GetUnreadCount gets count of unread notifications for a user
func (r *GormNotificationRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.DB.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

// MarkAsRead marks a notification as read
func (r *GormNotificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.DB.WithContext(ctx).
		Model(&model.Notification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// MarkAllAsRead marks all notifications as read for a user
func (r *GormNotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return r.db.DB.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// Delete deletes a notification
func (r *GormNotificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.DB.WithContext(ctx).Delete(&model.Notification{}, "id = ?", id).Error
}

// DeleteAllByUserID deletes all notifications for a user
func (r *GormNotificationRepository) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.Notification{}).Error
}

// CreatePushSubscription creates a push subscription
func (r *GormNotificationRepository) CreatePushSubscription(ctx context.Context, sub *model.PushSubscription) error {
	// First try to find existing subscription with same endpoint
	var existing model.PushSubscription
	err := r.db.DB.WithContext(ctx).
		Where("user_id = ? AND endpoint = ?", sub.UserID, sub.Endpoint).
		First(&existing).Error

	if err == nil {
		// Update existing subscription
		existing.P256DH = sub.P256DH
		existing.Auth = sub.Auth
		existing.UserAgent = sub.UserAgent
		return r.db.DB.Save(&existing).Error
	}

	// Create new subscription
	return r.db.DB.WithContext(ctx).Create(sub).Error
}

// GetPushSubscriptions gets push subscriptions for a user
func (r *GormNotificationRepository) GetPushSubscriptions(ctx context.Context, userID uuid.UUID) ([]model.PushSubscription, error) {
	var subs []model.PushSubscription
	err := r.db.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&subs).Error
	return subs, err
}

// DeletePushSubscription deletes a push subscription by endpoint
func (r *GormNotificationRepository) DeletePushSubscription(ctx context.Context, userID uuid.UUID, endpoint string) error {
	return r.db.DB.WithContext(ctx).
		Where("user_id = ? AND endpoint = ?", userID, endpoint).
		Delete(&model.PushSubscription{}).Error
}
