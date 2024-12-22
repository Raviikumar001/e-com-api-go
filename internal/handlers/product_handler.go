package handlers

import (
	"math"
	"strconv"

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

    var product models.Product
    if err := c.BodyParser(&product); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    // Validate required fields
    if product.Name == "" || product.Price <= 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Name and price are required, and price must be greater than 0",
        })
    }

    // Set the appropriate ID based on user role
    if user.Role.Name == models.WholesalerRole {
        wholesalerID := user.ID
        product.WholesalerID = &wholesalerID
        product.SellerID = nil
    } else if user.Role.Name == models.SellerRole {
        sellerID := user.ID
        product.SellerID = &sellerID
        product.WholesalerID = nil
    } else {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Only sellers and wholesalers can create products",
        })
    }

    // Set default values
    product.IsPublished = false // Products are unpublished by default

    // Save the product
    if err := database.DB.Create(&product).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not create product",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Product created successfully",
        "product": product,
    })
}




func GetProducts(c *fiber.Ctx) error {
    user := c.Locals("user").(*models.User)
    
    // Pagination parameters
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 10)
    offset := (page - 1) * limit

    query := database.DB.Model(&models.Product{})

    // Apply role-based filtering
    switch user.Role.Name {
    case models.WholesalerRole:
        query = query.Where("wholesaler_id = ?", user.ID)
    case models.SellerRole:
        query = query.Where("seller_id = ?", user.ID)
    case models.CustomerRole:
        query = query.Where("is_published = ?", true)
    }

    // Get total count
    var total int64
    if err := query.Count(&total).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error counting products",
        })
    }

    // Get products
    var products []models.Product
    if err := query.
        Limit(limit).
        Offset(offset).
        Find(&products).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not fetch products",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Products retrieved successfully",
        "products": products,
        "pagination": fiber.Map{
            "current_page": page,
            "per_page":    limit,
            "total":       total,
            "total_pages": int(math.Ceil(float64(total) / float64(limit))),
        },
    })
}






func UpdateProduct(c *fiber.Ctx) error {
    user := c.Locals("user").(*models.User)

    // Get product ID from params
    productID, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid product ID",
        })
    }

    // Find existing product
    var existingProduct models.Product
    if err := database.DB.First(&existingProduct, productID).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Product not found",
        })
    }

    // Check ownership based on role
    if user.Role.Name == models.SellerRole && (existingProduct.SellerID == nil || *existingProduct.SellerID != user.ID) {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "You don't have permission to update this product",
        })
    }

    if user.Role.Name == models.WholesalerRole && (existingProduct.WholesalerID == nil || *existingProduct.WholesalerID != user.ID) {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "You don't have permission to update this product",
        })
    }

    // Parse update data
    var updateData models.Product
    if err := c.BodyParser(&updateData); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    // Validate update data
    if updateData.Price <= 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Price must be greater than 0",
        })
    }

    // Prevent updating sensitive fields
    updateData.ID = existingProduct.ID
    updateData.SellerID = existingProduct.SellerID
    updateData.WholesalerID = existingProduct.WholesalerID
    updateData.CreatedAt = existingProduct.CreatedAt

    // Update only allowed fields
    if err := database.DB.Model(&existingProduct).Updates(map[string]interface{}{
        "name":         updateData.Name,
        "description":  updateData.Description,
        "price":       updateData.Price,
        "stock":       updateData.Stock,
        "category":    updateData.Category,
        "image_url":   updateData.ImageURL,
        "is_published": updateData.IsPublished,
    }).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not update product",
        })
    }

    // Refresh the product data
    if err := database.DB.First(&existingProduct, productID).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not fetch updated product",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Product updated successfully",
        "product": existingProduct,
    })
}

func GetProductDetails(c *fiber.Ctx) error {
    user := c.Locals("user").(*models.User)

    // Get product ID from params
    productID, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid product ID",
        })
    }

    // Find product with related data
    var product models.Product
    query := database.DB

    // If admin, include seller and wholesaler details
    if user.Role.Name == models.AdminRole {
        query = query.Preload("Seller").Preload("Wholesaler")
    }

    if err := query.First(&product, productID).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Product not found",
        })
    }

    // Handle role-based access and response
    switch user.Role.Name {
    case models.CustomerRole:
        // Customers can only view published products
        if !product.IsPublished {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Product not available",
            })
        }
        // Return limited product details for customers
        return c.JSON(fiber.Map{
            "message": "Product details retrieved successfully",
            "product": fiber.Map{
                "id":          product.ID,
                "name":        product.Name,
                "description": product.Description,
                "price":       product.Price,
                "category":    product.Category,
                "image_url":   product.ImageURL,
            },
        })

    case models.SellerRole:
        // Sellers can only view their own products
        if product.SellerID == nil || *product.SellerID != user.ID {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "You don't have permission to view this product",
            })
        }

    case models.WholesalerRole:
        // Wholesalers can only view their own products
        if product.WholesalerID == nil || *product.WholesalerID != user.ID {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "You don't have permission to view this product",
            })
        }

    case models.AdminRole:
        // Admins can view all products with full details
        break

    default:
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Invalid role",
        })
    }

    // Prepare full response for admin, seller, or wholesaler
    response := fiber.Map{
        "id":           product.ID,
        "name":         product.Name,
        "description":  product.Description,
        "price":        product.Price,
        "stock":        product.Stock,
        "category":     product.Category,
        "image_url":    product.ImageURL,
        "is_published": product.IsPublished,
        "created_at":   product.CreatedAt,
        "updated_at":   product.UpdatedAt,
    }

    // Add seller/wholesaler info for admins
    if user.Role.Name == models.AdminRole {
        // Check if Seller exists before adding to response
        if product.Seller != nil {
            response["seller"] = fiber.Map{
                "id":         product.Seller.ID,
                "email":      product.Seller.Email,
                "first_name": product.Seller.FirstName,
                "last_name":  product.Seller.LastName,
            }
        }
        
        // Check if Wholesaler exists before adding to response
        if product.Wholesaler != nil {
            response["wholesaler"] = fiber.Map{
                "id":         product.Wholesaler.ID,
                "email":      product.Wholesaler.Email,
                "first_name": product.Wholesaler.FirstName,
                "last_name":  product.Wholesaler.LastName,
            }
        }
    }

    return c.JSON(fiber.Map{
        "message": "Product details retrieved successfully",
        "product": response,
    })
}