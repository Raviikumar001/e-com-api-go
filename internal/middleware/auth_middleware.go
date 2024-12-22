// package middleware

// import (
// 	"os"

// 	"github.com/Raviikumar001/e-com-api-go/internal/database"
// 	"github.com/Raviikumar001/e-com-api-go/internal/models"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/golang-jwt/jwt/v4"
// )

// func AuthMiddleware() fiber.Handler {
//     return func(c *fiber.Ctx) error {
//         token := c.Get("Authorization")
//         if token == "" {
//             return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
//                 "error": "Authorization header required",
//             })
//         }

//         // Remove Bearer prefix
//         tokenString := token[7:] // Remove "Bearer "

//         // Parse token
//         claims := jwt.MapClaims{}
//         parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
//             return []byte(os.Getenv("JWT_SECRET")), nil
//         })

//         if err != nil || !parsedToken.Valid {
//             return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
//                 "error": "Invalid token",
//             })
//         }

//         // Get user from database
//         var user models.User
//         result := database.DB.Preload("Role.Permissions").First(&user, claims["user_id"])
//         if result.Error != nil {
//             return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
//                 "error": "User not found",
//             })
//         }

//         // Set user in context
//         c.Locals("user", &user)
//         return c.Next()
//     }
// }

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
        // Get the Authorization header
        authHeader := c.Get("Authorization")
        
        // Check if the header is empty or doesn't start with "Bearer "
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Missing or invalid authorization token",
            })
        }

        // Extract the token
        token := strings.TrimPrefix(authHeader, "Bearer ")

        // Validate the token
        userID, err := utils.ValidateToken(token)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid token",
            })
        }

        // Get the user from the database
        var user models.User
        if err := database.DB.Preload("Role").First(&user, userID).Error; err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "User not found",
            })
        }

        // Set the user in the context
        c.Locals("user", &user)

        return c.Next()
    }
}