package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	AvatarURL   *string   `json:"avatar_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=50"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=100"`
	AvatarURL   string `json:"avatar_url,omitempty"`
}
