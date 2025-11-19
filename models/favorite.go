package models

import (
	"time"

	"gorm.io/gorm"
)

// Favorite represents the favorites table.
type Favorite struct {
	gorm.Model
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64    `gorm:"not null" json:"user_id"`
	ProductItemID uint64    `gorm:"not null" json:"product_item_id"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User        User        `gorm:"foreignKey:UserID"`
	ProductItem ProductItem `gorm:"foreignKey:ProductItemID"`
}
