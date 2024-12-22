// internal/models/product.go
package models

type Product struct {
    Base
    Name          string  `json:"name" gorm:"not null"`
    Description   string  `json:"description"`
    Price         float64 `json:"price" gorm:"not null"`
    CostPrice     float64 `json:"cost_price"`
    Stock         int     `json:"stock" gorm:"not null"`
    Category      string  `json:"category"`
    ImageURL      string  `json:"image_url"`
    IsPublished   bool    `json:"is_published" gorm:"default:false"`
    SellerID      *uint   `json:"seller_id"`       
    WholesalerID  *uint   `json:"wholesaler_id"`  
    Seller        *User    `json:"seller" gorm:"foreignKey:SellerID"`
    Wholesaler    *User    `json:"wholesaler" gorm:"foreignKey:WholesalerID"`
}