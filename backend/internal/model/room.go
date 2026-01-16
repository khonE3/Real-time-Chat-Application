package model

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	IsPrivate   bool       `json:"is_private"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type RoomMember struct {
	RoomID     uuid.UUID `json:"room_id"`
	UserID     uuid.UUID `json:"user_id"`
	JoinedAt   time.Time `json:"joined_at"`
	LastReadAt time.Time `json:"last_read_at"`
}

type RoomWithMembers struct {
	Room
	MemberCount int `json:"member_count"`
	OnlineCount int `json:"online_count"`
	UnreadCount int `json:"unread_count,omitempty"`
}

type CreateRoomRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description,omitempty"`
	IsPrivate   bool   `json:"is_private"`
}

type JoinRoomRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}
