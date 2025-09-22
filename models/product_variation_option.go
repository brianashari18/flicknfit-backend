package models

import "gorm.io/gorm"

type ProductVariationOption struct {
	gorm.Model
	ID                 uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductAttributeID uint64 `gorm:"not null" json:"product_attribute_id"`
	Value              string `gorm:"size:50;not null" json:"value"`

	// Relationships
	ProductVariation ProductVariation `gorm:"foreignKey:ProductAttributeID"`
}
