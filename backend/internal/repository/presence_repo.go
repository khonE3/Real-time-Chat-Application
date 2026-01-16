package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/model"
	redisclient "github.com/khonE3/chat-backend/pkg/redis"
	"github.com/redis/go-redis/v9"
)

type PresenceRepository struct {
	redis *redisclient.Redis
}

func NewPresenceRepository(redis *redisclient.Redis) *PresenceRepository {
	return &PresenceRepository{redis: redis}
}

// Key formats
func onlineUsersKey(roomID string) string {
	return fmt.Sprintf("chat:online:%s", roomID)
}

func userInfoKey(userID string) string {
	return fmt.Sprintf("chat:user:%s", userID)
}

// SetOnline adds a user to the online users sorted set for a room
func (r *PresenceRepository) SetOnline(ctx context.Context, roomID, userID string, user *model.User) error {
	key := onlineUsersKey(roomID)
	score := float64(time.Now().Unix())

	// Add to sorted set with current timestamp as score
	if err := r.redis.Client.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: userID,
	}).Err(); err != nil {
		return err
	}

	// Set expiry on the key (auto cleanup after 24 hours of inactivity)
	r.redis.Client.Expire(ctx, key, 24*time.Hour)

	// Store user info in a hash for quick lookup
	if user != nil {
		userData, _ := json.Marshal(model.OnlineUser{
			UserID:      userID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			AvatarURL:   user.AvatarURL,
			LastSeen:    time.Now(),
		})
		r.redis.Client.Set(ctx, userInfoKey(userID), userData, 24*time.Hour)
	}

	return nil
}

// SetOffline removes a user from the online users sorted set
func (r *PresenceRepository) SetOffline(ctx context.Context, roomID, userID string) error {
	key := onlineUsersKey(roomID)
	return r.redis.Client.ZRem(ctx, key, userID).Err()
}

// UpdateHeartbeat updates the user's last seen timestamp
func (r *PresenceRepository) UpdateHeartbeat(ctx context.Context, roomID, userID string) error {
	key := onlineUsersKey(roomID)
	score := float64(time.Now().Unix())

	return r.redis.Client.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: userID,
	}).Err()
}

// GetOnlineUsers returns all users who have been active in the last 5 minutes
func (r *PresenceRepository) GetOnlineUsers(ctx context.Context, roomID string) ([]model.OnlineUser, error) {
	key := onlineUsersKey(roomID)

	// Get users active in last 5 minutes
	cutoff := time.Now().Add(-5 * time.Minute).Unix()

	userIDs, err := r.redis.Client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: strconv.FormatInt(cutoff, 10),
		Max: "+inf",
	}).Result()

	if err != nil {
		return nil, err
	}

	var onlineUsers []model.OnlineUser
	for _, userID := range userIDs {
		// Get user info from cache
		data, err := r.redis.Client.Get(ctx, userInfoKey(userID)).Result()
		if err != nil {
			// If no cached info, create basic entry
			onlineUsers = append(onlineUsers, model.OnlineUser{
				UserID:   userID,
				LastSeen: time.Now(),
			})
			continue
		}

		var user model.OnlineUser
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			continue
		}
		onlineUsers = append(onlineUsers, user)
	}

	return onlineUsers, nil
}

// GetOnlineCount returns the number of online users in a room
func (r *PresenceRepository) GetOnlineCount(ctx context.Context, roomID string) (int64, error) {
	key := onlineUsersKey(roomID)
	cutoff := time.Now().Add(-5 * time.Minute).Unix()

	return r.redis.Client.ZCount(ctx, key, strconv.FormatInt(cutoff, 10), "+inf").Result()
}

// IsOnline checks if a user is currently online in a room
func (r *PresenceRepository) IsOnline(ctx context.Context, roomID, userID string) (bool, error) {
	key := onlineUsersKey(roomID)
	cutoff := time.Now().Add(-5 * time.Minute).Unix()

	score, err := r.redis.Client.ZScore(ctx, key, userID).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return int64(score) >= cutoff, nil
}

// CleanupStaleUsers removes users who haven't been active for more than 5 minutes
func (r *PresenceRepository) CleanupStaleUsers(ctx context.Context, roomID string) ([]string, error) {
	key := onlineUsersKey(roomID)
	cutoff := time.Now().Add(-5 * time.Minute).Unix()

	// Get stale users
	staleUsers, err := r.redis.Client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: "-inf",
		Max: strconv.FormatInt(cutoff, 10),
	}).Result()

	if err != nil {
		return nil, err
	}

	if len(staleUsers) > 0 {
		// Remove stale users
		members := make([]interface{}, len(staleUsers))
		for i, u := range staleUsers {
			members[i] = u
		}
		r.redis.Client.ZRem(ctx, key, members...)
	}

	return staleUsers, nil
}

// GetUserRooms returns all rooms where a user is currently online
func (r *PresenceRepository) GetUserRooms(ctx context.Context, userID string) ([]uuid.UUID, error) {
	// This would require scanning all room keys - for simplicity, we'll skip this
	// In production, you might want to maintain a reverse index
	return nil, nil
}
