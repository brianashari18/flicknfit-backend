package models

import (
	"time"

	"gorm.io/gorm"
)

type SavedItemsList struct {
	gorm.Model
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SavedItemsID  uint64    `gorm:"not null" json:"saved_items_id"`
	ProductItemID uint64    `gorm:"not null" json:"product_item_id"`
	Quantity      int       `gorm:"not null;default:1" json:"quantity"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	SavedItems  SavedItems  `gorm:"foreignKey:SavedItemsID;references:ID"`
	ProductItem ProductItem `gorm:"foreignKey:ProductItemID"`
}

func (SavedItemsList) TableName() string {
	return "saved_items_items"
}
