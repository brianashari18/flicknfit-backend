package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the users table in the database.
// This model is used by GORM to map to the 'users' table.
type User struct {
	gorm.Model
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Email       string     `gorm:"unique;not null" json:"email"`
	Username    string     `gorm:"unique;not null" json:"username"`
	Password    string     `gorm:"default:NULL" json:"-"`            // Nullable for OAuth users
	PhoneNumber string     `gorm:"default:NULL" json:"phone_number"` // Optional
	Gender      Gender     `gorm:"type:ENUM('male', 'female', 'other');default:'other'" json:"gender"`
	Birthday    *time.Time `gorm:"default:NULL" json:"birthday"`
	Region      string     `json:"region"`
	Role        Role       `gorm:"type:ENUM('admin', 'user');not null;default:'user'" json:"role"`

	// OAuth Authentication Fields
	AuthProvider          AuthProvider `gorm:"type:ENUM('local', 'google', 'facebook');default:'local'" json:"auth_provider"`
	AuthProviderID        string       `gorm:"default:NULL" json:"-"` // Google/Facebook User ID
	ProfilePictureURL     string       `gorm:"default:NULL" json:"profile_picture_url"`
	IsEmailVerified       bool         `gorm:"default:false" json:"is_email_verified"`
	LastLogin             *time.Time   `gorm:"default:NULL" json:"last_login"`
	RefreshToken          string       `gorm:"default:NULL" json:"-"`
	RefreshTokenExpiredAt *time.Time   `gorm:"default:NULL" json:"-"`
	OTP                   string       `gorm:"default:NULL" json:"-"`
	OTPExpiredAt          *time.Time   `gorm:"default:NULL" json:"-"`
	ResetToken            string       `gorm:"default:NULL" json:"-"`
	ResetTokenExpAt       *time.Time   `gorm:"default:NULL" json:"-"`
	CreatedAt             time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time    `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	SavedItems        []SavedItems      `gorm:"foreignKey:UserID" json:"-"`
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

// IsOAuthUser checks if user signed up via OAuth (Google/Facebook)
func (u *User) IsOAuthUser() bool {
	return u.AuthProvider == GoogleAuthProvider || u.AuthProvider == FacebookAuthProvider
}

// IsLocalUser checks if user signed up with email/password
func (u *User) IsLocalUser() bool {
	return u.AuthProvider == LocalAuthProvider
}

// NeedsPassword returns true if user should have a password
func (u *User) NeedsPassword() bool {
	return u.IsLocalUser() && u.Password == ""
}
