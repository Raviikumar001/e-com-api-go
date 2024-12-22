// internal/models/storefront.go
package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Settings represents the JSON structure for storefront settings
type Settings map[string]interface{}

// Value implements the driver.Valuer interface
func (s Settings) Value() (driver.Value, error) {
    return json.Marshal(s)
}

// Scan implements the sql.Scanner interface
func (s *Settings) Scan(value interface{}) error {
    if value == nil {
        *s = nil
        return nil
    }
    bytes, ok := value.([]byte)
    if !ok {
        return errors.New("type assertion to []byte failed")
    }
    return json.Unmarshal(bytes, &s)
}

type Storefront struct {
    Base
    Name        string   `json:"name" gorm:"not null"`
    Description string   `json:"description"`
    SellerID    uint     `json:"seller_id" gorm:"not null"`
    Theme       string   `json:"theme" gorm:"default:'default'"`
    Domain      string   `json:"domain" gorm:"unique"`
    Settings    Settings `json:"settings" gorm:"type:jsonb"`
}