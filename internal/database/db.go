// internal/database/db.go

package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB() (*DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Connected to database")

	return &DB{db}, nil
}

func (db *DB) Migrate() error {
	return db.AutoMigrate(
		&User{},
		&Role{},
		&Permission{},
		&UserRole{},
		&RolePermission{},
	)
}

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
}

type Role struct {
	gorm.Model
	Name        string
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	gorm.Model
	Name  string
	Roles []Role `gorm:"many2many:role_permissions;"`
}

type UserRole struct {
	gorm.Model
	UserID uint
	RoleID uint
}

type RolePermission struct {
	gorm.Model
	RoleID       uint
	PermissionID uint
}
