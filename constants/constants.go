package constants

import "time"

// API Constants
const (
	APIVersion = "v1"
	APIPrefix  = "/api/" + APIVersion
)

// Authentication Constants
const (
	JWTAccessTokenExpiry  = time.Hour * 24     // 24 hours
	JWTRefreshTokenExpiry = time.Hour * 24 * 7 // 7 days
	OTPExpiry             = time.Minute * 15   // 15 minutes
	OTPLength             = 6
)

// User Roles
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// Database Constants
const (
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// File Upload Constants
const (
	MaxFileSize      = 10 << 20 // 10MB
	AllowedImageExts = "jpg,jpeg,png,webp"
)

// Cache Constants
const (
	CacheKeyPrefix     = "flicknfit:"
	CacheDefaultExpiry = time.Hour * 1
)

// Validation Constants
const (
	MinPasswordLength = 8
	MaxPasswordLength = 72
	MinUsernameLength = 3
	MaxUsernameLength = 50
)

// Email Templates
const (
	EmailTemplateOTP           = "otp"
	EmailTemplateWelcome       = "welcome"
	EmailTemplatePasswordReset = "password_reset"
)

// Product Constants
const (
	ProductStatusActive   = "active"
	ProductStatusInactive = "inactive"
	ProductStatusDraft    = "draft"
)

// Shopping Cart Constants
const (
	MaxCartItems      = 100
	CartItemMaxQty    = 999
	CartSessionExpiry = time.Hour * 24 * 30 // 30 days
)

// HTTP Headers
const (
	HeaderContentType   = "Content-Type"
	HeaderAuthorization = "Authorization"
	HeaderUserAgent     = "User-Agent"
	HeaderXRealIP       = "X-Real-IP"
	HeaderXForwardedFor = "X-Forwarded-For"
)

// Response Messages
const (
	MsgSuccess            = "Operation completed successfully"
	MsgCreated            = "Resource created successfully"
	MsgUpdated            = "Resource updated successfully"
	MsgDeleted            = "Resource deleted successfully"
	MsgNotFound           = "Resource not found"
	MsgUnauthorized       = "Authentication required"
	MsgForbidden          = "Access denied"
	MsgInvalidInput       = "Invalid input provided"
	MsgInternalError      = "Internal server error"
	MsgEmailAlreadyExists = "Email already exists"
	MsgInvalidCredentials = "Invalid email or password"
	MsgTokenExpired       = "Token has expired"
	MsgTokenInvalid       = "Invalid token"
)
