package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MessageType defines types of messages
type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeSystem   MessageType = "system"
	MessageTypeTyping   MessageType = "typing"
	MessageTypePresence MessageType = "presence"
	MessageTypeFile     MessageType = "file"
)

// Message represents a chat message
type Message struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	RoomID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"room_id"`
	UserID      *uuid.UUID     `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Content     string         `gorm:"type:text" json:"content"`
	MessageType MessageType    `gorm:"type:varchar(20);default:'text'" json:"message_type"`
	CreatedAt   time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Room      *Room      `gorm:"foreignKey:RoomID" json:"-"`
	User      *User      `gorm:"foreignKey:UserID" json:"-"`
	Files     []File     `gorm:"foreignKey:MessageID" json:"files,omitempty"`
	Reactions []Reaction `gorm:"foreignKey:MessageID" json:"reactions,omitempty"`
}

// TableName overrides the default table name
func (Message) TableName() string {
	return "messages"
}

// BeforeCreate hook to set UUID if not set
func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// MessageWithUser includes user information
type MessageWithUser struct {
	Message
	Username    string  `json:"username,omitempty"`
	DisplayName string  `json:"display_name,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
}

// SendMessageRequest for sending messages
type SendMessageRequest struct {
	Content string   `json:"content" validate:"required,min=1,max=4000"`
	FileIDs []string `json:"file_ids,omitempty"`
}

// WebSocket message types
type WSMessageType string

const (
	WSTypeMessage     WSMessageType = "message"
	WSTypeTyping      WSMessageType = "typing"
	WSTypeStopTyping  WSMessageType = "stop_typing"
	WSTypePresence    WSMessageType = "presence"
	WSTypeHistory     WSMessageType = "history"
	WSTypeOnlineUsers WSMessageType = "online_users"
	WSTypeError       WSMessageType = "error"
	WSTypeJoin        WSMessageType = "join"
	WSTypeLeave       WSMessageType = "leave"
	WSTypeReaction    WSMessageType = "reaction"
	WSTypeReactionAdd WSMessageType = "reaction_add"
	WSTypeReactionDel WSMessageType = "reaction_remove"
	WSTypeFileUpload  WSMessageType = "file_upload"
)

type WSMessage struct {
	Type    WSMessageType `json:"type"`
	Payload interface{}   `json:"payload"`
}

type WSIncomingMessage struct {
	Type      WSMessageType `json:"type"`
	Content   string        `json:"content,omitempty"`
	UserID    string        `json:"user_id,omitempty"`
	MessageID string        `json:"message_id,omitempty"`
	Emoji     string        `json:"emoji,omitempty"`
}

type TypingPayload struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	IsTyping    bool   `json:"is_typing"`
}

type PresencePayload struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	IsOnline    bool   `json:"is_online"`
}

type OnlineUser struct {
	UserID      string    `json:"user_id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	AvatarURL   *string   `json:"avatar_url,omitempty"`
	LastSeen    time.Time `json:"last_seen"`
}
