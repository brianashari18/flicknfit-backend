package models

import (
	"time"

	"gorm.io/gorm"
)

// FaceScanHistory represents the face_scan_histories table.
type FaceScanHistory struct {
	gorm.Model
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"not null" json:"user_id"`
	SkinTone  string    `gorm:"not null" json:"skin_tone"`
	ImageURL  string    `gorm:"not null;size:512" json:"image_url"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User                     User                      `gorm:"foreignKey:UserID"`
	ColorToneRecommendations []ColorToneRecommendation `gorm:"foreignKey:FaceScanHistoryID"`
}
