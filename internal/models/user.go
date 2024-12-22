// internal/models/user.go
package models

import (
	"github.com/Raviikumar001/e-com-api-go/internal/utils"
	"gorm.io/gorm"
)

type User struct {
    Base
    Email     string `json:"email" gorm:"unique;not null"`
    Password  string `json:"-" gorm:"not null"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    RoleID    uint   `json:"role_id"`
    Role      Role   `json:"role"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    hashedPassword, err := utils.HashPassword(u.Password)
    if err != nil {
        return err
    }
    u.Password = hashedPassword
    return nil
}