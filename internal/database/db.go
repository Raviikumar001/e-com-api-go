// internal/database/db.go
package database

import (
	"log"
	"os"
	"time"

	"github.com/Raviikumar001/e-com-api-go/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("DATABASE_URL environment variable is not set")
    }
    
    // Add logging configuration
    dbConfig := &gorm.Config{
        Logger: logger.New(
            log.New(os.Stdout, "\r\n", log.LstdFlags),
            logger.Config{
                SlowThreshold:             time.Second,
                LogLevel:                  logger.Info,
                IgnoreRecordNotFoundError: true,
                Colorful:                  true,
            },
        ),
    }

    // Connect to database
    db, err := gorm.Open(postgres.Open(dsn), dbConfig)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Configure connection pool
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatal("Failed to get database instance:", err)
    }

    // Set reasonable defaults for connection pool
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    DB = db

    // Auto Migrate schemas
    log.Println("Starting database migration...")
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
    log.Println("Database migration completed successfully")

    // Seed initial data
    log.Println("Starting database seeding...")
    if err := SeedRolesAndPermissions(); err != nil {
        log.Printf("Warning: Error seeding database: %v", err)
    }
    log.Println("Database seeding completed")
}

