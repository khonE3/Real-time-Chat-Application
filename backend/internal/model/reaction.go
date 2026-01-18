package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Reaction represents an emoji reaction on a message
type Reaction struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	MessageID uuid.UUID      `gorm:"type:uuid;not null;index" json:"message_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Emoji     string         `gorm:"not null;size:20" json:"emoji"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Message *Message `gorm:"foreignKey:MessageID" json:"-"`
	User    *User    `gorm:"foreignKey:UserID" json:"-"`
}

// TableName overrides the default table name
func (Reaction) TableName() string {
	return "message_reactions"
}

// BeforeCreate hook to set UUID if not set
func (r *Reaction) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// ReactionWithUser includes user information
type ReactionWithUser struct {
	Reaction
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

// ReactionCount groups reactions by emoji
type ReactionCount struct {
	Emoji string   `json:"emoji"`
	Count int      `json:"count"`
	Users []string `json:"users"` // user IDs
}

// AddReactionRequest for adding reactions
type AddReactionRequest struct {
	MessageID string `json:"message_id" validate:"required,uuid"`
	Emoji     string `json:"emoji" validate:"required,min=1,max=20"`
}

// RemoveReactionRequest for removing reactions
type RemoveReactionRequest struct {
	MessageID string `json:"message_id" validate:"required,uuid"`
	Emoji     string `json:"emoji" validate:"required,min=1,max=20"`
}

// Common emoji reactions
var CommonReactions = []string{"👍", "❤️", "😂", "😮", "😢", "😡", "🎉", "🤔"}
