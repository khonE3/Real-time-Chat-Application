package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/khonE3/chat-backend/internal/model"
	"github.com/khonE3/chat-backend/internal/repository"
	"github.com/khonE3/chat-backend/internal/service"
)

// Client represents a WebSocket client connection
type Client struct {
	ID          string
	UserID      uuid.UUID
	Username    string
	DisplayName string
	RoomID      string
	Conn        *websocket.Conn
	Hub         *Hub
	Send        chan []byte
	mu          sync.Mutex
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients per room
	rooms map[string]map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Broadcast to room
	broadcast chan *RoomMessage

	// Services
	chatService     *service.ChatService
	presenceService *service.PresenceService
	pubsubRepo      *repository.PubSubRepository

	mu sync.RWMutex
}

type RoomMessage struct {
	RoomID  string
	Message []byte
}

func NewHub(chatService *service.ChatService, presenceService *service.PresenceService, pubsubRepo *repository.PubSubRepository) *Hub {
	return &Hub{
		rooms:           make(map[string]map[*Client]bool),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		broadcast:       make(chan *RoomMessage),
		chatService:     chatService,
		presenceService: presenceService,
		pubsubRepo:      pubsubRepo,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case roomMsg := <-h.broadcast:
			h.broadcastToRoom(roomMsg)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[client.RoomID]; !ok {
		h.rooms[client.RoomID] = make(map[*Client]bool)
	}
	h.rooms[client.RoomID][client] = true

	log.Printf("👤 Client %s joined room %s", client.Username, client.RoomID)

	// Update presence
	go func() {
		ctx := context.Background()
		h.presenceService.UserJoined(ctx, client.RoomID, client.UserID.String(), &model.User{
			ID:          client.UserID,
			Username:    client.Username,
			DisplayName: client.DisplayName,
		})

		// Publish presence update
		h.pubsubRepo.PublishPresence(ctx, client.RoomID, &model.PresencePayload{
			UserID:      client.UserID.String(),
			Username:    client.Username,
			DisplayName: client.DisplayName,
			IsOnline:    true,
		})

		// Send online users list to the new client
		onlineUsers, err := h.presenceService.GetOnlineUsers(ctx, client.RoomID)
		if err != nil {
			log.Printf("Failed to get online users for room %s: %v", client.RoomID, err)
		}
		h.sendToClient(client, model.WSMessage{
			Type:    model.WSTypeOnlineUsers,
			Payload: onlineUsers,
		})

		// Send recent message history
		messages, err := h.chatService.GetRecentMessages(ctx, client.RoomID, 50)
		if err != nil {
			log.Printf("Failed to get recent messages for room %s: %v", client.RoomID, err)
		}
		h.sendToClient(client, model.WSMessage{
			Type:    model.WSTypeHistory,
			Payload: messages,
		})
	}()
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	if clients, ok := h.rooms[client.RoomID]; ok {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			close(client.Send)

			if len(clients) == 0 {
				delete(h.rooms, client.RoomID)
			}
		}
	}
	h.mu.Unlock()

	log.Printf("👋 Client %s left room %s", client.Username, client.RoomID)

	// Update presence
	go func() {
		ctx := context.Background()
		h.presenceService.UserLeft(ctx, client.RoomID, client.UserID.String())

		// Publish presence update
		h.pubsubRepo.PublishPresence(ctx, client.RoomID, &model.PresencePayload{
			UserID:      client.UserID.String(),
			Username:    client.Username,
			DisplayName: client.DisplayName,
			IsOnline:    false,
		})
	}()
}

func (h *Hub) broadcastToRoom(roomMsg *RoomMessage) {
	h.mu.RLock()
	clients, ok := h.rooms[roomMsg.RoomID]
	h.mu.RUnlock()

	if !ok {
		return
	}

	for client := range clients {
		select {
		case client.Send <- roomMsg.Message:
		default:
			h.unregister <- client
		}
	}
}

func (h *Hub) sendToClient(client *Client, msg model.WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	select {
	case client.Send <- data:
	default:
		// Client buffer full
	}
}

