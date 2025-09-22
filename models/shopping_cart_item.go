package models

import (
	"time"

	"gorm.io/gorm"
)

type ShoppingCartItem struct {
	gorm.Model
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ShoppingCartID uint64    `gorm:"not null" json:"shopping_cart_id"`
	ProductItemID  uint64    `gorm:"not null" json:"product_item_id"`
	Quantity       int       `gorm:"not null;default:1" json:"quantity"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	ShoppingCart ShoppingCart `gorm:"foreignKey:ShoppingCartID"`
	ProductItem  ProductItem  `gorm:"foreignKey:ProductItemID"`
}
