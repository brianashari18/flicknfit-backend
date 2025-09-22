package services

import (
	"errors"
	"flicknfit_backend/config"
	"flicknfit_backend/dtos"
	"flicknfit_backend/models"
	"flicknfit_backend/repositories"
	"flicknfit_backend/utils"
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// CustomClaims defines the claims for the JWT.
type CustomClaims struct {
	UserID uint64 `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// UserService defines the interface for business logic related to users.
type UserService interface {
	AdminCreateUser(dto *dtos.UserAdminCreateRequestDTO) error
	GetAllUsers() ([]*models.User, error)
	GetUserByID(id uint64) (*models.User, error)
	AdminUpdateUser(id uint64, dto *dtos.UserAdminUpdateRequestDTO) error
	AdminDeleteUser(id uint64) error
	RegisterUser(dto *dtos.UserRegisterRequestDTO) error
	LoginUser(dto *dtos.UserLoginRequestDTO) (*models.LoginToken, error)
	ForgotPassword(dto *dtos.ForgotPasswordDTO) error
	VerifyOTP(dto *dtos.VerifyOTPDTO) (*models.ResetToken, error)
	ResetPassword(dto *dtos.UserResetPasswordRequestDTO) error
	LogoutUser(userID uint64) error
	RefreshToken(token string) (*models.LoginToken, error)
	EditProfile(id uint64, dto *dtos.UserEditProfileRequestDTO) (*models.User, error)
}

// userService is the implementation of UserService.
type userService struct {
	userRepository repositories.UserRepository
	cfg            *config.Config
}

// NewUserService creates a new instance of UserService.
// It now accepts the jwtSecret from the application's configuration.
func NewUserService(userRepository repositories.UserRepository, cfg *config.Config) UserService {
	return &userService{
		userRepository: userRepository,
		cfg:            cfg,
	}
}

// generateAccessToken creates a new JWT access token.
func (s *userService) generateAccessToken(user *models.User) (string, error) {
	logger := utils.GetLogger()
	claims := &CustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), // 1 hour expiration
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := tok.SignedString([]byte(s.cfg.JwtSecretKey))
	if err != nil {
		logger.Error("failed signing JWT: ", err)
		return "", err
	}

	return signed, nil

}

// generateRefreshToken creates a new JWT refresh token.
func (s *userService) generateRefreshToken(user *models.User) (string, error) {
	logger := utils.GetLogger()
	claims := &CustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hour expiration
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := tok.SignedString([]byte(s.cfg.JwtSecretKey))
	if err != nil {
		logger.Error("failed signing JWT: ", err)
		return "", err
	}

	return signed, nil
}

// AdminCreateUser handles creating a new user by an admin.
func (s *userService) AdminCreateUser(dto *dtos.UserAdminCreateRequestDTO) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:       dto.Email,
		Username:    dto.Username,
		Password:    string(hashedPassword),
		PhoneNumber: dto.PhoneNumber,
		Gender:      models.Gender(dto.Gender),
		Region:      dto.Region,
		Role:        models.Role(dto.Role),
	}

	return s.userRepository.CreateUser(user)
}

// GetAllUsers retrieves all users.
func (s *userService) GetAllUsers() ([]*models.User, error) {
	return s.userRepository.GetAllUsers()
}

// GetUserByID retrieves a user by their ID.
func (s *userService) GetUserByID(id uint64) (*models.User, error) {
	return s.userRepository.GetUserByID(id)
}

// AdminUpdateUser handles updating a user by an admin.
func (s *userService) AdminUpdateUser(id uint64, dto *dtos.UserAdminUpdateRequestDTO) error {
	user, err := s.userRepository.GetUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	if dto.Email != "" {
		user.Email = dto.Email
	}
	if dto.Username != "" {
		user.Username = dto.Username
	}
	if dto.PhoneNumber != "" {
		user.PhoneNumber = dto.PhoneNumber
	}
	if dto.Gender != "" {
		user.Gender = models.Gender(dto.Gender)
	}
	if dto.Region != "" {
		user.Region = dto.Region
	}
	return s.userRepository.UpdateUser(user)
}

// AdminDeleteUser handles deleting a user by an admin.
func (s *userService) AdminDeleteUser(id uint64) error {
	user, err := s.userRepository.GetUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	return s.userRepository.DeleteUser(user)
}

// RegisterUser handles a new user registration.
func (s *userService) RegisterUser(dto *dtos.UserRegisterRequestDTO) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Email:       dto.Email,
		Username:    dto.Username,
		Password:    string(hashedPassword),
		PhoneNumber: dto.PhoneNumber,
		Role:        "user",
	}
	return s.userRepository.CreateUser(user)
}

// LoginUser handles user authentication and generates new JWTs.
func (s *userService) LoginUser(dto *dtos.UserLoginRequestDTO) (*models.LoginToken, error) {
	logger := utils.GetLogger()
	logger.Info("START")

	user, err := s.userRepository.GetUserByEmail(dto.Email)
	if err != nil {
		return nil, errors.New("user not found or invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	// Generate new access and refresh tokens.
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	// Store the refresh token in the database for revocation.
	user.RefreshToken = refreshToken
	user.RefreshTokenExpiredAt = &expirationTime
	s.userRepository.UpdateUser(user)

	return &models.LoginToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ForgotPassword initiates the password reset process.
func (s *userService) ForgotPassword(dto *dtos.ForgotPasswordDTO) error {
	user, err := s.userRepository.GetUserByEmail(dto.Email)
	if err != nil {
		return errors.New("email not found")
	}

	// Generate a new 6-digit OTP
	otp := fmt.Sprintf("%06d", rand.Intn(999999))

	expirationTime := time.Now().Add(10 * time.Minute)
	user.OTP = otp
	user.OTPExpiredAt = &expirationTime

	err = s.sendOTPByEmail(user.Email, otp)
	if err != nil {
		fmt.Printf("failed to send OTP email to %s: %v\n", user.Email, err)
	}

	return s.userRepository.UpdateUser(user)
}

// sendOTPByEmail is a helper function to send the OTP.
func (s *userService) sendOTPByEmail(toEmail, otp string) error {
	// Ganti dengan konfigurasi SMTP server Anda
	smtpHost := s.cfg.SmtpHost
	smtpPort := s.cfg.SmtpPort
	fromEmail := s.cfg.SmtpUser
	password := s.cfg.SmtpPassword

	// Message body
	subject := "Subject: Your OTP Code\n"
	body := "Your OTP is: " + otp

	msg := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", fromEmail, password, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{toEmail}, msg)
}

func generateRandomString(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if n <= 0 {
		return "", errors.New("invalid token length")
	}
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	for i := range buf {
		buf[i] = letters[int(buf[i])%len(letters)]
	}
	return string(buf), nil
}

// VerifyOTP handles OTP verification.
func (s *userService) VerifyOTP(dto *dtos.VerifyOTPDTO) (*models.ResetToken, error) {
	user, err := s.userRepository.GetUserByEmail(dto.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.OTP != dto.OTP || time.Now().After(*user.OTPExpiredAt) {
		return nil, errors.New("invalid or expired OTP")
	}

	resetToken, err := generateRandomString(10)
	if err != nil {
		return nil, errors.New("failed to generate reset token: %w" + err.Error())
	}
	resetExp := time.Now().Add(10 * time.Minute)

	user.ResetToken = resetToken
	user.ResetTokenExpAt = &resetExp

	if err := s.userRepository.UpdateUser(user); err != nil {
		return nil, errors.New("failed to update reset token: %w" + err.Error())
	}

	return &models.ResetToken{
		ResetToken:          resetToken,
		ResetTokenExpiredAt: &resetExp,
	}, nil
}

// ResetPassword allows the user to set a new password.
func (s *userService) ResetPassword(dto *dtos.UserResetPasswordRequestDTO) error {
	user, err := s.userRepository.GetUserByEmail(dto.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// Validate the reset token and check for expiration.
	if user.ResetToken != dto.ResetToken || time.Now().After(*user.ResetTokenExpAt) {
		return errors.New("invalid or expired reset token")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.OTP = ""
	user.OTPExpiredAt = nil
	user.ResetToken = "" // Clear the reset token after use
	user.ResetTokenExpAt = nil
	return s.userRepository.UpdateUser(user)
}

// LogoutUser logs out the current user by invalidating their refresh token in the database.
func (s *userService) LogoutUser(userID uint64) error {
	user, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	user.RefreshToken = ""
	user.RefreshTokenExpiredAt = nil // Set to nil to properly clear the value in the database
	return s.userRepository.UpdateUser(user)
}

// RefreshToken generates a new access token using a valid refresh token.
func (s *userService) RefreshToken(tokenString string) (*models.LoginToken, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.cfg.JwtSecretKey, nil
	})
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepository.GetUserByRefreshToken(tokenString)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Generate new access token
	newAccessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, errors.New("failed to generate new access token")
	}

	return &models.LoginToken{
		AccessToken:  newAccessToken,
		RefreshToken: tokenString, // The refresh token itself is not changed here
	}, nil
}

func (s *userService) EditProfile(id uint64, dto *dtos.UserEditProfileRequestDTO) (*models.User, error) {
	user, err := s.userRepository.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if dto.Username != "" {
		user.Username = dto.Username
	}
	if dto.Email != "" {
		user.Email = dto.Email
	}
	if dto.Gender != "" {
		user.Gender = models.Gender(dto.Gender)
	}
	if dto.Birthday != nil {
		user.Birthday = dto.Birthday
	}
	if dto.Region != "" {
		user.Region = dto.Region
	}

	if err := s.userRepository.UpdateUser(user); err != nil {
		return nil, errors.New("failed to update user")
	}

	return user, nil
}
