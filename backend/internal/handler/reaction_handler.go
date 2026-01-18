package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/model"
	"github.com/khonE3/chat-backend/internal/repository"
)

// ReactionHandler handles reaction operations
type ReactionHandler struct {
	reactionRepo *repository.GormReactionRepository
	messageRepo  *repository.GormMessageRepository
}

// NewReactionHandler creates a new reaction handler
func NewReactionHandler(reactionRepo *repository.GormReactionRepository, messageRepo *repository.GormMessageRepository) *ReactionHandler {
	return &ReactionHandler{
		reactionRepo: reactionRepo,
		messageRepo:  messageRepo,
	}
}

// AddReaction adds a reaction to a message
func (h *ReactionHandler) AddReaction(c *fiber.Ctx) error {
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userId is required",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid userId",
		})
	}

	var req model.AddReactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	messageID, err := uuid.Parse(req.MessageID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid message ID",
		})
	}

	// Check if message exists
	_, err = h.messageRepo.GetByID(c.Context(), messageID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Message not found",
		})
	}

	// Validate emoji
	if req.Emoji == "" || len(req.Emoji) > 20 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid emoji",
		})
	}

	reaction := &model.Reaction{
		MessageID: messageID,
		UserID:    userID,
		Emoji:     req.Emoji,
	}

	if err := h.reactionRepo.Add(c.Context(), reaction); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add reaction",
		})
	}

	// Get updated reaction counts
	counts, _ := h.reactionRepo.GetReactionCounts(c.Context(), messageID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"reaction": reaction,
		"counts":   counts,
	})
}

// RemoveReaction removes a reaction from a message
func (h *ReactionHandler) RemoveReaction(c *fiber.Ctx) error {
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userId is required",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid userId",
		})
	}

	var req model.RemoveReactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	messageID, err := uuid.Parse(req.MessageID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid message ID",
		})
	}

	if err := h.reactionRepo.Remove(c.Context(), messageID, userID, req.Emoji); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove reaction",
		})
	}

	// Get updated reaction counts
	counts, _ := h.reactionRepo.GetReactionCounts(c.Context(), messageID)

	return c.JSON(fiber.Map{
		"message": "Reaction removed",
		"counts":  counts,
	})
}

// ToggleReaction toggles a reaction (add if not exists, remove if exists)
func (h *ReactionHandler) ToggleReaction(c *fiber.Ctx) error {
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userId is required",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid userId",
		})
	}

	messageIDStr := c.Params("messageId")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid message ID",
		})
	}

	emoji := c.Params("emoji")
	if emoji == "" || len(emoji) > 20 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid emoji",
		})
	}

	added, err := h.reactionRepo.Toggle(c.Context(), messageID, userID, emoji)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to toggle reaction",
		})
	}

	// Get updated reaction counts
	counts, _ := h.reactionRepo.GetReactionCounts(c.Context(), messageID)

	action := "removed"
	if added {
		action = "added"
	}

	return c.JSON(fiber.Map{
		"action": action,
		"counts": counts,
	})
}

// GetReactions gets all reactions for a message
func (h *ReactionHandler) GetReactions(c *fiber.Ctx) error {
	messageIDStr := c.Params("messageId")
	messageID, err := uuid.Parse(messageIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid message ID",
		})
	}

	counts, err := h.reactionRepo.GetReactionCounts(c.Context(), messageID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get reactions",
		})
	}

	return c.JSON(counts)
}

// GetCommonEmojis returns the common reaction emojis
func (h *ReactionHandler) GetCommonEmojis(c *fiber.Ctx) error {
	return c.JSON(model.CommonReactions)
}
