package models

import (
	"time"

	"gorm.io/gorm"
)

// UserWardrobe represents the user_wardrobes table.
type UserWardrobe struct {
	gorm.Model
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"not null" json:"user_id"`
	Category  string    `gorm:"not null" json:"category"`
	ImageURL  string    `gorm:"not null;size:512" json:"image_url"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID"`
}
