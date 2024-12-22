// internal/handlers/storefront_handler.go
package handlers

import (
	"fmt"
	"math"

	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/models"
	"github.com/gofiber/fiber/v2"
)

type CreateStorefrontRequest struct {
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Theme       string            `json:"theme"`
    Domain      string            `json:"domain"`
    Settings    models.Settings   `json:"settings"`
}

func CreateStorefront(c *fiber.Ctx) error {
    user := c.Locals("user").(*models.User)
    
    var storefront models.Storefront
    if err := c.BodyParser(&storefront); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    storefront.SellerID = user.ID
    if err := database.DB.Create(&storefront).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not create storefront",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(storefront)
}


// func CreateStorefront(c *fiber.Ctx) error {
//     user := c.Locals("user").(*models.User)
    
//     if user.Role.Name != models.SellerRole {
//         return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
//             "error": "Only sellers can create storefronts",
//         })
//     }

//     var req CreateStorefrontRequest
//     if err := c.BodyParser(&req); err != nil {
//         return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//             "error": "Invalid request body",
//         })
//     }

//     storefront := models.Storefront{
//         Name:        req.Name,
//         Description: req.Description,
//         SellerID:    user.ID,
//         Theme:       req.Theme,
//         Domain:      req.Domain,
//         Settings:    req.Settings,
//     }

//     if err := database.DB.Create(&storefront).Error; err != nil {
//         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//             "error": "Could not create storefront",
//         })
//     }

//     return c.Status(fiber.StatusCreated).JSON(storefront)
// }

func GetStorefront(c *fiber.Ctx) error {
    user := c.Locals("user").(*models.User)
    
    // Pagination parameters
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 10)
    offset := (page - 1) * limit
    
    // Sorting parameters
    sortBy := c.Query("sort_by", "created_at")
    sortOrder := c.Query("sort_order", "desc")

    var storefronts []models.Storefront
    var total int64

    // Validate sort parameters
    allowedSortFields := map[string]bool{
        "created_at": true,
        "name":      true,
        "domain":    true,
    }
    if !allowedSortFields[sortBy] {
        sortBy = "created_at"
    }
    if sortOrder != "asc" && sortOrder != "desc" {
        sortOrder = "desc"
    }

    // Get total count
    if err := database.DB.Model(&models.Storefront{}).Where("seller_id = ?", user.ID).Count(&total).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error counting storefronts",
        })
    }

    // Get paginated and sorted storefronts
    if err := database.DB.Where("seller_id = ?", user.ID).
        Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
        Limit(limit).
        Offset(offset).
        Find(&storefronts).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error fetching storefronts",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Storefronts retrieved successfully",
        "storefronts": storefronts,
        "pagination": fiber.Map{
            "current_page": page,
            "per_page":    limit,
            "total":       total,
            "total_pages": int(math.Ceil(float64(total) / float64(limit))),
        },
    })
}


// func GetStorefront(c *fiber.Ctx) error {
//     user := c.Locals("user").(*models.User)
    
//     var storefront models.Storefront
//     if err := database.DB.Where("seller_id = ?", user.ID).First(&storefront).Error; err != nil {
//         return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
//             "error": "Storefront not found",
//         })
//     }

//     return c.JSON(storefront)
// }

// func GetStorefront(c *fiber.Ctx) error {
//     id := c.Params("id")
    
//     var storefront models.Storefront
//     if err := database.DB.First(&storefront, id).Error; err != nil {
//         return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
//             "error": "Storefront not found",
//         })
//     }

//     return c.JSON(storefront)
// }

func UpdateStorefront(c *fiber.Ctx) error {
    user := c.Locals("user").(*models.User)
    id := c.Params("id")

    var storefront models.Storefront
    if err := database.DB.First(&storefront, id).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Storefront not found",
        })
    }

    if storefront.SellerID != user.ID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Not authorized to update this storefront",
        })
    }

    var updateData CreateStorefrontRequest
    if err := c.BodyParser(&updateData); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    storefront.Name = updateData.Name
    storefront.Description = updateData.Description
    storefront.Theme = updateData.Theme
    storefront.Domain = updateData.Domain
    storefront.Settings = updateData.Settings

    if err := database.DB.Save(&storefront).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not update storefront",
        })
    }

    return c.JSON(storefront)
}

func DeleteStorefront(c *fiber.Ctx) error {
    user := c.Locals("user").(*models.User)
    id := c.Params("id")

    var storefront models.Storefront
    if err := database.DB.First(&storefront, id).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Storefront not found",
        })
    }

    if storefront.SellerID != user.ID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Not authorized to delete this storefront",
        })
    }

    if err := database.DB.Delete(&storefront).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not delete storefront",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Storefront deleted successfully",
    })
}