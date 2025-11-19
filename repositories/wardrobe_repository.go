package repositories

import (
	"flicknfit_backend/models"

	"gorm.io/gorm"
)

// WardrobeRepository defines methods for wardrobe data access
type WardrobeRepository interface {
	GetUserWardrobe(userID uint64) ([]models.UserWardrobe, error)
	GetWardrobeByCategory(userID uint64, category string) ([]models.UserWardrobe, error)
	CreateWardrobeItem(item *models.UserWardrobe) error
	UpdateWardrobeItem(item *models.UserWardrobe) error
	DeleteWardrobeItem(userID, itemID uint64) error
	GetWardrobeItemByID(id uint64) (*models.UserWardrobe, error)
	GetWardrobeCategories(userID uint64) ([]string, error)
}

// wardrobeRepository implements WardrobeRepository interface
type wardrobeRepository struct {
	db *gorm.DB
}

// NewWardrobeRepository creates a new wardrobe repository
func NewWardrobeRepository(db *gorm.DB) WardrobeRepository {
	return &wardrobeRepository{db: db}
}

// GetUserWardrobe retrieves all wardrobe items for a user
func (r *wardrobeRepository) GetUserWardrobe(userID uint64) ([]models.UserWardrobe, error) {
	var items []models.UserWardrobe
	err := r.db.Where("user_id = ?", userID).
		Order("category ASC, created_at DESC").
		Find(&items).Error
	return items, err
}

// GetWardrobeByCategory retrieves wardrobe items by category for a user
func (r *wardrobeRepository) GetWardrobeByCategory(userID uint64, category string) ([]models.UserWardrobe, error) {
	var items []models.UserWardrobe
	err := r.db.Where("user_id = ? AND category = ?", userID, category).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

// CreateWardrobeItem creates a new wardrobe item
func (r *wardrobeRepository) CreateWardrobeItem(item *models.UserWardrobe) error {
	return r.db.Create(item).Error
}

// UpdateWardrobeItem updates an existing wardrobe item
func (r *wardrobeRepository) UpdateWardrobeItem(item *models.UserWardrobe) error {
	return r.db.Save(item).Error
}

// DeleteWardrobeItem deletes a wardrobe item
func (r *wardrobeRepository) DeleteWardrobeItem(userID, itemID uint64) error {
	return r.db.Where("user_id = ? AND id = ?", userID, itemID).
		Delete(&models.UserWardrobe{}).Error
}

// GetWardrobeItemByID retrieves a wardrobe item by ID
func (r *wardrobeRepository) GetWardrobeItemByID(id uint64) (*models.UserWardrobe, error) {
	var item models.UserWardrobe
	err := r.db.First(&item, id).Error
	return &item, err
}

// GetWardrobeCategories retrieves unique categories in user's wardrobe
func (r *wardrobeRepository) GetWardrobeCategories(userID uint64) ([]string, error) {
	var categories []string
	err := r.db.Model(&models.UserWardrobe{}).
		Select("DISTINCT category").
		Where("user_id = ?", userID).
		Order("category ASC").
		Pluck("category", &categories).Error
	return categories, err
}
