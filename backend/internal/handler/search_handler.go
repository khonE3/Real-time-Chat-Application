package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/repository"
)

// SearchHandler handles search operations
type SearchHandler struct {
	messageRepo *repository.GormMessageRepository
	userRepo    *repository.GormUserRepository
	roomRepo    *repository.GormRoomRepository
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(
	messageRepo *repository.GormMessageRepository,
	userRepo *repository.GormUserRepository,
	roomRepo *repository.GormRoomRepository,
) *SearchHandler {
	return &SearchHandler{
		messageRepo: messageRepo,
		userRepo:    userRepo,
		roomRepo:    roomRepo,
	}
}

// SearchMessages searches messages in a room
func (h *SearchHandler) SearchMessages(c *fiber.Ctx) error {
	roomIDStr := c.Params("roomId")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid room ID",
		})
	}

	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	messages, err := h.messageRepo.Search(c.Context(), roomID, query, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search messages",
		})
	}

	return c.JSON(fiber.Map{
		"query":    query,
		"count":    len(messages),
		"messages": messages,
	})
}

// GetMessagesBefore gets messages before a specific message ID (for infinite scroll)
func (h *SearchHandler) GetMessagesBefore(c *fiber.Ctx) error {
	roomIDStr := c.Params("roomId")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid room ID",
		})
	}

	beforeIDStr := c.Query("before")
	if beforeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "before parameter is required",
		})
	}

	beforeID, err := uuid.Parse(beforeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid before ID",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	messages, err := h.messageRepo.GetByRoomBefore(c.Context(), roomID, beforeID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get messages",
		})
	}

	return c.JSON(fiber.Map{
		"count":    len(messages),
		"messages": messages,
		"has_more": len(messages) == limit,
	})
}

// SearchUsers searches users by username or display name
func (h *SearchHandler) SearchUsers(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.JSON([]interface{}{})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	users, err := h.userRepo.Search(c.Context(), query, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search users",
		})
	}

	return c.JSON(users)
}

// GlobalSearch searches across messages and users
func (h *SearchHandler) GlobalSearch(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	// Search users
	users, _ := h.userRepo.Search(c.Context(), query, 5)

	// Search in specific room if provided
	var messages interface{}
	roomIDStr := c.Query("room")
	if roomIDStr != "" {
		roomID, err := uuid.Parse(roomIDStr)
		if err == nil {
			msgs, _ := h.messageRepo.Search(c.Context(), roomID, query, 10)
			messages = msgs
		}
	}

	return c.JSON(fiber.Map{
		"query":    query,
		"users":    users,
		"messages": messages,
	})
}
