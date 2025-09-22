package models

import (
	"time"

	"gorm.io/gorm"
)

// Brand represents the brands table in the database.
// This model is used by GORM to map to the 'brands' table.
type Brand struct {
	gorm.Model
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"not null" json:"name"`
	Description   string    `gorm:"not null" json:"description"`
	Rating        float64   `gorm:"default:0.0" json:"rating"`
	Reviewer      uint      `gorm:"default:0" json:"reviewer"`
	LogoURL       string    `json:"logo_url"`
	WebsiteURL    string    `json:"website_url"`
	TotalProducts uint      `gorm:"default:0" json:"total_products"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Products []Product `gorm:"foreignKey:BrandID"`
}
