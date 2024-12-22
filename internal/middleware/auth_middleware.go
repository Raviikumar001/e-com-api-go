package middleware

import (
	"os"

	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := c.Get("Authorization")
        if token == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Authorization header required",
            })
        }

        // Remove Bearer prefix
        tokenString := token[7:] // Remove "Bearer "

        // Parse token
        claims := jwt.MapClaims{}
        parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !parsedToken.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid token",
            })
        }

        // Get user from database
        var user models.User
        result := database.DB.Preload("Role.Permissions").First(&user, claims["user_id"])
        if result.Error != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "User not found",
            })
        }

        // Set user in context
        c.Locals("user", &user)
        return c.Next()
    }
}