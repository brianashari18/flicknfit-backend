package models

import (
	"time"

	"gorm.io/gorm"
)

// ColorToneRecommendation represents the color_tone_recommendations table.
type ColorToneRecommendation struct {
	gorm.Model
	ID                uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	FaceScanHistoryID uint64    `gorm:"not null" json:"face_scan_history_id"`
	ColorTone         string    `gorm:"not null" json:"color_tone"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	FaceScanHistory FaceScanHistory `gorm:"foreignKey:FaceScanHistoryID"`
}
