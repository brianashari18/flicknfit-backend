package models

import (
	"time"

	"gorm.io/gorm"
)

// Product represents the products table.
type Product struct {
	gorm.Model
	ID          uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	BrandID     uint64  `gorm:"not null" json:"brandId"`
	Name        string  `gorm:"not null" json:"name"`
	Description string  `gorm:"not null" json:"description"`
	Discount    float64 `gorm:"not null" json:"discount"`
	Rating      float64 `gorm:"default:0.0" json:"rating"`
	Reviewer    int     `gorm:"default:0" json:"reviewer"`
	Sold        int     `gorm:"default:0" json:"sold"`

	// External links for product discovery platform
	BrandProductURL     string `gorm:"size:512" json:"brand_product_url"`     // Primary link to brand store
	WhatsAppTemplate    string `gorm:"type:text" json:"whatsapp_template"`    // Custom WhatsApp message template
	InstagramProductURL string `gorm:"size:255" json:"instagram_product_url"` // Instagram shop post
	TokopediaProductURL string `gorm:"size:255" json:"tokopedia_product_url"` // Tokopedia product page
	ShopeeProductURL    string `gorm:"size:255" json:"shopee_product_url"`    // Shopee product page

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	Brand             Brand             `gorm:"foreignKey:BrandID"`
	ProductCategories []ProductCategory `gorm:"foreignKey:ProductID"`
	ProductStyles     []ProductStyle    `gorm:"foreignKey:ProductID"`
	ProductItems      []ProductItem     `gorm:"foreignKey:ProductID"`
	Reviews           []Review          `gorm:"foreignKey:ProductID"`
}
