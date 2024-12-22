// cmd/main.go
package main

import (
	"log"
	"os"

	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {

    if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found")
    }


    database.InitDB()

    app := fiber.New()


    app.Use(logger.New())
    app.Use(cors.New())


    routes.SetupRoutes(app)

    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    log.Printf("Server starting on port %s", port)
    log.Fatal(app.Listen(":" + port))
}

