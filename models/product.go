package models

import (
	"time"

	"gorm.io/gorm"
)

// Product represents the products table.
type Product struct {
	gorm.Model
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	BrandID     uint64    `gorm:"not null" json:"brandId"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"not null" json:"description"`
	Discount    float64   `gorm:"not null" json:"discount"`
	Rating      float64   `gorm:"default:0.0" json:"rating"`
	Reviewer    int       `gorm:"default:0" json:"reviewer"`
	Sold        int       `gorm:"default:0" json:"sold"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	Brand             Brand             `gorm:"foreignKey:BrandID"`
	ProductCategories []ProductCategory `gorm:"foreignKey:ProductID"`
	ProductStyles     []ProductStyle    `gorm:"foreignKey:ProductID"`
	ProductItems      []ProductItem     `gorm:"foreignKey:ProductID"`
	Reviews           []Review          `gorm:"foreignKey:ProductID"`
}
