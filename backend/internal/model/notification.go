package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationType defines types of notifications
type NotificationType string

const (
	NotificationTypeMessage  NotificationType = "message"
	NotificationTypeMention  NotificationType = "mention"
	NotificationTypeReaction NotificationType = "reaction"
	NotificationTypeDM       NotificationType = "dm"
	NotificationTypeSystem   NotificationType = "system"
)

// Notification represents a user notification
type Notification struct {
	ID        uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	Type      NotificationType `gorm:"type:varchar(50);not null" json:"type"`
	Title     string           `gorm:"not null;size:255" json:"title"`
	Body      string           `gorm:"type:text" json:"body"`
	Data      string           `gorm:"type:jsonb" json:"data,omitempty"`
	IsRead    bool             `gorm:"default:false;index" json:"is_read"`
	RoomID    *uuid.UUID       `gorm:"type:uuid;index" json:"room_id,omitempty"`
	MessageID *uuid.UUID       `gorm:"type:uuid" json:"message_id,omitempty"`
	FromUser  *uuid.UUID       `gorm:"type:uuid" json:"from_user,omitempty"`
	CreatedAt time.Time        `gorm:"autoCreateTime" json:"created_at"`
	ReadAt    *time.Time       `json:"read_at,omitempty"`
	DeletedAt gorm.DeletedAt   `gorm:"index" json:"-"`

	// Relationships
	User        *User    `gorm:"foreignKey:UserID" json:"-"`
	Room        *Room    `gorm:"foreignKey:RoomID" json:"-"`
	Message     *Message `gorm:"foreignKey:MessageID" json:"-"`
	FromUserObj *User    `gorm:"foreignKey:FromUser" json:"-"`
}

// TableName overrides the default table name
func (Notification) TableName() string {
	return "notifications"
}

// BeforeCreate hook to set UUID if not set
func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

// NotificationWithUser includes from user information
type NotificationWithUser struct {
	Notification
	FromUsername    string  `json:"from_username,omitempty"`
	FromDisplayName string  `json:"from_display_name,omitempty"`
	FromAvatarURL   *string `json:"from_avatar_url,omitempty"`
}

// PushSubscription stores web push subscriptions
type PushSubscription struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Endpoint  string         `gorm:"type:text;not null" json:"endpoint"`
	P256DH    string         `gorm:"type:text" json:"p256dh"`
	Auth      string         `gorm:"type:text" json:"auth"`
	UserAgent string         `gorm:"size:500" json:"user_agent,omitempty"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName overrides the default table name
func (PushSubscription) TableName() string {
	return "push_subscriptions"
}

// CreateNotificationRequest for creating notifications
type CreateNotificationRequest struct {
	Type      NotificationType `json:"type" validate:"required"`
	Title     string           `json:"title" validate:"required"`
	Body      string           `json:"body,omitempty"`
	RoomID    string           `json:"room_id,omitempty"`
	MessageID string           `json:"message_id,omitempty"`
	FromUser  string           `json:"from_user,omitempty"`
}

// SubscribePushRequest for subscribing to push notifications
type SubscribePushRequest struct {
	Endpoint  string `json:"endpoint" validate:"required"`
	P256DH    string `json:"p256dh"`
	Auth      string `json:"auth"`
	UserAgent string `json:"user_agent,omitempty"`
}

// NotificationSettings for user notification preferences
type NotificationSettings struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID            uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	EnablePush        bool      `gorm:"default:true" json:"enable_push"`
	EnableSound       bool      `gorm:"default:true" json:"enable_sound"`
	EnableDesktop     bool      `gorm:"default:true" json:"enable_desktop"`
	MuteDMs           bool      `gorm:"default:false" json:"mute_dms"`
	MuteGroupChats    bool      `gorm:"default:false" json:"mute_group_chats"`
	MuteMentionsOnly  bool      `gorm:"default:false" json:"mute_mentions_only"`
	QuietHoursEnabled bool      `gorm:"default:false" json:"quiet_hours_enabled"`
	QuietHoursStart   string    `gorm:"size:5" json:"quiet_hours_start,omitempty"`
	QuietHoursEnd     string    `gorm:"size:5" json:"quiet_hours_end,omitempty"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName overrides the default table name
func (NotificationSettings) TableName() string {
	return "notification_settings"
}
