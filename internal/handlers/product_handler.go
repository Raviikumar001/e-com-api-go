package handlers

import (
	"errors"
	"math"
	"strconv"

	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

   
    if product.Name == "" || product.Price <= 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Name and price are required, and price must be greater than 0",
        })
    }

    // Admin can create products for any seller or wholesaler
    if user.Role.Name == models.AdminRole {
       
    } else {
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
    }

    
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
    
    
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 10)
    offset := (page - 1) * limit

    query := database.DB.Model(&models.Product{})

   
   if user.Role.Name != models.AdminRole {
    switch user.Role.Name {
    case models.WholesalerRole:
        query = query.Where("wholesaler_id = ?", user.ID)
    case models.SellerRole:
        query = query.Where("seller_id = ?", user.ID)
    case models.CustomerRole:
        query = query.Where("is_published = ?", true)
    }
}
    
    var total int64
    if err := query.Count(&total).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error counting products",
        })
    }

    
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

    
    productID, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid product ID",
        })
    }

    
    var existingProduct models.Product
    if err := database.DB.First(&existingProduct, productID).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Product not found",
        })
    }

      
      if user.Role.Name != models.AdminRole {
        // Check ownership based on role
        if user.Role.Name == models.SellerRole && 
           (existingProduct.SellerID == nil || *existingProduct.SellerID != user.ID) {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "You don't have permission to update this product",
            })
        }

        if user.Role.Name == models.WholesalerRole && 
           (existingProduct.WholesalerID == nil || *existingProduct.WholesalerID != user.ID) {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "You don't have permission to update this product",
            })
        }
    }

    
    var updateData models.Product
    if err := c.BodyParser(&updateData); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    
    if updateData.Price <= 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Price must be greater than 0",
        })
    }

    
    updateData.ID = existingProduct.ID
    updateData.SellerID = existingProduct.SellerID
    updateData.WholesalerID = existingProduct.WholesalerID
    updateData.CreatedAt = existingProduct.CreatedAt

    
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

    
    productID, err := strconv.ParseUint(c.Params("id"), 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid product ID",
        })
    }

   
    var product models.Product
    query := database.DB

    
    query = query.Preload("Seller").Preload("Wholesaler")

    if err := query.First(&product, productID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "Product not found",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error fetching product",
        })
    }

    // Base response for all roles
    baseResponse := fiber.Map{
        "id":          product.ID,
        "name":        product.Name,
        "description": product.Description,
        "price":       product.Price,
        "category":    product.Category,
        "image_url":   product.ImageURL,
    }

    // Handle role-based access and response
    switch user.Role.Name {
    case models.CustomerRole:
        if !product.IsPublished {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Product not available",
            })
        }
        return c.JSON(fiber.Map{
            "message": "Product details retrieved successfully",
            "product": baseResponse,
        })

    case models.SellerRole:
        if product.SellerID == nil || *product.SellerID != user.ID {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "You don't have permission to view this product",
            })
        }
        // Add seller-specific fields
        baseResponse["stock"] = product.Stock
        baseResponse["is_published"] = product.IsPublished
        baseResponse["created_at"] = product.CreatedAt
        baseResponse["updated_at"] = product.UpdatedAt

    case models.WholesalerRole:
        if product.WholesalerID == nil || *product.WholesalerID != user.ID {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "You don't have permission to view this product",
            })
        }
        // Add wholesaler-specific fields
        baseResponse["stock"] = product.Stock
        baseResponse["cost_price"] = product.CostPrice
        baseResponse["created_at"] = product.CreatedAt
        baseResponse["updated_at"] = product.UpdatedAt

    case models.AdminRole:
        // Add all fields for admin
        baseResponse["stock"] = product.Stock
        baseResponse["cost_price"] = product.CostPrice
        baseResponse["is_published"] = product.IsPublished
        baseResponse["created_at"] = product.CreatedAt
        baseResponse["updated_at"] = product.UpdatedAt
        
        if product.Seller != nil {
            baseResponse["seller"] = fiber.Map{
                "id":         product.Seller.ID,
                "email":      product.Seller.Email,
                "first_name": product.Seller.FirstName,
                "last_name":  product.Seller.LastName,
            }
        }
        
        if product.Wholesaler != nil {
            baseResponse["wholesaler"] = fiber.Map{
                "id":         product.Wholesaler.ID,
                "email":      product.Wholesaler.Email,
                "first_name": product.Wholesaler.FirstName,
                "last_name":  product.Wholesaler.LastName,
            }
        }

    default:
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Invalid role",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Product details retrieved successfully",
        "product": baseResponse,
    })
}