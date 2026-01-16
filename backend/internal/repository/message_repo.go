package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/model"
	"github.com/khonE3/chat-backend/pkg/database"
	redisclient "github.com/khonE3/chat-backend/pkg/redis"
	"github.com/redis/go-redis/v9"
)

type MessageRepository struct {
	db    *database.Postgres
	redis *redisclient.Redis
}

func NewMessageRepository(db *database.Postgres, redis *redisclient.Redis) *MessageRepository {
	return &MessageRepository{db: db, redis: redis}
}

// PostgreSQL operations

func (r *MessageRepository) Create(ctx context.Context, roomID, userID uuid.UUID, content string, msgType model.MessageType) (*model.Message, error) {
	msg := &model.Message{
		ID:          uuid.New(),
		RoomID:      roomID,
		UserID:      &userID,
		Content:     content,
		MessageType: msgType,
		CreatedAt:   time.Now(),
	}

	query := `
		INSERT INTO messages (id, room_id, user_id, content, message_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, room_id, user_id, content, message_type, created_at
	`

	err := r.db.Pool.QueryRow(ctx, query,
		msg.ID, msg.RoomID, msg.UserID, msg.Content, msg.MessageType, msg.CreatedAt,
	).Scan(&msg.ID, &msg.RoomID, &msg.UserID, &msg.Content, &msg.MessageType, &msg.CreatedAt)

	if err != nil {
		return nil, err
	}

	// Also add to Redis Stream for recent messages cache
	go r.addToStream(context.Background(), msg)

	return msg, nil
}

func (r *MessageRepository) GetByRoom(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]model.MessageWithUser, error) {
	query := `
		SELECT m.id, m.room_id, m.user_id, m.content, m.message_type, m.created_at,
			   COALESCE(u.username, 'deleted') as username,
			   COALESCE(u.display_name, 'Deleted User') as display_name,
			   u.avatar_url
		FROM messages m
		LEFT JOIN users u ON m.user_id = u.id
		WHERE m.room_id = $1
		ORDER BY m.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, roomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.MessageWithUser
	for rows.Next() {
		var msg model.MessageWithUser
		err := rows.Scan(
			&msg.ID, &msg.RoomID, &msg.UserID, &msg.Content, &msg.MessageType, &msg.CreatedAt,
			&msg.Username, &msg.DisplayName, &msg.AvatarURL,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *MessageRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.MessageWithUser, error) {
	msg := &model.MessageWithUser{}

	query := `
		SELECT m.id, m.room_id, m.user_id, m.content, m.message_type, m.created_at,
			   COALESCE(u.username, 'deleted') as username,
			   COALESCE(u.display_name, 'Deleted User') as display_name,
			   u.avatar_url
		FROM messages m
		LEFT JOIN users u ON m.user_id = u.id
		WHERE m.id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&msg.ID, &msg.RoomID, &msg.UserID, &msg.Content, &msg.MessageType, &msg.CreatedAt,
		&msg.Username, &msg.DisplayName, &msg.AvatarURL,
	)

	if err != nil {
		return nil, err
	}

	return msg, nil
}

// Redis Stream operations

func (r *MessageRepository) addToStream(ctx context.Context, msg *model.Message) error {
	streamKey := fmt.Sprintf("chat:stream:%s", msg.RoomID.String())

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = r.redis.Client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		MaxLen: 100, // Keep only last 100 messages in stream
		Approx: true,
		Values: map[string]interface{}{
			"data": string(data),
		},
	}).Result()

	return err
}

func (r *MessageRepository) GetRecentFromStream(ctx context.Context, roomID uuid.UUID, count int64) ([]model.Message, error) {
	streamKey := fmt.Sprintf("chat:stream:%s", roomID.String())

	result, err := r.redis.Client.XRevRange(ctx, streamKey, "+", "-").Result()
	if err != nil {
		return nil, err
	}

	var messages []model.Message
	for i, stream := range result {
		if int64(i) >= count {
			break
		}

		data, ok := stream.Values["data"].(string)
		if !ok {
			continue
		}

		var msg model.Message
		if err := json.Unmarshal([]byte(data), &msg); err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
