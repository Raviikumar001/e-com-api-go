// internal/models/role.go
package models

type Role struct {
    Base
    Name        string       `json:"name" gorm:"uniqueIndex;not null"`
    Description string       `json:"description"`
    Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
}