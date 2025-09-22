package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the users table in the database.
// This model is used by GORM to map to the 'users' table.
type User struct {
	gorm.Model
	ID                    uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Email                 string     `gorm:"unique;not null" json:"email"`
	Username              string     `gorm:"not null" json:"username"`
	Password              string     `gorm:"not null" json:"-"`
	PhoneNumber           string     `gorm:"not null" json:"phone_number"`
	Gender                Gender     `gorm:"type:ENUM('male', 'female', 'other');default:'other'" json:"gender"`
	Birthday              *time.Time `gorm:"default:NULL" json:"birthday"`
	Region                string     `json:"region"`
	Role                  Role       `gorm:"type:ENUM('admin', 'user');not null;default:'user'" json:"role"`
	RefreshToken          string     `gorm:"default:NULL" json:"-"`
	RefreshTokenExpiredAt *time.Time `gorm:"default:NULL" json:"-"`
	OTP                   string     `gorm:"default:NULL" json:"-"`
	OTPExpiredAt          *time.Time `gorm:"default:NULL" json:"-"`
	ResetToken            string     `gorm:"default:NULL" json:"-"`
	ResetTokenExpAt       *time.Time `gorm:"default:NULL" json:"-"`
	CreatedAt             time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	ShoppingCart      []ShoppingCart    `gorm:"foreignKey:UserID" json:"-"`
	Favorites         []Favorite        `gorm:"foreignKey:UserID" json:"-"`
	UserWardrobes     []UserWardrobe    `gorm:"foreignKey:UserID" json:"-"`
	BodyScanHistories []BodyScanHistory `gorm:"foreignKey:UserID" json:"-"`
	FaceScanHistories []FaceScanHistory `gorm:"foreignKey:UserID" json:"-"`
}

type LoginToken struct {
	AccessToken  string
	RefreshToken string
}

type ResetToken struct {
	ResetToken          string
	ResetTokenExpiredAt *time.Time
}
