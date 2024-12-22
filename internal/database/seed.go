// internal/database/seed.go
package database

import (
	"log"

	"github.com/Raviikumar001/e-com-api-go/internal/models"
)

func SeedRolesAndPermissions() error {
    // Create permissions
    permissions := []models.Permission{
        {Name: "access_web_builder", Description: "Can access web builder module"},
        {Name: "create_product", Description: "Can create products"},
        {Name: "update_product", Description: "Can update products"},
        {Name: "delete_product", Description: "Can delete products"},
        {Name: "view_product", Description: "Can view products"},
    }

    for _, permission := range permissions {
        DB.FirstOrCreate(&permission, models.Permission{Name: permission.Name})
    }

    // Create roles with permissions
    roles := []models.Role{
        {
            Name:        "admin",
            Description: "Administrator",
        },
        {
            Name:        "wholesaler",
            Description: "Wholesaler",
        },
        {
            Name:        "seller",
            Description: "Seller",
        },
        {
            Name:        "customer",
            Description: "Customer",
        },
    }

    for _, role := range roles {
        DB.FirstOrCreate(&role, models.Role{Name: role.Name})
    }

    // Assign permissions to roles
    var sellerRole models.Role
    var webBuilderPermission models.Permission

    // Find seller role
    if err := DB.Where("name = ?", "seller").First(&sellerRole).Error; err != nil {
        return err
    }

    // Find web builder permission
    if err := DB.Where("name = ?", "access_web_builder").First(&webBuilderPermission).Error; err != nil {
        return err
    }

    // Assign web builder permission to seller role
    if err := DB.Model(&sellerRole).Association("Permissions").Append(&webBuilderPermission); err != nil {
        return err
    }

    log.Println("Successfully seeded roles and permissions")
    return nil
}