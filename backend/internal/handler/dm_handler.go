package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/model"
	"github.com/khonE3/chat-backend/internal/repository"
)

// DMHandler handles direct message operations
type DMHandler struct {
	roomRepo *repository.GormRoomRepository
	userRepo *repository.GormUserRepository
}

// NewDMHandler creates a new DM handler
func NewDMHandler(roomRepo *repository.GormRoomRepository, userRepo *repository.GormUserRepository) *DMHandler {
	return &DMHandler{
		roomRepo: roomRepo,
		userRepo: userRepo,
	}
}

// StartDM creates or gets a DM room between two users
func (h *DMHandler) StartDM(c *fiber.Ctx) error {
	// Get current user ID from query
	currentUserIDStr := c.Query("userId")
	if currentUserIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userId is required",
		})
	}

	currentUserID, err := uuid.Parse(currentUserIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid userId",
		})
	}

	// Parse request body
	var req model.StartDMRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	targetUserID, err := uuid.Parse(req.TargetUserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid target user ID",
		})
	}

	// Cannot DM yourself
	if currentUserID == targetUserID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot start DM with yourself",
		})
	}

	// Check if target user exists
	targetUser, err := h.userRepo.GetByID(c.Context(), targetUserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Target user not found",
		})
	}

	// Create or get DM room
	room, err := h.roomRepo.CreateDM(c.Context(), currentUserID, targetUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create DM room",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"room":        room,
		"target_user": targetUser,
	})
}

// ListDMs gets all DM conversations for current user
func (h *DMHandler) ListDMs(c *fiber.Ctx) error {
	// Get current user ID from query
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

	rooms, err := h.roomRepo.ListDMs(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list DMs",
		})
	}

	// Enrich with user information
	type DMResponse struct {
		model.Room
		OtherUser *model.User `json:"other_user"`
	}

	var response []DMResponse
	for _, room := range rooms {
		dm := DMResponse{Room: room}

		// Get the other user in the DM
		var otherUserID uuid.UUID
		if room.DMUser1 != nil && *room.DMUser1 != userID {
			otherUserID = *room.DMUser1
		} else if room.DMUser2 != nil && *room.DMUser2 != userID {
			otherUserID = *room.DMUser2
		}

		if otherUserID != uuid.Nil {
			otherUser, err := h.userRepo.GetByID(c.Context(), otherUserID)
			if err == nil {
				dm.OtherUser = otherUser
			}
		}

		response = append(response, dm)
	}

	return c.JSON(response)
}

// GetDM gets a specific DM room with a user
func (h *DMHandler) GetDM(c *fiber.Ctx) error {
	// Get current user ID from query
	currentUserIDStr := c.Query("userId")
	if currentUserIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "userId is required",
		})
	}

	currentUserID, err := uuid.Parse(currentUserIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid userId",
		})
	}

	// Get target user ID from params
	targetUserIDStr := c.Params("targetUserId")
	targetUserID, err := uuid.Parse(targetUserIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid target user ID",
		})
	}

	// Get or create DM room
	room, err := h.roomRepo.CreateDM(c.Context(), currentUserID, targetUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get DM room",
		})
	}

	// Get target user info
	targetUser, _ := h.userRepo.GetByID(c.Context(), targetUserID)

	return c.JSON(fiber.Map{
		"room":        room,
		"target_user": targetUser,
	})
}
