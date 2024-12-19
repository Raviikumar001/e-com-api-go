// main.go

package main

import (
	"log"
	"os"

	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/repositories"
	"github.com/Raviikumar001/e-com-api-go/internal/services"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize database
	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)

	// Initialize services
	jwtService := services.NewJWTService(os.Getenv("JWT_SECRET_KEY"), os.Getenv("JWT_ISSUER"))
	authService := services.NewAuthService(userRepo, jwtService)
	roleService := services.NewRoleService(*roleRepo, authService)
	permissionService := services.NewPermissionService(*permissionRepo, authService)

	// Create a new Fiber app
	app := fiber.New()

	// Register routes
	// We'll add routes here later

	// Start the server
	app.Listen(":8080")
}
