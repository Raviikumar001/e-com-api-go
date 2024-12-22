// internal/handlers/storefront_handler.go
package handlers

import (
	"fmt"
	"math"
	"strings"

	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/models"
	"github.com/gofiber/fiber/v2"
)

type CreateStorefrontRequest struct {
    Name        string          `json:"name" validate:"required,min=3,max=100"`
    Description string          `json:"description" validate:"max=500"`
    Theme       string          `json:"theme" validate:"required"`
    Domain      string          `json:"domain" validate:"required,hostname"`
    Settings    models.Settings `json:"settings"`
}


func validateStorefrontRequest(req CreateStorefrontRequest) error {
    if strings.TrimSpace(req.Name) == "" {
        return fmt.Errorf("name is required")
    }
    if strings.TrimSpace(req.Domain) == "" {
        return fmt.Errorf("domain is required")
    }
    if len(req.Name) > 100 {
        return fmt.Errorf("name must be less than 100 characters")
    }
    if len(req.Description) > 500 {
        return fmt.Errorf("description must be less than 500 characters")
    }
    return nil
}

func CreateStorefront(c *fiber.Ctx) error {
    user := c.Locals("user").(*models.User)
    
    var req CreateStorefrontRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }


    if err := validateStorefrontRequest(req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }


    var existingStore models.Storefront
    if err := database.DB.Where("domain = ?", req.Domain).First(&existingStore).Error; err == nil {
        return c.Status(fiber.StatusConflict).JSON(fiber.Map{
            "error": "Domain is already in use",
        })
    }

    storefront := models.Storefront{
        Name:        req.Name,
        Description: req.Description,
        Theme:       req.Theme,
        Domain:      req.Domain,
        Settings:    req.Settings,
        SellerID:    user.ID,
    }

    //transaction
    tx := database.DB.Begin()
    if err := tx.Create(&storefront).Error; err != nil {
        tx.Rollback()
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not create storefront",
        })
    }

    if err := tx.Commit().Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not commit transaction",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Storefront created successfully",
        "storefront": storefront,
    })
}

func GetStorefront(c *fiber.Ctx) error {
    user := c.Locals("user").(*models.User)
    

    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 10)
    if limit > 100 {
        limit = 100 // Maximum limit
    }
    offset := (page - 1) * limit
    

    sortBy := c.Query("sort_by", "created_at")
    sortOrder := c.Query("sort_order", "desc")
    

    search := c.Query("search", "")

    theme := c.Query("theme", "")
    status := c.Query("status", "")

    query := database.DB.Model(&models.Storefront{}).Where("seller_id = ?", user.ID)


    if search != "" {
        query = query.Where("name ILIKE ? OR description ILIKE ?", 
            "%"+search+"%", "%"+search+"%")
    }


    if theme != "" {
        query = query.Where("theme = ?", theme)
    }
    if status != "" {
        query = query.Where("status = ?", status)
    }


    var total int64
    if err := query.Count(&total).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Error counting storefronts",
        })
    }


    var storefronts []models.Storefront
    if err := query.
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

func UpdateStorefront(c *fiber.Ctx) error {
    user := c.Locals("user").(*models.User)
    id := c.Params("id")

    // transaction
    tx := database.DB.Begin()

    var storefront models.Storefront
    if err := tx.First(&storefront, id).Error; err != nil {
        tx.Rollback()
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Storefront not found",
        })
    }

    if storefront.SellerID != user.ID {
        tx.Rollback()
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Not authorized to update this storefront",
        })
    }

    var updateData CreateStorefrontRequest
    if err := c.BodyParser(&updateData); err != nil {
        tx.Rollback()
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

  
    if err := validateStorefrontRequest(updateData); err != nil {
        tx.Rollback()
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    if updateData.Domain != storefront.Domain {
        var existingStore models.Storefront
        if err := tx.Where("domain = ? AND id != ?", updateData.Domain, id).
            First(&existingStore).Error; err == nil {
            tx.Rollback()
            return c.Status(fiber.StatusConflict).JSON(fiber.Map{
                "error": "Domain is already in use",
            })
        }
    }


    storefront.Name = updateData.Name
    storefront.Description = updateData.Description
    storefront.Theme = updateData.Theme
    storefront.Domain = updateData.Domain
    storefront.Settings = updateData.Settings

    if err := tx.Save(&storefront).Error; err != nil {
        tx.Rollback()
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not update storefront",
        })
    }

    if err := tx.Commit().Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not commit transaction",
        })
    }

    return c.JSON(fiber.Map{
        "message": "Storefront updated successfully",
        "storefront": storefront,
    })
}

func DeleteStorefront(c *fiber.Ctx) error {
    user := c.Locals("user").(*models.User)
    id := c.Params("id")

    tx := database.DB.Begin()

    var storefront models.Storefront
    if err := tx.First(&storefront, id).Error; err != nil {
        tx.Rollback()
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Storefront not found",
        })
    }

    if storefront.SellerID != user.ID {
        tx.Rollback()
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Not authorized to delete this storefront",
        })
    }


    if err := tx.Delete(&storefront).Error; err != nil {
        tx.Rollback()
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not delete storefront",
        })
    }

    if err := tx.Commit().Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not commit transaction",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Storefront deleted successfully",
    })
}