// internal/controllers/user_controller.go

package controllers

import (
	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/services"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService *services.UserService
	authService *services.AuthService
}

func NewUserController(userService *services.UserService, authService *services.AuthService) *UserController {
	return &UserController{userService: userService, authService: authService}
}

func (uc *UserController) Register(c *fiber.Ctx) error {
	var user database.User
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid request body"})
	}

	err = uc.userService.Create(c, user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to create user"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "User created successfully"})
}

func (uc *UserController) Login(c *fiber.Ctx) error {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := c.BodyParser(&credentials)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid request body"})
	}

	token, err := uc.authService.Login(c, credentials.Username, credentials.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Invalid username or password"})
	}

	return c.Status(200).JSON(fiber.Map{"token": token})
}

func (uc *UserController) GetProfile(c *fiber.Ctx) error {
	claims, err := uc.authService.ValidateToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Invalid token"})
	}

	user, err := uc.userService.FindByID(c, claims.UserID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "User not found"})
	}

	return c.Status(200).JSON(user)
}
