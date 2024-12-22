package handlers

import (
	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/models"
	"github.com/Raviikumar001/e-com-api-go/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type RegisterRequest struct {
    Email     string `json:"email"`
    Password  string `json:"password"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    RoleID    uint   `json:"role_id"`
}

func Login(c *fiber.Ctx) error {
    var req LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    var user models.User
    result := database.DB.Where("email = ?", req.Email).First(&user)
    if result.Error != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid credentials",
        })
    }

    if !utils.CheckPassword(req.Password, user.Password) {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid credentials",
        })
    }

    token, err := utils.GenerateJWT(user.ID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not generate token",
        })
    }

    return c.JSON(fiber.Map{
        "token": token,
    })
}

func Register(c *fiber.Ctx) error {
    var req RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    user := models.User{
        Email:     req.Email,
        Password:  req.Password,
        FirstName: req.FirstName,
        LastName:  req.LastName,
        RoleID:    req.RoleID,
    }

    result := database.DB.Create(&user)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not create user",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "User created successfully",
    })
}