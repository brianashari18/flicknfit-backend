package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// StringArray is a custom type for storing JSON string arrays in database
type StringArray []string

// Scan implements sql.Scanner interface
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan StringArray")
	}
	return json.Unmarshal(bytes, s)
}

// Value implements driver.Valuer interface
func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

// FaceScanHistory represents the face_scan_histories table.
type FaceScanHistory struct {
	gorm.Model
	ID                   uint64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID               uint64      `gorm:"not null" json:"user_id"`
	ScanName             string      `gorm:"not null;size:100" json:"scan_name"`
	SkinTone             string      `gorm:"not null" json:"skin_tone"`
	ImagePath            string      `gorm:"not null;size:512" json:"image_path"`
	ColorRecommendations StringArray `gorm:"type:json" json:"color_recommendations"`
	Confidence           *float64    `gorm:"type:decimal(5,4)" json:"confidence,omitempty"`
	CreatedAt            time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time   `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID"`
}
