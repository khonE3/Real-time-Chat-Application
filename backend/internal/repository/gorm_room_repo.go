package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/model"
	"github.com/khonE3/chat-backend/pkg/database"
	"gorm.io/gorm"
)

// GormRoomRepository implements room repository with GORM
type GormRoomRepository struct {
	db *database.GormDB
}

// NewGormRoomRepository creates a new GORM room repository
func NewGormRoomRepository(db *database.GormDB) *GormRoomRepository {
	return &GormRoomRepository{db: db}
}

// Create creates a new room
func (r *GormRoomRepository) Create(ctx context.Context, room *model.Room) error {
	return r.db.DB.WithContext(ctx).Create(room).Error
}

// GetByID retrieves a room by ID
func (r *GormRoomRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Room, error) {
	var room model.Room
	if err := r.db.DB.WithContext(ctx).First(&room, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

// List retrieves all public rooms
func (r *GormRoomRepository) List(ctx context.Context, includePrivate bool) ([]model.RoomWithMembers, error) {
	var rooms []model.Room
	query := r.db.DB.WithContext(ctx).Where("is_dm = ?", false)

	if !includePrivate {
		query = query.Where("is_private = ?", false)
	}

	if err := query.Order("created_at DESC").Find(&rooms).Error; err != nil {
		return nil, err
	}

	// Get member counts
	result := make([]model.RoomWithMembers, len(rooms))
	for i, room := range rooms {
		var memberCount int64
		r.db.DB.Model(&model.RoomMember{}).Where("room_id = ?", room.ID).Count(&memberCount)

		result[i] = model.RoomWithMembers{
			Room:        room,
			MemberCount: int(memberCount),
		}
	}

	return result, nil
}

// ListWithUnread retrieves rooms with unread counts for a user
func (r *GormRoomRepository) ListWithUnread(ctx context.Context, userID uuid.UUID, includePrivate bool) ([]model.RoomWithMembers, error) {
	rooms, err := r.List(ctx, includePrivate)
	if err != nil {
		return nil, err
	}

	for i := range rooms {
		// Get unread count
		var member model.RoomMember
		if err := r.db.DB.Where("room_id = ? AND user_id = ?", rooms[i].ID, userID).First(&member).Error; err == nil {
			var unreadCount int64
			r.db.DB.Model(&model.Message{}).
				Where("room_id = ? AND created_at > ?", rooms[i].ID, member.LastReadAt).
				Count(&unreadCount)
			rooms[i].UnreadCount = int(unreadCount)
		}
	}

	return rooms, nil
}

// AddMember adds a user to a room
func (r *GormRoomRepository) AddMember(ctx context.Context, roomID, userID uuid.UUID) error {
	member := &model.RoomMember{
		RoomID: roomID,
		UserID: userID,
	}
	return r.db.DB.WithContext(ctx).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		FirstOrCreate(member).Error
}

// RemoveMember removes a user from a room
func (r *GormRoomRepository) RemoveMember(ctx context.Context, roomID, userID uuid.UUID) error {
	return r.db.DB.WithContext(ctx).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Delete(&model.RoomMember{}).Error
}

// GetMembers retrieves all members of a room
func (r *GormRoomRepository) GetMembers(ctx context.Context, roomID uuid.UUID) ([]model.User, error) {
	var members []model.RoomMember
	if err := r.db.DB.WithContext(ctx).Where("room_id = ?", roomID).Find(&members).Error; err != nil {
		return nil, err
	}

	var userIDs []uuid.UUID
	for _, m := range members {
		userIDs = append(userIDs, m.UserID)
	}

	var users []model.User
	if len(userIDs) > 0 {
		if err := r.db.DB.WithContext(ctx).Where("id IN ?", userIDs).Find(&users).Error; err != nil {
			return nil, err
		}
	}

	return users, nil
}

// MarkAsRead marks messages as read for a user
func (r *GormRoomRepository) MarkAsRead(ctx context.Context, roomID, userID uuid.UUID) error {
	return r.db.DB.WithContext(ctx).
		Model(&model.RoomMember{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Update("last_read_at", gorm.Expr("NOW()")).Error
}

// CreateDM creates or gets a DM room between two users
func (r *GormRoomRepository) CreateDM(ctx context.Context, userID1, userID2 uuid.UUID) (*model.Room, error) {
	// Check if DM already exists
	var room model.Room
	err := r.db.DB.WithContext(ctx).
		Where("is_dm = ? AND ((dm_user_1 = ? AND dm_user_2 = ?) OR (dm_user_1 = ? AND dm_user_2 = ?))",
			true, userID1, userID2, userID2, userID1).
		First(&room).Error

	if err == nil {
		return &room, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Create new DM room
	room = model.Room{
		Name:      "DM",
		IsPrivate: true,
		IsDM:      true,
		DMUser1:   &userID1,
		DMUser2:   &userID2,
	}

	if err := r.db.DB.WithContext(ctx).Create(&room).Error; err != nil {
		return nil, err
	}

	// Add both users as members
	r.AddMember(ctx, room.ID, userID1)
	r.AddMember(ctx, room.ID, userID2)

	return &room, nil
}

// ListDMs retrieves all DM conversations for a user
func (r *GormRoomRepository) ListDMs(ctx context.Context, userID uuid.UUID) ([]model.Room, error) {
	var rooms []model.Room
	err := r.db.DB.WithContext(ctx).
		Where("is_dm = ? AND (dm_user_1 = ? OR dm_user_2 = ?)", true, userID, userID).
		Find(&rooms).Error
	return rooms, err
}