func (h *Hub) HandleWebSocket(c *websocket.Conn) {
	roomID := c.Params("roomId")
	userID := c.Query("userId")
	username := c.Query("username", "anonymous")
	displayName := c.Query("displayName", username)

	if roomID == "" || userID == "" {
		c.WriteJSON(model.WSMessage{
			Type:    model.WSTypeError,
			Payload: "Missing roomId or userId",
		})
		c.Close()
		return
	}

	if _, err := uuid.Parse(roomID); err != nil {
		c.WriteJSON(model.WSMessage{
			Type:    model.WSTypeError,
			Payload: "Invalid roomId format",
		})
		c.Close()
		return
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		c.WriteJSON(model.WSMessage{
			Type:    model.WSTypeError,
			Payload: "Invalid userId format",
		})
		c.Close()
		return
	}

	client := &Client{
		ID:          uuid.New().String(),
		UserID:      parsedUserID,
		Username:    username,
		DisplayName: displayName,
		RoomID:      roomID,
		Conn:        c,
		Hub:         h,
		Send:        make(chan []byte, 256),
	}

	h.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var incoming model.WSIncomingMessage
		if err := json.Unmarshal(message, &incoming); err != nil {
			continue
		}

		c.handleMessage(&incoming)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.mu.Lock()
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			c.mu.Unlock()

			if err != nil {
				return
			}

		case <-ticker.C:
			c.mu.Lock()
			err := c.Conn.WriteMessage(websocket.PingMessage, nil)
			c.mu.Unlock()

			if err != nil {
				return
			}

			// Update heartbeat
			go c.Hub.presenceService.UpdateHeartbeat(context.Background(), c.RoomID, c.UserID.String())
		}
	}
}

func (c *Client) handleMessage(msg *model.WSIncomingMessage) {
	ctx := context.Background()

	switch msg.Type {
	case model.WSTypeMessage:
		// Save message and broadcast
		savedMsg, err := c.Hub.chatService.SendMessage(ctx, c.RoomID, c.UserID, msg.Content)
		if err != nil {
			log.Printf("Error saving message: %v", err)
			return
		}

		// Broadcast to room via hub
		wsMsg := model.WSMessage{
			Type:    model.WSTypeMessage,
			Payload: savedMsg,
		}
		data, _ := json.Marshal(wsMsg)
		c.Hub.broadcast <- &RoomMessage{
			RoomID:  c.RoomID,
			Message: data,
		}

	case model.WSTypeTyping:
		// Broadcast typing indicator
		c.Hub.pubsubRepo.PublishTyping(ctx, c.RoomID, &model.TypingPayload{
			UserID:      c.UserID.String(),
			Username:    c.Username,
			DisplayName: c.DisplayName,
			IsTyping:    true,
		})

		// Broadcast to local clients
		wsMsg := model.WSMessage{
			Type: model.WSTypeTyping,
			Payload: model.TypingPayload{
				UserID:      c.UserID.String(),
				Username:    c.Username,
				DisplayName: c.DisplayName,
				IsTyping:    true,
			},
		}
		data, _ := json.Marshal(wsMsg)
		c.Hub.broadcast <- &RoomMessage{
			RoomID:  c.RoomID,
			Message: data,
		}

	case model.WSTypeStopTyping:
		// Broadcast stop typing
		c.Hub.pubsubRepo.PublishTyping(ctx, c.RoomID, &model.TypingPayload{
			UserID:      c.UserID.String(),
			Username:    c.Username,
			DisplayName: c.DisplayName,
			IsTyping:    false,
		})

		wsMsg := model.WSMessage{
			Type: model.WSTypeStopTyping,
			Payload: model.TypingPayload{
				UserID:      c.UserID.String(),
				Username:    c.Username,
				DisplayName: c.DisplayName,
				IsTyping:    false,
			},
		}
		data, _ := json.Marshal(wsMsg)
		c.Hub.broadcast <- &RoomMessage{
			RoomID:  c.RoomID,
			Message: data,
		}
	}
}
