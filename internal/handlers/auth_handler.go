// internal/handlers/auth_handler.go
package handlers

import (
	"log"

	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/models"
	"github.com/Raviikumar001/e-com-api-go/internal/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {
    var input struct {
        Email     string `json:"email"`
        Password  string `json:"password"`
        FirstName string `json:"first_name"`
        LastName  string `json:"last_name"`
        RoleID    uint   `json:"role_id"`
    }

    if err := c.BodyParser(&input); err != nil {
        log.Printf("Error parsing registration input: %v", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }


    log.Printf("Attempting to register user with email: %s", input.Email)

 
    var existingUser models.User
    if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
        return c.Status(fiber.StatusConflict).JSON(fiber.Map{
            "error": "Email already registered",
        })
    }


    user := models.User{
        Email:     input.Email,
        Password:  input.Password, 
        FirstName: input.FirstName,
        LastName:  input.LastName,
        RoleID:    input.RoleID,
    }

    if err := database.DB.Create(&user).Error; err != nil {
        log.Printf("Error creating user: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not create user",
        })
    }


    if err := database.DB.Preload("Role").First(&user, user.ID).Error; err != nil {
        log.Printf("Error loading role: %v", err)
    }

    log.Printf("Successfully registered user with ID: %d", user.ID)

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "User registered successfully",
        "user": fiber.Map{
            "id":         user.ID,
            "email":      user.Email,
            "first_name": user.FirstName,
            "last_name":  user.LastName,
            "role_id":    user.RoleID,
        },
    })
}

func Login(c *fiber.Ctx) error {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.BodyParser(&input); err != nil {
        log.Printf("Error parsing login input: %v", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    log.Printf("Login attempt for email: %s", input.Email)

    var user models.User
    if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid email or password",
            })
        }
        log.Printf("Database error during login: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Login error",
        })
    }


    if !utils.CheckPassword(input.Password, user.Password) {
        log.Printf("Invalid password attempt for user: %s", input.Email)
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid email or password",
        })
    }

    
    token, err := utils.GenerateToken(user.ID)
    if err != nil {
        log.Printf("Error generating token: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not generate token",
        })
    }

    log.Printf("Successful login for user ID: %d", user.ID)

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Login successful",
        "token":   token,
        "user": fiber.Map{
            "id":         user.ID,
            "email":      user.Email,
            "first_name": user.FirstName,
            "last_name":  user.LastName,
            "role_id":    user.RoleID,
        },
    })
}

