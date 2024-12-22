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

    // Create all permissions
    for _, permission := range permissions {
        if err := DB.Where(models.Permission{Name: permission.Name}).
            FirstOrCreate(&permission).Error; err != nil {
            return err
        }
    }

    // Create roles
    roles := []models.Role{
        {Name: "admin", Description: "Administrator"},
        {Name: "wholesaler", Description: "Wholesaler"},
        {Name: "seller", Description: "Seller"},
        {Name: "customer", Description: "Customer"},
    }

    // Create all roles
    for _, role := range roles {
        if err := DB.Where(models.Role{Name: role.Name}).
            FirstOrCreate(&role).Error; err != nil {
            return err
        }
    }

    // Fetch all created roles and permissions
    var adminRole, wholesalerRole, sellerRole, customerRole models.Role
    var allPermissions []models.Permission

    if err := DB.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
        return err
    }
    if err := DB.Where("name = ?", "wholesaler").First(&wholesalerRole).Error; err != nil {
        return err
    }
    if err := DB.Where("name = ?", "seller").First(&sellerRole).Error; err != nil {
        return err
    }
    if err := DB.Where("name = ?", "customer").First(&customerRole).Error; err != nil {
        return err
    }
    if err := DB.Find(&allPermissions).Error; err != nil {
        return err
    }

    // Get specific permissions
    var createProduct, updateProduct, deleteProduct, viewProduct, accessWebBuilder models.Permission
    if err := DB.Where("name = ?", "create_product").First(&createProduct).Error; err != nil {
        return err
    }
    if err := DB.Where("name = ?", "update_product").First(&updateProduct).Error; err != nil {
        return err
    }
    if err := DB.Where("name = ?", "delete_product").First(&deleteProduct).Error; err != nil {
        return err
    }
    if err := DB.Where("name = ?", "view_product").First(&viewProduct).Error; err != nil {
        return err
    }
    if err := DB.Where("name = ?", "access_web_builder").First(&accessWebBuilder).Error; err != nil {
        return err
    }

    // Assign permissions to roles
    
    // Admin gets all permissions
    if err := DB.Model(&adminRole).Association("Permissions").Replace(&allPermissions); err != nil {
        return err
    }

    // Wholesaler permissions
    wholesalerPermissions := []models.Permission{
        createProduct,
        updateProduct,
        deleteProduct,
        viewProduct,
    }
    if err := DB.Model(&wholesalerRole).Association("Permissions").Replace(&wholesalerPermissions); err != nil {
        return err
    }

    // Seller permissions
    sellerPermissions := []models.Permission{
        createProduct,
        updateProduct,
        deleteProduct,
        viewProduct,
        accessWebBuilder,
    }
    if err := DB.Model(&sellerRole).Association("Permissions").Replace(&sellerPermissions); err != nil {
        return err
    }

    // Customer permissions
    customerPermissions := []models.Permission{
        viewProduct,
    }
    if err := DB.Model(&customerRole).Association("Permissions").Replace(&customerPermissions); err != nil {
        return err
    }

    log.Println("Successfully seeded roles and permissions")
    return nil
}


// // internal/database/seed.go
// package database

// import (
// 	"log"

// 	"github.com/Raviikumar001/e-com-api-go/internal/models"
// )

// func SeedRolesAndPermissions() error {
//     // Create permissions
//     permissions := []models.Permission{
//         {Name: "access_web_builder", Description: "Can access web builder module"},
//         {Name: "create_product", Description: "Can create products"},
//         {Name: "update_product", Description: "Can update products"},
//         {Name: "delete_product", Description: "Can delete products"},
//         {Name: "view_product", Description: "Can view products"},
//     }

//     for _, permission := range permissions {
//         DB.FirstOrCreate(&permission, models.Permission{Name: permission.Name})
//     }

//     // Create roles with permissions
//     roles := []models.Role{
//         {
//             Name:        "admin",
//             Description: "Administrator",
//         },
//         {
//             Name:        "wholesaler",
//             Description: "Wholesaler",
//         },
//         {
//             Name:        "seller",
//             Description: "Seller",
//         },
//         {
//             Name:        "customer",
//             Description: "Customer",
//         },
//     }

//     for _, role := range roles {
//         DB.FirstOrCreate(&role, models.Role{Name: role.Name})
//     }

//     // Assign permissions to roles
//     var sellerRole models.Role
//     var webBuilderPermission models.Permission

//     // Find seller role
//     if err := DB.Where("name = ?", "seller").First(&sellerRole).Error; err != nil {
//         return err
//     }

//     // Find web builder permission
//     if err := DB.Where("name = ?", "access_web_builder").First(&webBuilderPermission).Error; err != nil {
//         return err
//     }

//     // Assign web builder permission to seller role
//     if err := DB.Model(&sellerRole).Association("Permissions").Append(&webBuilderPermission); err != nil {
//         return err
//     }

//     log.Println("Successfully seeded roles and permissions")
//     return nil
// }