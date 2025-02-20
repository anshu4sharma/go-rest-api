package app

import (
	"log"
	"github.com/anshu4sharma/go-rest-api/internal/config"
	"github.com/anshu4sharma/go-rest-api/internal/handlers"
	"github.com/anshu4sharma/go-rest-api/internal/middleware"
	"github.com/anshu4sharma/go-rest-api/internal/models"
	"github.com/anshu4sharma/go-rest-api/internal/repository"
	"github.com/anshu4sharma/go-rest-api/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	router *fiber.App
	db     *gorm.DB
}

func NewApp() *App {
	app := &App{}
	app.initialize()
	return app
}

func (a *App) initialize() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize database
	a.db, err = gorm.Open(mysql.Open(cfg.DBConfig.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate models
	if err := a.db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(a.db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTConfig.Secret, cfg.JWTConfig.Expiration)
	googleAuthService := services.NewGoogleAuthService(cfg.GoogleConfig, userRepo, authService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, googleAuthService)

	// Initialize router with config
	a.router = fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Add middleware
	a.router.Use(recover.New())
	a.router.Use(logger.New())

	// Public routes
	auth := a.router.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Get("/google/login", authHandler.GoogleLogin)
	auth.Get("/google/callback", authHandler.GoogleCallback)

	// Protected routes
	protected := a.router.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTConfig.Secret))

	// Profile route
	protected.Get("/profile", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uint)
		return c.JSON(fiber.Map{"user_id": userID})
	})

	// Admin routes
	// admin := protected.Group("/admin")
	// admin.Use(middleware.RoleMiddleware("admin"))
	// admin.Get("/users", func(c *fiber.Ctx) error {
	// 	users, err := userRepo.GetAllUsers()
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	// 	}
	// 	return c.JSON(users)
	// })
}

func (a *App) Start(addr string) error {
	return a.router.Listen(addr)
}
