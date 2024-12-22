// internal/models/role.go
package models

type Role struct {
    Base
    Name        string       `json:"name" gorm:"unique;not null"`
    Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
}

const (
    AdminRole      = "admin"
    WholesalerRole = "wholesaler"
    SellerRole     = "seller"
    CustomerRole   = "customer"
)
