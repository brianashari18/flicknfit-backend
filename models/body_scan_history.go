package models

import (
	"time"

	"gorm.io/gorm"
)

// BodyScanHistory represents the body_scan_histories table.
type BodyScanHistory struct {
	gorm.Model
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"not null" json:"user_id"`
	BodyShape string    `gorm:"not null" json:"body_shape"`
	ImageURL  string    `gorm:"not null;size:512" json:"image_url"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User                 User                  `gorm:"foreignKey:UserID"`
	StyleRecommendations []StyleRecommendation `gorm:"foreignKey:BodyScanHistoryID"`
}
