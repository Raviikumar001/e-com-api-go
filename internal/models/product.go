// internal/models/product.go
package models

type Product struct {
    Base
    Name         string  `json:"name" gorm:"not null"`
    Description  string  `json:"description"`
    Price        float64 `json:"price" gorm:"not null"`
    Stock        int     `json:"stock" gorm:"not null"`
    WholesalerID uint    `json:"wholesaler_id" gorm:"not null"`
    SellerID     uint    `json:"seller_id"`
    Status       string  `json:"status" gorm:"default:'active'"` // active, inactive, deleted
    Images       []string `json:"images" gorm:"type:text[]"`
    Category     string   `json:"category"`
}

