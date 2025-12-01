package models

// Role represents the role enum for a user.
// It can be 'admin' or 'user'.
type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

// Gender represents the gender enum for a user.
// It can be 'male', 'female', or 'other'.
type Gender string

const (
	MaleGender   Gender = "male"
	FemaleGender Gender = "female"
	OtherGender  Gender = "other"
)

// AuthProvider represents the authentication provider enum for a user.
// It can be 'local' (email/password), 'google' (Google OAuth), or 'facebook' (Facebook OAuth).
type AuthProvider string

const (
	LocalAuthProvider    AuthProvider = "local"
	GoogleAuthProvider   AuthProvider = "google"
	FacebookAuthProvider AuthProvider = "facebook"
)
