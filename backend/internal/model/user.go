package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a chat user
type User struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username    string         `gorm:"uniqueIndex;not null;size:50" json:"username"`
	DisplayName string         `gorm:"not null;size:100" json:"display_name"`
	AvatarURL   *string        `gorm:"size:500" json:"avatar_url,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Messages  []Message  `gorm:"foreignKey:UserID" json:"-"`
	Rooms     []Room     `gorm:"many2many:room_members;" json:"-"`
	Files     []File     `gorm:"foreignKey:UserID" json:"-"`
	Reactions []Reaction `gorm:"foreignKey:UserID" json:"-"`
}

// TableName overrides the default table name
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook to set UUID if not set
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// CreateUserRequest for creating new users
type CreateUserRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=50"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=100"`
	AvatarURL   string `json:"avatar_url,omitempty"`
}
