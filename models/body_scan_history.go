package models

import (
	"time"

	"gorm.io/gorm"
)

// BodyScanHistory represents the body_scan_histories table.
type BodyScanHistory struct {
	gorm.Model
	ID                   uint64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID               uint64      `gorm:"not null" json:"user_id"`
	ScanName             string      `gorm:"not null;size:100" json:"scan_name"`
	BodyType             string      `gorm:"not null" json:"body_type"`
	Gender               string      `gorm:"not null;size:10" json:"gender"` // "woman" or "man"
	ImagePath            string      `gorm:"not null;size:512" json:"image_path"`
	StyleRecommendations StringArray `gorm:"type:json" json:"style_recommendations"`
	Confidence           float64     `gorm:"not null;type:decimal(5,4)" json:"confidence"`
	CreatedAt            time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time   `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID"`
}
