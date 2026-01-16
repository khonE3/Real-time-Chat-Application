package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/repository"
)

type MessageHandler struct {
	messageRepo *repository.MessageRepository
}

func NewMessageHandler(messageRepo *repository.MessageRepository) *MessageHandler {
	return &MessageHandler{messageRepo: messageRepo}
}

// GetByRoom gets messages for a room with pagination
func (h *MessageHandler) GetByRoom(c *fiber.Ctx) error {
	roomIDStr := c.Params("id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid room ID",
		})
	}

	// Parse pagination params
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 50
	}

	ctx := context.Background()
	messages, err := h.messageRepo.GetByRoom(ctx, roomID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch messages",
		})
	}

	return c.JSON(fiber.Map{
		"messages": messages,
		"limit":    limit,
		"offset":   offset,
	})
}
