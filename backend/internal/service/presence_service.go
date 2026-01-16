package service

import (
	"context"

	"github.com/khonE3/chat-backend/internal/model"
	"github.com/khonE3/chat-backend/internal/repository"
)

type PresenceService struct {
	presenceRepo *repository.PresenceRepository
}

func NewPresenceService(presenceRepo *repository.PresenceRepository) *PresenceService {
	return &PresenceService{presenceRepo: presenceRepo}
}

func (s *PresenceService) UserJoined(ctx context.Context, roomID, userID string, user *model.User) error {
	return s.presenceRepo.SetOnline(ctx, roomID, userID, user)
}

func (s *PresenceService) UserLeft(ctx context.Context, roomID, userID string) error {
	return s.presenceRepo.SetOffline(ctx, roomID, userID)
}

func (s *PresenceService) UpdateHeartbeat(ctx context.Context, roomID, userID string) error {
	return s.presenceRepo.UpdateHeartbeat(ctx, roomID, userID)
}

func (s *PresenceService) GetOnlineUsers(ctx context.Context, roomID string) ([]model.OnlineUser, error) {
	return s.presenceRepo.GetOnlineUsers(ctx, roomID)
}

func (s *PresenceService) GetOnlineCount(ctx context.Context, roomID string) (int64, error) {
	return s.presenceRepo.GetOnlineCount(ctx, roomID)
}

func (s *PresenceService) IsUserOnline(ctx context.Context, roomID, userID string) (bool, error) {
	return s.presenceRepo.IsOnline(ctx, roomID, userID)
}
