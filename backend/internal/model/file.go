package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// File represents an uploaded file
type File struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	MessageID    *uuid.UUID     `gorm:"type:uuid;index" json:"message_id,omitempty"`
	Filename     string         `gorm:"not null;size:255" json:"filename"`
	OriginalName string         `gorm:"not null;size:255" json:"original_name"`
	MimeType     string         `gorm:"not null;size:100" json:"mime_type"`
	Size         int64          `gorm:"not null" json:"size"`
	URL          string         `gorm:"not null;size:500" json:"url"`
	ThumbnailURL *string        `gorm:"size:500" json:"thumbnail_url,omitempty"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User    *User    `gorm:"foreignKey:UserID" json:"-"`
	Message *Message `gorm:"foreignKey:MessageID" json:"-"`
}

// TableName overrides the default table name
func (File) TableName() string {
	return "files"
}

// BeforeCreate hook to set UUID if not set
func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// FileUploadResponse after successful upload
type FileUploadResponse struct {
	ID           string `json:"id"`
	Filename     string `json:"filename"`
	OriginalName string `json:"original_name"`
	MimeType     string `json:"mime_type"`
	Size         int64  `json:"size"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
}

// IsImage checks if the file is an image
func (f *File) IsImage() bool {
	switch f.MimeType {
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/svg+xml":
		return true
	}
	return false
}

// IsVideo checks if the file is a video
func (f *File) IsVideo() bool {
	switch f.MimeType {
	case "video/mp4", "video/webm", "video/ogg":
		return true
	}
	return false
}
