package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/model"
	"github.com/khonE3/chat-backend/pkg/database"
)

// GormReactionRepository implements reaction repository with GORM
type GormReactionRepository struct {
	db *database.GormDB
}

// NewGormReactionRepository creates a new GORM reaction repository
func NewGormReactionRepository(db *database.GormDB) *GormReactionRepository {
	return &GormReactionRepository{db: db}
}

// Add adds a reaction to a message
func (r *GormReactionRepository) Add(ctx context.Context, reaction *model.Reaction) error {
	// Use FirstOrCreate to handle duplicate reactions
	return r.db.DB.WithContext(ctx).
		Where("message_id = ? AND user_id = ? AND emoji = ?", reaction.MessageID, reaction.UserID, reaction.Emoji).
		FirstOrCreate(reaction).Error
}

// Remove removes a reaction from a message
func (r *GormReactionRepository) Remove(ctx context.Context, messageID, userID uuid.UUID, emoji string) error {
	return r.db.DB.WithContext(ctx).
		Where("message_id = ? AND user_id = ? AND emoji = ?", messageID, userID, emoji).
		Delete(&model.Reaction{}).Error
}

// GetByMessageID gets all reactions for a message
func (r *GormReactionRepository) GetByMessageID(ctx context.Context, messageID uuid.UUID) ([]model.Reaction, error) {
	var reactions []model.Reaction
	err := r.db.DB.WithContext(ctx).
		Where("message_id = ?", messageID).
		Find(&reactions).Error
	return reactions, err
}

// GetByMessageIDWithUsers gets reactions with user info
func (r *GormReactionRepository) GetByMessageIDWithUsers(ctx context.Context, messageID uuid.UUID) ([]model.ReactionWithUser, error) {
	var reactions []model.Reaction
	err := r.db.DB.WithContext(ctx).
		Where("message_id = ?", messageID).
		Find(&reactions).Error
	if err != nil {
		return nil, err
	}

	result := make([]model.ReactionWithUser, len(reactions))
	for i, reaction := range reactions {
		result[i] = model.ReactionWithUser{Reaction: reaction}

		var user model.User
		if err := r.db.DB.First(&user, "id = ?", reaction.UserID).Error; err == nil {
			result[i].Username = user.Username
			result[i].DisplayName = user.DisplayName
		}
	}

	return result, nil
}

// GetReactionCounts gets grouped reaction counts for a message
func (r *GormReactionRepository) GetReactionCounts(ctx context.Context, messageID uuid.UUID) ([]model.ReactionCount, error) {
	var results []struct {
		Emoji string
		Count int
	}

	err := r.db.DB.WithContext(ctx).
		Model(&model.Reaction{}).
		Select("emoji, COUNT(*) as count").
		Where("message_id = ?", messageID).
		Group("emoji").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make([]model.ReactionCount, len(results))
	for i, res := range results {
		counts[i] = model.ReactionCount{
			Emoji: res.Emoji,
			Count: res.Count,
		}

		// Get user IDs for this emoji
		var userIDs []string
		r.db.DB.Model(&model.Reaction{}).
			Select("user_id").
			Where("message_id = ? AND emoji = ?", messageID, res.Emoji).
			Pluck("user_id", &userIDs)
		counts[i].Users = userIDs
	}

	return counts, nil
}

// HasUserReacted checks if a user has already reacted with this emoji
func (r *GormReactionRepository) HasUserReacted(ctx context.Context, messageID, userID uuid.UUID, emoji string) (bool, error) {
	var count int64
	err := r.db.DB.WithContext(ctx).
		Model(&model.Reaction{}).
		Where("message_id = ? AND user_id = ? AND emoji = ?", messageID, userID, emoji).
		Count(&count).Error
	return count > 0, err
}

// Toggle toggles a reaction (add if not exists, remove if exists)
func (r *GormReactionRepository) Toggle(ctx context.Context, messageID, userID uuid.UUID, emoji string) (added bool, err error) {
	exists, err := r.HasUserReacted(ctx, messageID, userID, emoji)
	if err != nil {
		return false, err
	}

	if exists {
		err = r.Remove(ctx, messageID, userID, emoji)
		return false, err
	}

	reaction := &model.Reaction{
		MessageID: messageID,
		UserID:    userID,
		Emoji:     emoji,
	}
	err = r.Add(ctx, reaction)
	return true, err
}

// GetMessageReactionsSummary gets a summary of reactions for multiple messages
func (r *GormReactionRepository) GetMessageReactionsSummary(ctx context.Context, messageIDs []uuid.UUID) (map[uuid.UUID][]model.ReactionCount, error) {
	if len(messageIDs) == 0 {
		return map[uuid.UUID][]model.ReactionCount{}, nil
	}

	var reactions []model.Reaction
	err := r.db.DB.WithContext(ctx).
		Where("message_id IN ?", messageIDs).
		Find(&reactions).Error
	if err != nil {
		return nil, err
	}

	// Group by message and emoji
	summary := make(map[uuid.UUID]map[string][]string)
	for _, reaction := range reactions {
		if _, ok := summary[reaction.MessageID]; !ok {
			summary[reaction.MessageID] = make(map[string][]string)
		}
		summary[reaction.MessageID][reaction.Emoji] = append(
			summary[reaction.MessageID][reaction.Emoji],
			reaction.UserID.String(),
		)
	}

	// Convert to ReactionCount
	result := make(map[uuid.UUID][]model.ReactionCount)
	for msgID, emojiMap := range summary {
		for emoji, users := range emojiMap {
			result[msgID] = append(result[msgID], model.ReactionCount{
				Emoji: emoji,
				Count: len(users),
				Users: users,
			})
		}
	}

	return result, nil
}
