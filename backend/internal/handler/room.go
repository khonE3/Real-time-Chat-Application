package handler

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/model"
	"github.com/khonE3/chat-backend/internal/repository"
)

type RoomHandler struct {
	roomRepo *repository.RoomRepository
	userRepo *repository.UserRepository
}

func NewRoomHandler(roomRepo *repository.RoomRepository, userRepo *repository.UserRepository) *RoomHandler {
	return &RoomHandler{
		roomRepo: roomRepo,
		userRepo: userRepo,
	}
}

// List returns all public rooms
func (h *RoomHandler) List(c *fiber.Ctx) error {
	ctx := context.Background()

	// Check if userId is provided for unread counts
	userIDStr := c.Query("userId")
	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err == nil {
			rooms, err := h.roomRepo.ListWithUnread(ctx, userID, false)
			if err != nil {
				log.Printf("❌ Error fetching rooms with unread: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to fetch rooms",
				})
			}
			if rooms == nil {
				rooms = []model.RoomWithMembers{}
			}
			return c.JSON(rooms)
		}
	}

	rooms, err := h.roomRepo.List(ctx, false)
	if err != nil {
		log.Printf("❌ Error fetching rooms: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch rooms",
		})
	}

	if rooms == nil {
		rooms = []model.RoomWithMembers{}
	}

	return c.JSON(rooms)
}

// Create creates a new room
func (h *RoomHandler) Create(c *fiber.Ctx) error {
	var req model.CreateRoomRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Room name is required",
		})
	}

	// Get creator ID from query (simplified auth)
	creatorID := c.Query("userId")
	var createdBy *uuid.UUID
	if creatorID != "" {
		if id, err := uuid.Parse(creatorID); err == nil {
			createdBy = &id
		}
	}

	ctx := context.Background()
	room, err := h.roomRepo.Create(ctx, &req, createdBy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create room",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(room)
}

// GetByID gets a room by ID
func (h *RoomHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid room ID",
		})
	}

	ctx := context.Background()
	room, err := h.roomRepo.GetByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Room not found",
		})
	}

	return c.JSON(room)
}

// Join adds a user to a room
func (h *RoomHandler) Join(c *fiber.Ctx) error {
	roomIDStr := c.Params("id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid room ID",
		})
	}

	var req model.JoinRoomRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	ctx := context.Background()
	if err := h.roomRepo.AddMember(ctx, roomID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to join room",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully joined room",
	})
}

// GetMembers gets all members of a room
func (h *RoomHandler) GetMembers(c *fiber.Ctx) error {
	roomIDStr := c.Params("id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid room ID",
		})
	}

	ctx := context.Background()
	members, err := h.roomRepo.GetMembers(ctx, roomID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch members",
		})
	}

	if members == nil {
		members = []model.User{}
	}

	return c.JSON(members)
}

// MarkAsRead marks all messages in a room as read for a user
func (h *RoomHandler) MarkAsRead(c *fiber.Ctx) error {
	roomIDStr := c.Params("id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid room ID",
		})
	}

	userIDStr := c.Query("userId")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userId is required",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	ctx := context.Background()

	// Ensure user is a member first (add if not)
	if err := h.roomRepo.AddMember(ctx, roomID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add member",
		})
	}

	// Update last_read_at
	if err := h.roomRepo.UpdateLastRead(ctx, roomID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to mark as read",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Marked as read",
	})
}

// GetUnreadCount gets unread message count for a user in a room
func (h *RoomHandler) GetUnreadCount(c *fiber.Ctx) error {
	roomIDStr := c.Params("id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid room ID",
		})
	}

	userIDStr := c.Query("userId")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userId is required",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	ctx := context.Background()
	count, err := h.roomRepo.GetUnreadCount(ctx, roomID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get unread count",
		})
	}

	return c.JSON(fiber.Map{
		"unread_count": count,
	})
}
