package models

import "gorm.io/gorm"

type ProductConfiguration struct {
	gorm.Model
	ProductItemID           uint64 `gorm:"primaryKey;not null" json:"product_item_id"`
	ProductAttributeValueID uint64 `gorm:"primaryKey;not null" json:"product_attribute_value_id"`

	// Relationships
	ProductItem            ProductItem            `gorm:"foreignKey:ProductItemID"`
	ProductVariationOption ProductVariationOption `gorm:"foreignKey:ProductAttributeValueID"`
}
