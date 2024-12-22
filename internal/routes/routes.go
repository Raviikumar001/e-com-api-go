package routes

import (
	"github.com/Raviikumar001/e-com-api-go/internal/handlers"
	"github.com/Raviikumar001/e-com-api-go/internal/middleware"
	"github.com/Raviikumar001/e-com-api-go/internal/models"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
    // Public routes
    auth := app.Group("/auth")
    auth.Post("/login", handlers.Login)
    auth.Post("/register", handlers.Register)

    // Protected routes
    api := app.Group("/api", middleware.AuthMiddleware())

    // Product routes
    products := api.Group("/products")
    products.Post("/", middleware.RBACMiddleware(models.CreateProduct), handlers.CreateProduct)
    products.Get("/", handlers.GetProducts)

    // Add more routes as needed...
}