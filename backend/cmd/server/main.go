package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/khonE3/chat-backend/internal/config"
	"github.com/khonE3/chat-backend/internal/handler"
	"github.com/khonE3/chat-backend/internal/repository"
	"github.com/khonE3/chat-backend/internal/service"
	ws "github.com/khonE3/chat-backend/internal/websocket"
	"github.com/khonE3/chat-backend/pkg/database"
	redisclient "github.com/khonE3/chat-backend/pkg/redis"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize PostgreSQL
	db, err := database.NewPostgres(cfg.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()
	log.Println("✅ Connected to PostgreSQL")

	// Initialize Redis
	rdb, err := redisclient.NewRedis(cfg.RedisURL, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()
	log.Println("✅ Connected to Redis")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	messageRepo := repository.NewMessageRepository(db, rdb)
	presenceRepo := repository.NewPresenceRepository(rdb)
	pubsubRepo := repository.NewPubSubRepository(rdb)

	// Initialize services
	chatService := service.NewChatService(messageRepo, pubsubRepo, presenceRepo)
	presenceService := service.NewPresenceService(presenceRepo)

	// Initialize WebSocket hub
	hub := ws.NewHub(chatService, presenceService, pubsubRepo)
	go hub.Run()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Isan Chat - หนองบัวลำภู 🏯",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "สวัสดีจากหนองบัวลำภู! 🎋",
		})
	})

	// API routes
	api := app.Group("/api")

	// User routes
	userHandler := handler.NewUserHandler(userRepo)
	api.Post("/users", userHandler.Create)
	api.Get("/users/:id", userHandler.GetByID)
	api.Get("/users/username/:username", userHandler.GetByUsername)

	// Room routes
	roomHandler := handler.NewRoomHandler(roomRepo, userRepo)
	api.Get("/rooms", roomHandler.List)
	api.Post("/rooms", roomHandler.Create)
	api.Get("/rooms/:id", roomHandler.GetByID)
	api.Post("/rooms/:id/join", roomHandler.Join)
	api.Get("/rooms/:id/members", roomHandler.GetMembers)
	api.Post("/rooms/:id/read", roomHandler.MarkAsRead)
	api.Get("/rooms/:id/unread", roomHandler.GetUnreadCount)

	// Message routes
	messageHandler := handler.NewMessageHandler(messageRepo)
	api.Get("/rooms/:id/messages", messageHandler.GetByRoom)

	// WebSocket route
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:roomId", websocket.New(func(c *websocket.Conn) {
		hub.HandleWebSocket(c)
	}))

	// Graceful shutdown
	go func() {
		addr := cfg.ServerHost + ":" + cfg.ServerPort
		log.Printf("🚀 Server starting on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("👋 Server exited")
}
