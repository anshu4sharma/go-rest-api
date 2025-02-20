package handlers

import (
	"github.com/anshu4sharma/go-rest-api/internal/models"
	"github.com/anshu4sharma/go-rest-api/internal/services"
	"github.com/anshu4sharma/go-rest-api/internal/validation"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService       *services.AuthService
	googleAuthService *services.GoogleAuthService
}

func NewAuthHandler(authService *services.AuthService, googleAuthService *services.GoogleAuthService) *AuthHandler {
	return &AuthHandler{
		authService:       authService,
		googleAuthService: googleAuthService,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate request
	if err := validation.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate request
	if err := validation.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := h.authService.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(token)
}

func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	state := "random-state" // In production, generate a random state and store it in session/cookie
	url := h.googleAuthService.GetAuthURL(state)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"url": url})
}

func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Code is required"})
	}

	token, err := h.googleAuthService.HandleCallback(code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(token)
}
