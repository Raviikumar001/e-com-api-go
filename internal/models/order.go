// internal/models/order.go
package models

type Order struct {
    Base
    CustomerID  uint    `json:"customer_id" gorm:"not null"`
    SellerID    uint    `json:"seller_id" gorm:"not null"`
    TotalAmount float64 `json:"total_amount"`
    Status      string  `json:"status" gorm:"default:'pending'"`
    Items       []OrderItem `json:"items"`
}

type OrderItem struct {
    Base
    OrderID   uint    `json:"order_id" gorm:"not null"`
    ProductID uint    `json:"product_id" gorm:"not null"`
    Quantity  int     `json:"quantity" gorm:"not null"`
    Price     float64 `json:"price" gorm:"not null"`
}