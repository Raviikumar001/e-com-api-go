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
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found")
    }

    // Initialize database
    database.InitDB()

    app := fiber.New()

    // Middleware
    app.Use(logger.New())
    app.Use(cors.New())

    // Setup routes
    routes.SetupRoutes(app)

    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    log.Printf("Server starting on port %s", port)
    log.Fatal(app.Listen(":" + port))
}


// package main

// import (
// 	"log"
// 	"os"

// 	"github.com/Raviikumar001/e-com-api-go/internal/database"
// 	"github.com/Raviikumar001/e-com-api-go/internal/middleware"
// 	"github.com/Raviikumar001/e-com-api-go/internal/routes"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/fiber/v2/middleware/cors"
// 	"github.com/gofiber/fiber/v2/middleware/logger"
// 	"github.com/gofiber/fiber/v2/middleware/recover"
// 	"github.com/joho/godotenv"
// )

// func main() {
//     // Load environment variables
//     if err := godotenv.Load(); err != nil {
//         log.Println("No .env file found")
//     }

//     // Initialize database
//     database.InitDB()

//     app := fiber.New(fiber.Config{
//         ErrorHandler: middleware.ErrorHandler,
//     })

//     // Middleware
//     app.Use(recover.New())
//     app.Use(logger.New())
//     app.Use(cors.New())

//     // Setup routes
//     routes.SetupRoutes(app)

//     // Start server
//     port := os.Getenv("PORT")
//     if port == "" {
//         port = "3000"
//     }
    
//     log.Fatal(app.Listen(":" + port))
// }