package models

import (
	"time"

	"gorm.io/gorm"
)

// ProductClick represents a user's click on a product
type ProductClick struct {
	gorm.Model
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    *uint64   `gorm:"index" json:"user_id"` // Nullable for anonymous users
	ProductID uint64    `gorm:"not null;index" json:"product_id"`
	BrandID   uint64    `gorm:"not null;index" json:"brand_id"`
	ClickedAt time.Time `gorm:"autoCreateTime;index" json:"clicked_at"`
	IPAddress string    `gorm:"size:45" json:"ip_address"` // IPv6 max length
	UserAgent string    `gorm:"type:text" json:"user_agent"`

	// Relationships
	User    *User   `gorm:"foreignKey:UserID"`
	Product Product `gorm:"foreignKey:ProductID"`
	Brand   Brand   `gorm:"foreignKey:BrandID"`
}

// TableName specifies the table name for ProductClick model
func (ProductClick) TableName() string {
	return "product_clicks"
}
