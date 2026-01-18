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

	// Initialize GORM (parallel to pgx for gradual migration)
	gormDB, err := database.NewGormDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, "disable")
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL via GORM: %v", err)
	}
	defer gormDB.Close()
	log.Println("✅ Connected to PostgreSQL via GORM")

	// Auto-migrate new tables (File, Reaction, DM fields)
	// Note: Running migrations on existing tables is safe - GORM only adds missing columns
	// For production, you should use separate migration files
	// gormDB.AutoMigrate(&model.File{}, &model.Reaction{})

	// Initialize repositories (keeping old pgx repos for backward compatibility)
	userRepo := repository.NewUserRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	messageRepo := repository.NewMessageRepository(db, rdb)
	presenceRepo := repository.NewPresenceRepository(rdb)
	pubsubRepo := repository.NewPubSubRepository(rdb)

	// Initialize GORM-based repositories for new features
	gormUserRepo := repository.NewGormUserRepository(gormDB)
	gormRoomRepo := repository.NewGormRoomRepository(gormDB)
	gormFileRepo := repository.NewGormFileRepository(gormDB)

	// Initialize services
	chatService := service.NewChatService(messageRepo, pubsubRepo, presenceRepo)
	presenceService := service.NewPresenceService(presenceRepo)

	// Initialize Global WebSocket hub for homepage updates
	globalHub := ws.NewGlobalHub(roomRepo)
	go globalHub.Run()

	// Initialize WebSocket hub for chat rooms
	hub := ws.NewHub(chatService, presenceService, pubsubRepo, globalHub)
	go hub.Run()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:   "Isan Chat - หนองบัวลำภู 🏯",
		BodyLimit: 15 * 1024 * 1024, // 15MB for file uploads
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,Upgrade,Connection,Sec-WebSocket-Key,Sec-WebSocket-Version,Sec-WebSocket-Extensions",
		AllowCredentials: true,
	}))

	// Static file serving for uploads
	app.Static("/uploads", "./uploads")

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
	api.Get("/users/search", func(c *fiber.Ctx) error {
		query := c.Query("q")
		if query == "" {
			return c.JSON([]interface{}{})
		}
		users, err := gormUserRepo.Search(c.Context(), query, 20)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(users)
	})

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

	// File upload routes
	fileHandler := handler.NewFileHandler(gormFileRepo, "./uploads", "http://localhost:"+cfg.ServerPort)
	api.Post("/upload", fileHandler.Upload)
	api.Post("/upload/multiple", fileHandler.UploadMultiple)
	api.Get("/files/:id", fileHandler.GetFile)
	api.Delete("/files/:id", fileHandler.Delete)
	app.Get("/uploads/:filename", fileHandler.ServeFile)

	// DM routes
	dmHandler := handler.NewDMHandler(gormRoomRepo, gormUserRepo)
	api.Post("/dm/start", dmHandler.StartDM)
	api.Get("/dm", dmHandler.ListDMs)
	api.Get("/dm/:targetUserId", dmHandler.GetDM)

	// Initialize additional GORM repositories for Phase 2.4-2.6
	gormMessageRepo := repository.NewGormMessageRepository(gormDB)
	gormReactionRepo := repository.NewGormReactionRepository(gormDB)
	gormNotificationRepo := repository.NewGormNotificationRepository(gormDB)

	// Reaction routes (Phase 2.4)
	reactionHandler := handler.NewReactionHandler(gormReactionRepo, gormMessageRepo)
	api.Post("/reactions", reactionHandler.AddReaction)
	api.Delete("/reactions", reactionHandler.RemoveReaction)
	api.Post("/messages/:messageId/reactions/:emoji", reactionHandler.ToggleReaction)
	api.Get("/messages/:messageId/reactions", reactionHandler.GetReactions)
	api.Get("/reactions/emojis", reactionHandler.GetCommonEmojis)

	// Search routes (Phase 2.5)
	searchHandler := handler.NewSearchHandler(gormMessageRepo, gormUserRepo, gormRoomRepo)
	api.Get("/rooms/:roomId/search", searchHandler.SearchMessages)
	api.Get("/rooms/:roomId/messages/before", searchHandler.GetMessagesBefore)
	api.Get("/search/users", searchHandler.SearchUsers)
	api.Get("/search", searchHandler.GlobalSearch)

	// Notification routes (Phase 2.6)
	notificationHandler := handler.NewNotificationHandler(gormNotificationRepo)
	api.Get("/notifications", notificationHandler.GetNotifications)
	api.Get("/notifications/unread", notificationHandler.GetUnreadNotifications)
	api.Get("/notifications/count", notificationHandler.GetUnreadCount)
	api.Post("/notifications/:id/read", notificationHandler.MarkAsRead)
	api.Post("/notifications/read-all", notificationHandler.MarkAllAsRead)
	api.Delete("/notifications/:id", notificationHandler.Delete)
	api.Delete("/notifications", notificationHandler.DeleteAll)
	api.Post("/push/subscribe", notificationHandler.SubscribePush)
	api.Delete("/push/unsubscribe", notificationHandler.UnsubscribePush)

	// Global WebSocket route for homepage real-time updates (MUST be before /ws/:roomId)
	app.Get("/ws/global", func(c *fiber.Ctx) error {
		log.Printf("🌐 GET /ws/global - Global WebSocket request")
		if !websocket.IsWebSocketUpgrade(c) {
			return c.SendStatus(fiber.StatusUpgradeRequired)
		}

		return websocket.New(func(conn *websocket.Conn) {
			log.Printf("🌐 Global WebSocket connection!")
			globalHub.HandleGlobalWebSocket(conn)
		})(c)
	})

	// WebSocket route for chat rooms
	app.Get("/ws/:roomId", func(c *fiber.Ctx) error {
		log.Printf("🎯 GET /ws/%s - Headers: %v", c.Params("roomId"), c.GetReqHeaders())
		log.Printf("🎯 Upgrade header: %s", c.Get("Upgrade"))
		log.Printf("🎯 Connection header: %s", c.Get("Connection"))

		if !websocket.IsWebSocketUpgrade(c) {
			log.Printf("❌ NOT a WebSocket upgrade request")
			return c.SendStatus(fiber.StatusUpgradeRequired)
		}

		log.Printf("✅ IS a WebSocket upgrade request")
		return websocket.New(func(conn *websocket.Conn) {
			log.Printf("🔌 WebSocket handler called!")
			hub.HandleWebSocket(conn)
		})(c)
	})

	// Graceful shutdown
	go func() {
		addr := cfg.ServerHost + ":" + cfg.ServerPort
		// On Windows, "localhost" may resolve to ::1 first. If we bind IPv4-only (0.0.0.0),
		// WebSocket connections can fail with a generic client-side error. Binding ":port"
		// enables dual-stack on most systems.
		if cfg.ServerHost == "" || cfg.ServerHost == "0.0.0.0" {
			addr = ":" + cfg.ServerPort
		}
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
