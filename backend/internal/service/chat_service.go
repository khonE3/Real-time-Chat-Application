package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/model"
	"github.com/khonE3/chat-backend/internal/repository"
)

type ChatService struct {
	messageRepo  *repository.MessageRepository
	pubsubRepo   *repository.PubSubRepository
	presenceRepo *repository.PresenceRepository
}

func NewChatService(
	messageRepo *repository.MessageRepository,
	pubsubRepo *repository.PubSubRepository,
	presenceRepo *repository.PresenceRepository,
) *ChatService {
	return &ChatService{
		messageRepo:  messageRepo,
		pubsubRepo:   pubsubRepo,
		presenceRepo: presenceRepo,
	}
}

func (s *ChatService) SendMessage(ctx context.Context, roomID string, userID uuid.UUID, content string) (*model.MessageWithUser, error) {
	roomUUID, err := uuid.Parse(roomID)
	if err != nil {
		return nil, err
	}

	// Save message to database
	msg, err := s.messageRepo.Create(ctx, roomUUID, userID, content, model.MessageTypeText)
	if err != nil {
		return nil, err
	}

	// Get message with user info
	return s.messageRepo.GetByID(ctx, msg.ID)
}

func (s *ChatService) GetRecentMessages(ctx context.Context, roomID string, limit int) ([]model.MessageWithUser, error) {
	roomUUID, err := uuid.Parse(roomID)
	if err != nil {
		return nil, err
	}

	return s.messageRepo.GetByRoom(ctx, roomUUID, limit, 0)
}

func (s *ChatService) GetMessageHistory(ctx context.Context, roomID string, limit, offset int) ([]model.MessageWithUser, error) {
	roomUUID, err := uuid.Parse(roomID)
	if err != nil {
		return nil, err
	}

	return s.messageRepo.GetByRoom(ctx, roomUUID, limit, offset)
}
