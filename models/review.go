package models

import (
	"time"

	"gorm.io/gorm"
)

// Review represents the reviews table.
type Review struct {
	gorm.Model
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID  uint64    `gorm:"not null" json:"product_id"`
	Rating     int       `gorm:"default:0" json:"rating"`
	ReviewText string    `gorm:"not null;column:review" json:"review_text"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Product Product `gorm:"foreignKey:ProductID"`
}
