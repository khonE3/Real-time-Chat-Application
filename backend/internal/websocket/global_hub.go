package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/khonE3/chat-backend/internal/model"
	"github.com/khonE3/chat-backend/internal/repository"
)

// GlobalClient represents a client subscribed to global updates
type GlobalClient struct {
	ID     string
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *GlobalHub
	mu     sync.Mutex
}

// GlobalHub maintains global subscribers for homepage updates
type GlobalHub struct {
	// All connected global clients
	clients map[*GlobalClient]bool

	// Register requests
	register chan *GlobalClient

	// Unregister requests
	unregister chan *GlobalClient

	// Broadcast to all clients
	broadcast chan []byte

	// Room repository for fetching room stats
	roomRepo *repository.RoomRepository

	mu sync.RWMutex
}

// GlobalMessage types
type GlobalMessageType string

const (
	GlobalTypeRoomStats   GlobalMessageType = "room_stats"
	GlobalTypeRoomCreated GlobalMessageType = "room_created"
	GlobalTypeRoomDeleted GlobalMessageType = "room_deleted"
	GlobalTypePresence    GlobalMessageType = "global_presence"
)

type GlobalMessage struct {
	Type    GlobalMessageType `json:"type"`
	Payload interface{}       `json:"payload"`
}

type RoomStatsPayload struct {
	RoomID      string `json:"room_id"`
	OnlineCount int    `json:"online_count"`
	UnreadCount int    `json:"unread_count,omitempty"`
}

type GlobalPresencePayload struct {
	TotalOnline int `json:"total_online"`
}

func NewGlobalHub(roomRepo *repository.RoomRepository) *GlobalHub {
	return &GlobalHub{
		clients:    make(map[*GlobalClient]bool),
		register:   make(chan *GlobalClient),
		unregister: make(chan *GlobalClient),
		broadcast:  make(chan []byte, 256),
		roomRepo:   roomRepo,
	}
}

func (h *GlobalHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("🌐 Global client connected: %s (total: %d)", client.UserID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("🌐 Global client disconnected: %s (total: %d)", client.UserID, len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					// Client buffer full, remove it
					go func(c *GlobalClient) {
						h.unregister <- c
					}(client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastRoomStats sends room statistics update to all global clients
func (h *GlobalHub) BroadcastRoomStats(roomID string, onlineCount int) {
	msg := GlobalMessage{
		Type: GlobalTypeRoomStats,
		Payload: RoomStatsPayload{
			RoomID:      roomID,
			OnlineCount: onlineCount,
		},
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal room stats: %v", err)
		return
	}
	h.broadcast <- data
}

// BroadcastNewMessage notifies about new message for unread count
func (h *GlobalHub) BroadcastNewMessage(roomID string, senderUserID string) {
	// Broadcast that room has new message
	msg := GlobalMessage{
		Type: GlobalTypeRoomStats,
		Payload: map[string]interface{}{
			"room_id":     roomID,
			"has_new_msg": true,
			"sender_id":   senderUserID,
		},
	}
	data, _ := json.Marshal(msg)
	h.broadcast <- data
}

// BroadcastRoomCreated notifies about new room
func (h *GlobalHub) BroadcastRoomCreated(room *model.Room) {
	msg := GlobalMessage{
		Type:    GlobalTypeRoomCreated,
		Payload: room,
	}
	data, _ := json.Marshal(msg)
	h.broadcast <- data
}

// BroadcastTotalOnline sends total online count
func (h *GlobalHub) BroadcastTotalOnline(total int) {
	msg := GlobalMessage{
		Type: GlobalTypePresence,
		Payload: GlobalPresencePayload{
			TotalOnline: total,
		},
	}
	data, _ := json.Marshal(msg)
	h.broadcast <- data
}

// HandleGlobalWebSocket handles WebSocket connection for global updates
func (h *GlobalHub) HandleGlobalWebSocket(c *websocket.Conn) {
	userID := c.Query("userId", "anonymous")

	client := &GlobalClient{
		ID:     userID,
		UserID: userID,
		Conn:   c,
		Send:   make(chan []byte, 256),
		Hub:    h,
	}

	h.register <- client

	// Send initial room stats
	go func() {
		ctx := context.Background()
		rooms, err := h.roomRepo.List(ctx, false)
		if err == nil {
			msg := GlobalMessage{
				Type:    "rooms_init",
				Payload: rooms,
			}
			data, _ := json.Marshal(msg)
			client.Send <- data
		}
	}()

	// Start goroutines
	go client.writePump()
	client.readPump()
}

func (c *GlobalClient) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		// Global clients don't send messages, just receive
	}
}

func (c *GlobalClient) writePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		c.mu.Lock()
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		c.mu.Unlock()

		if err != nil {
			return
		}
	}
}

// GetOnlineCount returns number of global clients
func (h *GlobalHub) GetOnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
