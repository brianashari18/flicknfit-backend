package models

import (
	"gorm.io/gorm"
)

// ProductStyle represents the product_styles join table.
type ProductStyle struct {
	gorm.Model
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uint64 `gorm:"not null" json:"product_id"`
	Style     string `gorm:"not null" json:"style"`

	// Relationships
	Product Product `gorm:"foreignKey:ProductID"`
}
