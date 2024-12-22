// internal/models/permission.go
package models

type Permission struct {
    Base
    Name        string `json:"name" gorm:"uniqueIndex;not null"`
    Description string `json:"description"`
    Roles       []Role `json:"roles" gorm:"many2many:role_permissions;"`
}