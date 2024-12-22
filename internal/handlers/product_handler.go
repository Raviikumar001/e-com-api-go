package handlers

import (
	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/models"
	"github.com/gofiber/fiber/v2"
)

type CreateProductRequest struct {
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
}

func CreateProduct(c *fiber.Ctx) error {
    user := c.Locals("user").(*models.User)
    
    var req CreateProductRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    product := models.Product{
        Name:         req.Name,
        Description:  req.Description,
        Price:        req.Price,
        Stock:        req.Stock,
        WholesalerID: user.ID,
    }

    result := database.DB.Create(&product)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not create product",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(product)
}

func GetProducts(c *fiber.Ctx) error {
    var products []models.Product
    result := database.DB.Find(&products)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not fetch products",
        })
    }

    return c.JSON(products)
}


// internal/handlers/product_handler.go
func UpdateProduct(c *fiber.Ctx) error {
    user := c.Locals("user").(*models.User)
    productID := c.Params("id")

    var product models.Product
    if err := database.DB.First(&product, productID).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Product not found",
        })
    }

    // Check if user is authorized to update this product
    if user.Role.Name == models.WholesalerRole && product.WholesalerID != user.ID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Not authorized to update this product",
        })
    }

    if user.Role.Name == models.SellerRole && product.SellerID != user.ID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Not authorized to update this product",
        })
    }

    var updateData map[string]interface{}
    if err := c.BodyParser(&updateData); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    if err := database.DB.Model(&product).Updates(updateData).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not update product",
        })
    }

    return c.JSON(product)
}