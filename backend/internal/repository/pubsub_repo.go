package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/khonE3/chat-backend/internal/model"
	redisclient "github.com/khonE3/chat-backend/pkg/redis"
	"github.com/redis/go-redis/v9"
)

type PubSubRepository struct {
	redis *redisclient.Redis
}

func NewPubSubRepository(redis *redisclient.Redis) *PubSubRepository {
	return &PubSubRepository{redis: redis}
}

// Channel name formats
func roomChannel(roomID string) string {
	return fmt.Sprintf("chat:room:%s", roomID)
}

func typingChannel(roomID string) string {
	return fmt.Sprintf("chat:typing:%s", roomID)
}

func presenceChannel(roomID string) string {
	return fmt.Sprintf("chat:presence:%s", roomID)
}

// PublishMessage publishes a message to the room channel
func (r *PubSubRepository) PublishMessage(ctx context.Context, roomID string, msg *model.MessageWithUser) error {
	channel := roomChannel(roomID)

	wsMsg := model.WSMessage{
		Type:    model.WSTypeMessage,
		Payload: msg,
	}

	data, err := json.Marshal(wsMsg)
	if err != nil {
		return err
	}

	return r.redis.Client.Publish(ctx, channel, string(data)).Err()
}

// PublishTyping publishes typing status to the room
func (r *PubSubRepository) PublishTyping(ctx context.Context, roomID string, payload *model.TypingPayload) error {
	channel := typingChannel(roomID)

	msgType := model.WSTypeTyping
	if !payload.IsTyping {
		msgType = model.WSTypeStopTyping
	}

	wsMsg := model.WSMessage{
		Type:    msgType,
		Payload: payload,
	}

	data, err := json.Marshal(wsMsg)
	if err != nil {
		return err
	}

	return r.redis.Client.Publish(ctx, channel, string(data)).Err()
}

// PublishPresence publishes user presence (online/offline) to the room
func (r *PubSubRepository) PublishPresence(ctx context.Context, roomID string, payload *model.PresencePayload) error {
	channel := presenceChannel(roomID)

	wsMsg := model.WSMessage{
		Type:    model.WSTypePresence,
		Payload: payload,
	}

	data, err := json.Marshal(wsMsg)
	if err != nil {
		return err
	}

	return r.redis.Client.Publish(ctx, channel, string(data)).Err()
}

// Subscribe subscribes to all channels for a room
func (r *PubSubRepository) Subscribe(ctx context.Context, roomID string) *redis.PubSub {
	channels := []string{
		roomChannel(roomID),
		typingChannel(roomID),
		presenceChannel(roomID),
	}

	return r.redis.Client.Subscribe(ctx, channels...)
}

// Unsubscribe unsubscribes from all channels for a room
func (r *PubSubRepository) Unsubscribe(ctx context.Context, pubsub *redis.PubSub, roomID string) error {
	channels := []string{
		roomChannel(roomID),
		typingChannel(roomID),
		presenceChannel(roomID),
	}

	return pubsub.Unsubscribe(ctx, channels...)
}
