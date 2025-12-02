package repositories

import (
	"flicknfit_backend/models"

	"gorm.io/gorm"
)

// SavedItemsRepository defines the interface for data access operations on the SavedItems model.
//
//go:generate mockery --name SavedItemsRepository
type SavedItemsRepository interface {
	GetSavedItemsByUserID(userID uint64) (*models.SavedItems, error)
	CreateSavedItems(cart *models.SavedItems) error
	GetSavedItemByID(cartID uint64, productItemID uint64) (*models.SavedItemsList, error)
	AddSavedItem(item *models.SavedItemsList) error
	UpdateSavedItem(item *models.SavedItemsList) error
	DeleteSavedItem(id uint64) error
	GetSavedItemByItemID(id uint64) (*models.SavedItemsList, error)
	SaveSavedItems(cart *models.SavedItems) error
}

// SavedItemsRepository is the implementation of SavedItemsRepository.
type savedItemsRepository struct {
	BaseRepository
}

// NewSavedItemsRepository creates and returns a new instance of SavedItemsRepository.
func NewSavedItemsRepository(db *gorm.DB) SavedItemsRepository {
	return &savedItemsRepository{BaseRepository: BaseRepository{DB: db}}
}

// GetSavedItemsByUserID retrieves a shopping cart by its user ID.
func (r *savedItemsRepository) GetSavedItemsByUserID(userID uint64) (*models.SavedItems, error) {
	var cart models.SavedItems
	err := r.DB.Where("user_id = ?", userID).Preload("SavedItemsList.ProductItem").First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// CreateSavedItems creates a new shopping cart record in the database.
func (r *savedItemsRepository) CreateSavedItems(cart *models.SavedItems) error {
	return r.DB.Create(cart).Error
}

// GetSavedItemByID retrieves a shopping cart item by its cart ID and product item ID.
func (r *savedItemsRepository) GetSavedItemByID(cartID uint64, productItemID uint64) (*models.SavedItemsList, error) {
	var item models.SavedItemsList
	err := r.DB.Where("saved_items_id = ? AND product_item_id = ?", cartID, productItemID).Preload("ProductItem").First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// AddSavedItem adds a new item to the shopping cart.
func (r *savedItemsRepository) AddSavedItem(item *models.SavedItemsList) error {
	return r.DB.Create(item).Error
}

// UpdateSavedItem updates an existing item in the shopping cart.
func (r *savedItemsRepository) UpdateSavedItem(item *models.SavedItemsList) error {
	return r.DB.Save(item).Error
}

// DeleteSavedItem deletes an item from the shopping cart by its ID.
func (r *savedItemsRepository) DeleteSavedItem(id uint64) error {
	return r.DB.Delete(&models.SavedItemsList{}, id).Error
}

func (r *savedItemsRepository) GetSavedItemByItemID(id uint64) (*models.SavedItemsList, error) {
	var item models.SavedItemsList
	err := r.DB.Preload("ProductItem").First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *savedItemsRepository) SaveSavedItems(cart *models.SavedItems) error {
	return r.DB.Save(cart).Error
}
