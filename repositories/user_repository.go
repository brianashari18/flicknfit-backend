package repositories

import (
	"flicknfit_backend/models"

	"gorm.io/gorm"
)

// UserRepository defines the interface for data access operations on the User model.
//
//go:generate mockery --name UserRepository
type UserRepository interface {
	CreateUser(user *models.User) error
	GetAllUsers() ([]*models.User, error)
	GetUserByID(id uint64) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByRefreshToken(token string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(user *models.User) error
}

// userRepository is the implementation of UserRepository.
type userRepository struct {
	BaseRepository
}

// NewUserRepository creates and returns a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{BaseRepository{DB: db}}
}

// CreateUser creates a new user record in the database.
func (r *userRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *userRepository) GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	if err := r.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByID retrieves a user record by its ID.
func (r *userRepository) GetUserByID(id uint64) (*models.User, error) {
	var user models.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user record by their email.
func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByRefreshToken retrieves a user record by their stored refresh token.
func (r *userRepository) GetUserByRefreshToken(token string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("refresh_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user record in the database.
func (r *userRepository) UpdateUser(user *models.User) error {
	return r.DB.Save(user).Error
}

// DeleteUser deletes a user record from the database.
func (r *userRepository) DeleteUser(user *models.User) error {
	return r.DB.Delete(user).Error
}
