// internal/models/permission.go
package models

type Permission struct {
    Base
    Name        string `json:"name" gorm:"unique;not null"`
    Description string `json:"description"`
    Roles       []Role `json:"roles" gorm:"many2many:role_permissions;"`
}

const (
    CreateProduct    = "create_product"
    UpdateProduct    = "update_product"
    DeleteProduct    = "delete_product"
    ViewProduct      = "view_product"
    ManageUsers      = "manage_users"
    AccessWebBuilder = "access_web_builder"
)