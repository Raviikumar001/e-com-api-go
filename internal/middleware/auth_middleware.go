// internal/middleware/auth_middleware.go
package middleware

import (
	"strings"

	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/models"
	"github.com/Raviikumar001/e-com-api-go/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {

        authHeader := c.Get("Authorization")
        
   
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Missing or invalid authorization token",
            })
        }


        token := strings.TrimPrefix(authHeader, "Bearer ")


        userID, err := utils.ValidateToken(token)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid token",
            })
        }


        var user models.User
        if err := database.DB.Preload("Role").First(&user, userID).Error; err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "User not found",
            })
        }


        c.Locals("user", &user)

        return c.Next()
    }
}