package models

import (
	"gorm.io/gorm"
)

// ProductCategory represents the product_categories join table.
type ProductCategory struct {
	gorm.Model
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uint64 `gorm:"not null" json:"product_id"`
	Category  string `gorm:"not null" json:"category"`

	// Relationships
	Product Product `gorm:"foreignKey:ProductID"`
}
