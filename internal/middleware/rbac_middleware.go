// internal/middleware/rbac_middleware.go
package middleware

import (
	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/models"
	"github.com/gofiber/fiber/v2"
)

func RBACMiddleware(requiredPermission string) fiber.Handler {
    return func(c *fiber.Ctx) error {

        user, ok := c.Locals("user").(*models.User)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Unauthorized",
            })
        }

          // If user is admin, allow access immediately
          if user.Role.Name == models.AdminRole {
            return c.Next()
        }

        // Load user role with permissions
        var role models.Role
        if err := database.DB.Preload("Permissions").First(&role, user.RoleID).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Error loading user role",
            })
        }

        // Check if the role has the required permission
        hasPermission := false
        for _, permission := range role.Permissions {
            if permission.Name == requiredPermission {
                hasPermission = true
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