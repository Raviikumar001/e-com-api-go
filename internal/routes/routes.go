// internal/routes/routes.go
package routes

import (
	"github.com/Raviikumar001/e-com-api-go/internal/handlers"
	"github.com/Raviikumar001/e-com-api-go/internal/middleware"
	"github.com/Raviikumar001/e-com-api-go/internal/models"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
    // Public Authentication routes
    auth := app.Group("/auth")
    auth.Post("/login", handlers.Login)
    auth.Post("/register", handlers.Register)

    // Protected routes - requires authentication
    api := app.Group("/api", middleware.AuthMiddleware())

    // 1. Module Level Access Example: Web Builder (Sellers only)
    // webBuilder := api.Group("/web-builder", middleware.RBACMiddleware(models.AccessWebBuilder))
    // webBuilder.Post("/storefront", handlers.CreateStorefront)
    // webBuilder.Get("/storefront", handlers.GetStorefront)
    webBuilder := api.Group("/web-builder")
    webBuilder.Use(middleware.AuthMiddleware())
    webBuilder.Use(middleware.RBACMiddleware("access_web_builder"))
    webBuilder.Post("/storefront", handlers.CreateStorefront)
    webBuilder.Get("/storefront", handlers.GetStorefront)

    // 2. Endpoint Level Access Example: Product Management
    products := api.Group("/products")
    // Public product viewing (all authenticated users)
    products.Get("/", handlers.GetProducts)
    // Create product (sellers and wholesalers only)
    products.Post("/", middleware.RBACMiddleware(models.CreateProduct), handlers.CreateProduct)

    // 3. Data Level Access Example: Product Details with Inventory
    products.Get("/:id/details", handlers.GetProductDetails) // Will return different data based on role
}