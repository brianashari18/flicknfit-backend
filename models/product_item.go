package models

import "gorm.io/gorm"

type ProductItem struct {
	gorm.Model
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uint64 `gorm:"not null" json:"product_id"`
	SKU       string `gorm:"size:20;not null;unique" json:"sku"`
	Price     int    `gorm:"not null" json:"price"`
	Stock     int    `gorm:"not null" json:"stock"`
	Sold      int    `gorm:"default:0" json:"sold"`
	PhotoURL  string `gorm:"size:255" json:"photo_url"`

	// Relationships
	Product           Product                `gorm:"foreignKey:ProductID"`
	ShoppingCartItems []ShoppingCartItem     `gorm:"foreignKey:ProductItemID"`
	Favorites         []Favorite             `gorm:"foreignKey:ProductItemID"`
	Configurations    []ProductConfiguration `gorm:"foreignKey:ProductItemID"`
}
