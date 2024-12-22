// internal/database/db.go
package database

import (
	"log"
	"os"

	"github.com/Raviikumar001/e-com-api-go/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// func InitDB() {
//     dsn := os.Getenv("DATABASE_URL")

//     db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//     if err != nil {
//         log.Fatal("Failed to connect to database:", err)
//     }

//     DB = db

//     // Auto Migrate schemas
//     err = db.AutoMigrate(
//         &models.User{},
//         &models.Role{},
//         &models.Permission{},
//         &models.Product{},
//         &models.Storefront{},
//     )
//     if err != nil {
//         log.Fatal("Failed to migrate database:", err)
//     }

//     // Initialize RBAC data
//     initializeRBAC(db)
// }

// internal/database/db.go
func InitDB() {
    dsn := os.Getenv("DATABASE_URL")
    
    // Add logging configuration
    dbConfig := &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info), // Add this line for detailed logging
    }

    db, err := gorm.Open(postgres.Open(dsn), dbConfig)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    DB = db

    // Auto Migrate schemas
    err = db.AutoMigrate(
        &models.User{},
        &models.Role{},
        &models.Permission{},
        &models.Product{},
        &models.Storefront{},
    )
    if err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

     // After migrations, seed the database
     if err := SeedRolesAndPermissions(); err != nil {
        log.Printf("Warning: Error seeding database: %v", err)
    }

    log.Println("Database migration completed successfully")
}

// Initialize default roles and permissions
func initializeRBAC(db *gorm.DB) {
    // Create default roles
    roles := []models.Role{
        {Name: models.AdminRole},
        {Name: models.WholesalerRole},
        {Name: models.SellerRole},
        {Name: models.CustomerRole},
    }

    for _, role := range roles {
        db.FirstOrCreate(&role, models.Role{Name: role.Name})
    }

    // Create default permissions
    permissions := []models.Permission{
        {Name: models.CreateProduct, Description: "Can create products"},
        {Name: models.UpdateProduct, Description: "Can update products"},
        {Name: models.DeleteProduct, Description: "Can delete products"},
        {Name: models.ViewProduct, Description: "Can view products"},
        {Name: models.ManageUsers, Description: "Can manage users"},
        {Name: models.AccessWebBuilder, Description: "Can access web builder"},
    }

    for _, permission := range permissions {
        db.FirstOrCreate(&permission, models.Permission{Name: permission.Name})
    }

    // Fetch all roles and permissions
    var adminRole models.Role
    var wholesalerRole models.Role
    var sellerRole models.Role
    var customerRole models.Role
    var allPermissions []models.Permission

    db.First(&adminRole, "name = ?", models.AdminRole)
    db.First(&wholesalerRole, "name = ?", models.WholesalerRole)
    db.First(&sellerRole, "name = ?", models.SellerRole)
    db.First(&customerRole, "name = ?", models.CustomerRole)
    db.Find(&allPermissions)

    // Define permission sets for each role
    var wholesalerPermissions []models.Permission
    var sellerPermissions []models.Permission
    var customerPermissions []models.Permission

    // Assign permissions based on role
    for _, permission := range allPermissions {
        // Admin gets all permissions automatically

        // Assign wholesaler permissions
        if permission.Name == models.CreateProduct ||
           permission.Name == models.UpdateProduct ||
           permission.Name == models.DeleteProduct {
            wholesalerPermissions = append(wholesalerPermissions, permission)
        }

        // Assign seller permissions
        if permission.Name == models.CreateProduct ||
           permission.Name == models.UpdateProduct {
            sellerPermissions = append(sellerPermissions, permission)
        }

        // Assign customer permissions
        if permission.Name == models.ViewProduct {
            customerPermissions = append(customerPermissions, permission)
        }
    }

    // Assign permissions to roles
    db.Model(&adminRole).Association("Permissions").Replace(allPermissions)
    db.Model(&wholesalerRole).Association("Permissions").Replace(wholesalerPermissions)
    db.Model(&sellerRole).Association("Permissions").Replace(sellerPermissions)
    db.Model(&customerRole).Association("Permissions").Replace(customerPermissions)
}