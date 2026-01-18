package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Room represents a chat room
type Room struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"not null;size:100" json:"name"`
	Description *string        `gorm:"size:500" json:"description,omitempty"`
	IsPrivate   bool           `gorm:"default:false" json:"is_private"`
	IsDM        bool           `gorm:"default:false" json:"is_dm"`
	DMUser1     *uuid.UUID     `gorm:"type:uuid;index" json:"dm_user_1,omitempty"`
	DMUser2     *uuid.UUID     `gorm:"type:uuid;index" json:"dm_user_2,omitempty"`
	CreatedBy   *uuid.UUID     `gorm:"type:uuid" json:"created_by,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Creator  *User     `gorm:"foreignKey:CreatedBy" json:"-"`
	DMUserA  *User     `gorm:"foreignKey:DMUser1" json:"-"`
	DMUserB  *User     `gorm:"foreignKey:DMUser2" json:"-"`
	Members  []User    `gorm:"many2many:room_members;" json:"-"`
	Messages []Message `gorm:"foreignKey:RoomID" json:"-"`
}

// TableName overrides the default table name
func (Room) TableName() string {
	return "rooms"
}

// BeforeCreate hook to set UUID if not set
func (r *Room) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// RoomMember represents room membership
type RoomMember struct {
	RoomID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"room_id"`
	UserID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	JoinedAt   time.Time `gorm:"autoCreateTime" json:"joined_at"`
	LastReadAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"last_read_at"`

	// Relationships
	Room *Room `gorm:"foreignKey:RoomID" json:"-"`
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName overrides the default table name
func (RoomMember) TableName() string {
	return "room_members"
}

// RoomWithMembers includes member counts
type RoomWithMembers struct {
	Room
	MemberCount int `gorm:"-" json:"member_count"`
	OnlineCount int `gorm:"-" json:"online_count"`
	UnreadCount int `gorm:"-" json:"unread_count,omitempty"`
}

// CreateRoomRequest for creating new rooms
type CreateRoomRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description,omitempty"`
	IsPrivate   bool   `json:"is_private"`
}

// JoinRoomRequest for joining rooms
type JoinRoomRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}

// StartDMRequest for starting direct messages
type StartDMRequest struct {
	TargetUserID string `json:"target_user_id" validate:"required,uuid"`
}
