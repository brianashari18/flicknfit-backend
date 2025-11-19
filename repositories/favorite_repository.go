package repositories

import (
	"flicknfit_backend/models"

	"gorm.io/gorm"
)

// FavoriteRepository defines methods for favorite data access
type FavoriteRepository interface {
	GetUserFavorites(userID uint64) ([]models.Favorite, error)
	AddFavorite(favorite *models.Favorite) error
	RemoveFavorite(userID, productItemID uint64) error
	IsFavorite(userID, productItemID uint64) (bool, error)
	GetFavoriteByID(id uint64) (*models.Favorite, error)
}

// favoriteRepository implements FavoriteRepository interface
type favoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository creates a new favorite repository
func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepository{db: db}
}

// GetUserFavorites retrieves all favorites for a user
func (r *favoriteRepository) GetUserFavorites(userID uint64) ([]models.Favorite, error) {
	var favorites []models.Favorite
	err := r.db.Preload("ProductItem").Preload("ProductItem.Product").
		Where("user_id = ?", userID).Find(&favorites).Error
	return favorites, err
}

// AddFavorite adds a product to user's favorites
func (r *favoriteRepository) AddFavorite(favorite *models.Favorite) error {
	return r.db.Create(favorite).Error
}

// RemoveFavorite removes a product from user's favorites
func (r *favoriteRepository) RemoveFavorite(userID, productItemID uint64) error {
	return r.db.Where("user_id = ? AND product_item_id = ?", userID, productItemID).
		Delete(&models.Favorite{}).Error
}

// IsFavorite checks if a product is in user's favorites
func (r *favoriteRepository) IsFavorite(userID, productItemID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&models.Favorite{}).
		Where("user_id = ? AND product_item_id = ?", userID, productItemID).
		Count(&count).Error
	return count > 0, err
}

// GetFavoriteByID retrieves a favorite by ID
func (r *favoriteRepository) GetFavoriteByID(id uint64) (*models.Favorite, error) {
	var favorite models.Favorite
	err := r.db.Preload("ProductItem").Preload("ProductItem.Product").
		First(&favorite, id).Error
	return &favorite, err
}
