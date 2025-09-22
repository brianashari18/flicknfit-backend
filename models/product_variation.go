package models

import (
	"gorm.io/gorm"
)

// ProductVariant represents the product_variants table.
type ProductVariation struct {
	gorm.Model
	ID   uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"size:50;not null" json:"name"`
	// Hubungan ke ProductVariationOption
	Options []ProductVariationOption `gorm:"foreignKey:ProductAttributeID"`
}
