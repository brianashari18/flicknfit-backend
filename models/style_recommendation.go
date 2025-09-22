package models

import (
	"time"

	"gorm.io/gorm"
)

// StyleRecommendation represents the style_recommendations table.
type StyleRecommendation struct {
	gorm.Model
	ID                uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	BodyScanHistoryID uint64    `gorm:"not null" json:"body_scan_history_id"`
	Style             string    `gorm:"not null" json:"style"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	BodyScanHistory BodyScanHistory `gorm:"foreignKey:BodyScanHistoryID"`
}
