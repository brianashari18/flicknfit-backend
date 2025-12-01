package dtos

import (
	"flicknfit_backend/models"
	"time"
)

// UserLoginDTO is used for user authentication.
type UserLoginRequestDTO struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required"`
}

type UserRegisterRequestDTO struct {
	Email           string `json:"email" validate:"required,email"`
	Username        string `json:"username" validate:"required,min=3,max=30"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
	PhoneNumber     string `json:"phone_number" validate:"required"`
}

// ForgotPasswordDTO is used for initiating the forgot password process.
type ForgotPasswordDTO struct {
	Email string `json:"email" validate:"required,email"`
}

// VerifyOTPDTO is used for verifying the OTP.
type VerifyOTPDTO struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required"`
}

// ResetPasswordDTO is used for resetting the password after OTP verification.
type UserResetPasswordRequestDTO struct {
	Email           string `json:"email" validate:"required,email"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
	ResetToken      string `json:"reset_token" validate:"required"`
}

// ResponseDTO defines a standardized JSON response format.
type ResponseDTO struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// UserLoginResponseDTO is the response DTO for a successful login.
type UserLoginResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func ToUserLoginResponseDTO(t *models.LoginToken) UserLoginResponseDTO {
	return UserLoginResponseDTO{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
	}
}

type UserResetTokenResponseDTO struct {
	ResetToken          string     `json:"reset_token"`
	ResetTokenExpiredAt *time.Time `json:"expired_at"`
}

func ToUserResetTokenResponseDTO(t *models.ResetToken) UserResetTokenResponseDTO {
	return UserResetTokenResponseDTO{
		ResetToken:          t.ResetToken,
		ResetTokenExpiredAt: t.ResetTokenExpiredAt,
	}
}

// UserCreateDTO is used for creating a new user (admin dashboard).
type UserAdminCreateRequestDTO struct {
	Email       string `json:"email" validate:"required,email"`
	Username    string `json:"username" validate:"required,min=3"`
	Password    string `json:"password" validate:"required,min=8"`
	PhoneNumber string `json:"phone_number"`
	Gender      string `json:"gender"`
	Birthday    string `json:"birthday"`
	Region      string `json:"region"`
	Role        string `json:"role" validate:"required,oneof=admin user"`
}

// UserUpdateDTO is used for updating an existing user's details (admin or user).
type UserAdminUpdateRequestDTO struct {
	Email       string `json:"email,omitempty" validate:"omitempty,email"`
	Username    string `json:"username,omitempty" validate:"omitempty,min=3"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Birthday    string `json:"birthday,omitempty"`
	Region      string `json:"region,omitempty"`
}

type UserResponseDTO struct {
	ID          uint64     `json:"id"`
	Email       string     `json:"email"`
	Username    string     `json:"username"`
	PhoneNumber string     `json:"phone_number"`
	Gender      string     `json:"gender"`
	Birthday    *time.Time `json:"birthday,omitempty"`
	Region      string     `json:"region"`
	Role        string     `json:"role"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func ToUserResponseDTO(u *models.User) UserResponseDTO {
	return UserResponseDTO{
		ID:          u.ID,
		Email:       u.Email,
		Username:    u.Username,
		PhoneNumber: u.PhoneNumber,
		Gender:      string(u.Gender),
		Birthday:    u.Birthday,
		Region:      u.Region,
		Role:        string(u.Role),
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

type UserAdminResponseDTO struct {
	ID                    uint64     `json:"id"`
	Email                 string     `json:"email"`
	Username              string     `json:"username"`
	Password              string     `json:"password"`
	PhoneNumber           string     `json:"phone_number"`
	Gender                string     `json:"gender"`
	Birthday              *time.Time `json:"birthday,omitempty"`
	Region                string     `json:"region"`
	Role                  string     `json:"role"`
	RefreshToken          string     `json:"refresh_token"`
	RefreshTokenExpiredAt *time.Time `json:"refresh_token_expired_at"`
	OTP                   string     `json:"otp"`
	OTPExpiredAt          *time.Time `json:"otp_expired_at"`
	ResetToken            string     `json:"reset_token"`
	ResetTokenExpAt       *time.Time `json:"reset_token_exp_at"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

func ToUserAdminResponseDTO(u *models.User) UserAdminResponseDTO {
	return UserAdminResponseDTO{
		ID:                    u.ID,
		Email:                 u.Email,
		Username:              u.Username,
		Password:              u.Password,
		PhoneNumber:           u.PhoneNumber,
		Gender:                string(u.Gender),
		Birthday:              u.Birthday,
		Region:                u.Region,
		Role:                  string(u.Role),
		RefreshToken:          u.RefreshToken,
		RefreshTokenExpiredAt: u.RefreshTokenExpiredAt,
		OTP:                   u.OTP,
		OTPExpiredAt:          u.OTPExpiredAt,
		ResetToken:            u.ResetToken,
		ResetTokenExpAt:       u.ResetTokenExpAt,
		CreatedAt:             u.CreatedAt,
		UpdatedAt:             u.UpdatedAt,
	}
}

func ToUserAdminResponseDTOs(users []*models.User) []UserAdminResponseDTO {
	result := make([]UserAdminResponseDTO, 0, len(users))
	for _, u := range users {
		result = append(result, ToUserAdminResponseDTO(u))
	}
	return result
}

type UserEditProfileRequestDTO struct {
	Email       string     `json:"email" validate:"omitempty,email"`
	Username    string     `json:"username" validate:"omitempty,min=3,max=30"`
	PhoneNumber string     `json:"phone_number" validate:"omitempty"`
	Gender      string     `json:"gender" validate:"omitempty,oneof=male female other"`
	Birthday    *time.Time `json:"birthday,omitempty" validate:"omitempty"`
	Region      string     `json:"region" validate:"omitempty,min=2,max=100"`
}

type UserEditProfileResponseDTO struct {
	ID          uint64     `json:"id"`
	Email       string     `json:"email"`
	Username    string     `json:"username"`
	PhoneNumber string     `json:"phone_number"`
	Gender      string     `json:"gender"`
	Birthday    *time.Time `json:"birthday,omitempty"`
	Region      string     `json:"region"`
	Role        string     `json:"role"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func ToUserEditProfileResponseDTO(u models.User) UserEditProfileResponseDTO {
	return UserEditProfileResponseDTO{
		ID:          u.ID,
		Email:       u.Email,
		Username:    u.Username,
		PhoneNumber: u.PhoneNumber,
		Gender:      string(u.Gender),
		Birthday:    u.Birthday,
		Region:      u.Region,
		Role:        string(u.Role),
		UpdatedAt:   u.UpdatedAt,
	}
}
