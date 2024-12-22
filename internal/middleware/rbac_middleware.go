package middleware

import (
	"github.com/Raviikumar001/e-com-api-go/internal/models"
	"github.com/gofiber/fiber/v2"
)

func RBACMiddleware(requiredPermissions ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user := c.Locals("user").(*models.User)
        
        // Check if user has required permissions
        hasPermission := false
        for _, rolePermission := range user.Role.Permissions {
            for _, requiredPermission := range requiredPermissions {
                if rolePermission.Name == requiredPermission {
                    hasPermission = true
                    break
                }
            }
            if hasPermission {
                break
            }
        }

        if !hasPermission {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Insufficient permissions",
            })
        }

        return c.Next()
    }
}