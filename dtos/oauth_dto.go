package dtos

import "time"

// FirebaseAuthRequestDTO represents Firebase authentication request
type FirebaseAuthRequestDTO struct {
	FirebaseToken string `json:"firebase_token" validate:"required"` // Firebase ID token
}

// OAuthLoginRequestDTO represents OAuth login via Firebase
type OAuthLoginRequestDTO struct {
	FirebaseToken     string     `json:"firebase_token" validate:"required"`
	AuthProvider      string     `json:"auth_provider" validate:"required,oneof=google facebook"` // google or facebook
	Email             string     `json:"email" validate:"required,email"`
	Username          string     `json:"username" validate:"required"`
	ProfilePictureURL string     `json:"profile_picture_url"`
	AuthProviderID    string     `json:"auth_provider_id" validate:"required"` // Google/Facebook user ID
	PhoneNumber       string     `json:"phone_number"`
	Gender            string     `json:"gender" validate:"omitempty,oneof=male female other"`
	Birthday          *time.Time `json:"birthday"`
}

// OAuthRegisterRequestDTO represents OAuth registration
type OAuthRegisterRequestDTO struct {
	FirebaseToken     string     `json:"firebase_token" validate:"required"`
	AuthProvider      string     `json:"auth_provider" validate:"required,oneof=google facebook"`
	Email             string     `json:"email" validate:"required,email"`
	Username          string     `json:"username" validate:"required,min=3,max=50"`
	ProfilePictureURL string     `json:"profile_picture_url"`
	AuthProviderID    string     `json:"auth_provider_id" validate:"required"`
	PhoneNumber       string     `json:"phone_number"`
	Gender            string     `json:"gender" validate:"omitempty,oneof=male female other"`
	Birthday          *time.Time `json:"birthday"`
	Region            string     `json:"region"`
}

// FirebaseUserData represents parsed Firebase token data
type FirebaseUserData struct {
	UID           string `json:"uid"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Provider      string `json:"provider"` // google.com or facebook.com
}
