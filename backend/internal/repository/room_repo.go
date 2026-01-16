package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/model"
	"github.com/khonE3/chat-backend/pkg/database"
)

type RoomRepository struct {
	db *database.Postgres
}

func NewRoomRepository(db *database.Postgres) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(ctx context.Context, req *model.CreateRoomRequest, createdBy *uuid.UUID) (*model.Room, error) {
	room := &model.Room{
		ID:        uuid.New(),
		Name:      req.Name,
		IsPrivate: req.IsPrivate,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}

	if req.Description != "" {
		room.Description = &req.Description
	}

	query := `
		INSERT INTO rooms (id, name, description, is_private, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, description, is_private, created_by, created_at
	`

	err := r.db.Pool.QueryRow(ctx, query,
		room.ID, room.Name, room.Description, room.IsPrivate, room.CreatedBy, room.CreatedAt,
	).Scan(&room.ID, &room.Name, &room.Description, &room.IsPrivate, &room.CreatedBy, &room.CreatedAt)

	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *RoomRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Room, error) {
	room := &model.Room{}

	query := `
		SELECT id, name, description, is_private, created_by, created_at
		FROM rooms WHERE id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&room.ID, &room.Name, &room.Description, &room.IsPrivate, &room.CreatedBy, &room.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *RoomRepository) List(ctx context.Context, includePrivate bool) ([]model.RoomWithMembers, error) {
	query := `
		SELECT r.id, r.name, r.description, r.is_private, r.created_by, r.created_at,
			   COALESCE(COUNT(rm.user_id), 0) as member_count
		FROM rooms r
		LEFT JOIN room_members rm ON r.id = rm.room_id
	`

	if !includePrivate {
		query += ` WHERE r.is_private = false`
	}

	query += ` GROUP BY r.id ORDER BY r.created_at ASC`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []model.RoomWithMembers
	for rows.Next() {
		var room model.RoomWithMembers
		err := rows.Scan(
			&room.ID, &room.Name, &room.Description, &room.IsPrivate,
			&room.CreatedBy, &room.CreatedAt, &room.MemberCount,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *RoomRepository) AddMember(ctx context.Context, roomID, userID uuid.UUID) error {
	query := `
		INSERT INTO room_members (room_id, user_id, joined_at, last_read_at)
		VALUES ($1, $2, $3, $3)
		ON CONFLICT (room_id, user_id) DO NOTHING
	`

	_, err := r.db.Pool.Exec(ctx, query, roomID, userID, time.Now())
	return err
}

func (r *RoomRepository) RemoveMember(ctx context.Context, roomID, userID uuid.UUID) error {
	query := `DELETE FROM room_members WHERE room_id = $1 AND user_id = $2`
	_, err := r.db.Pool.Exec(ctx, query, roomID, userID)
	return err
}

func (r *RoomRepository) GetMembers(ctx context.Context, roomID uuid.UUID) ([]model.User, error) {
	query := `
		SELECT u.id, u.username, u.display_name, u.avatar_url, u.created_at, u.updated_at
		FROM users u
		INNER JOIN room_members rm ON u.id = rm.user_id
		WHERE rm.room_id = $1
		ORDER BY rm.joined_at ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Username, &user.DisplayName, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *RoomRepository) IsMember(ctx context.Context, roomID, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM room_members WHERE room_id = $1 AND user_id = $2)`

	var exists bool
	err := r.db.Pool.QueryRow(ctx, query, roomID, userID).Scan(&exists)
	return exists, err
}

func (r *RoomRepository) UpdateLastRead(ctx context.Context, roomID, userID uuid.UUID) error {
	query := `
		UPDATE room_members SET last_read_at = $3
		WHERE room_id = $1 AND user_id = $2
	`
	_, err := r.db.Pool.Exec(ctx, query, roomID, userID, time.Now())
	return err
}

func (r *RoomRepository) GetUnreadCount(ctx context.Context, roomID, userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM messages m
		INNER JOIN room_members rm ON m.room_id = rm.room_id AND rm.user_id = $2
		WHERE m.room_id = $1 AND m.created_at > rm.last_read_at
	`

	var count int
	err := r.db.Pool.QueryRow(ctx, query, roomID, userID).Scan(&count)
	return count, err
}

// ListWithUnread returns rooms with unread count for a specific user
func (r *RoomRepository) ListWithUnread(ctx context.Context, userID uuid.UUID, includePrivate bool) ([]model.RoomWithMembers, error) {
	query := `
		SELECT r.id, r.name, r.description, r.is_private, r.created_by, r.created_at,
		       COALESCE(COUNT(DISTINCT rm.user_id), 0) as member_count,
		       COALESCE((
		           SELECT COUNT(*)
		           FROM messages m
		           LEFT JOIN room_members urm ON m.room_id = urm.room_id AND urm.user_id = $1
		           WHERE m.room_id = r.id AND (urm.last_read_at IS NULL OR m.created_at > urm.last_read_at)
		       ), 0) as unread_count
		FROM rooms r
		LEFT JOIN room_members rm ON r.id = rm.room_id
	`

	if !includePrivate {
		query += ` WHERE r.is_private = false`
	}

	query += ` GROUP BY r.id ORDER BY r.created_at ASC`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []model.RoomWithMembers
	for rows.Next() {
		var room model.RoomWithMembers
		err := rows.Scan(
			&room.ID, &room.Name, &room.Description, &room.IsPrivate,
			&room.CreatedBy, &room.CreatedAt, &room.MemberCount, &room.UnreadCount,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}
