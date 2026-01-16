package model

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeSystem   MessageType = "system"
	MessageTypeTyping   MessageType = "typing"
	MessageTypePresence MessageType = "presence"
)

type Message struct {
	ID          uuid.UUID   `json:"id"`
	RoomID      uuid.UUID   `json:"room_id"`
	UserID      *uuid.UUID  `json:"user_id,omitempty"`
	Content     string      `json:"content"`
	MessageType MessageType `json:"message_type"`
	CreatedAt   time.Time   `json:"created_at"`
}

type MessageWithUser struct {
	Message
	Username    string  `json:"username,omitempty"`
	DisplayName string  `json:"display_name,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
}

type SendMessageRequest struct {
	Content string `json:"content" validate:"required,min=1,max=4000"`
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
)

type WSMessage struct {
	Type    WSMessageType `json:"type"`
	Payload interface{}   `json:"payload"`
}

type WSIncomingMessage struct {
	Type    WSMessageType `json:"type"`
	Content string        `json:"content,omitempty"`
	UserID  string        `json:"user_id,omitempty"`
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
