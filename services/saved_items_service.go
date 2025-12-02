package services

import (
	"errors"
	"flicknfit_backend/dtos"
	"flicknfit_backend/models"
	"flicknfit_backend/repositories"

	"gorm.io/gorm"
)

// SavedItemsService defines the interface for business logic related to shopping savedItemss.
type SavedItemsService interface {
	GetOrCreateSavedItems(userID uint64) (*models.SavedItems, error)
	AddProductItemToSavedItems(userID uint64, dto *dtos.AddProductItemToSavedItemsRequestDTO) (*models.SavedItems, error)
	UpdateProductItemInSavedItems(userID uint64, savedItemID uint64, dto *dtos.UpdateProductItemInSavedItemsRequestDTO) (*models.SavedItems, error)
	RemoveProductItemFromSavedItems(userID uint64, savedItemID uint64) error
	GetSavedItems(userID uint64) (*models.SavedItems, error)
}

// SavedItemsService is the implementation of SavedItemsService.
type savedItemsService struct {
	savedItemsRepository    repositories.SavedItemsRepository
	productRepository repositories.ProductRepository
}

// NewSavedItemsService creates and returns a new instance of SavedItemsService.
func NewSavedItemsService(savedItemsRepo repositories.SavedItemsRepository, productRepo repositories.ProductRepository) SavedItemsService {
	return &savedItemsService{
		savedItemsRepository:    savedItemsRepo,
		productRepository: productRepo,
	}
}

// GetOrCreateSavedItems retrieves a user's shopping savedItems or creates a new one if it doesn't exist.
func (s *savedItemsService) GetOrCreateSavedItems(userID uint64) (*models.SavedItems, error) {
	savedItems, err := s.savedItemsRepository.GetSavedItemsByUserID(userID)
	if err != nil {
		// If the savedItems doesn't exist, create a new one.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newsavedItems := &models.SavedItems{
				UserID: userID,
			}
			if err := s.savedItemsRepository.CreateSavedItems(newsavedItems); err != nil {
				return nil, errors.New("failed to create shopping savedItems")
			}
			return newsavedItems, nil
		}
		return nil, errors.New("failed to retrieve shopping savedItems")
	}
	return savedItems, nil
}

// GetSavedItems retrieves all items from a user's shopping savedItems.
func (s *savedItemsService) GetSavedItems(userID uint64) (*models.SavedItems, error) {
	savedItems, err := s.savedItemsRepository.GetSavedItemsByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to retrieve shopping savedItems")
	}
	return savedItems, nil
}

// AddProductItemToSavedItems handles adding a product item to the user's savedItems.
func (s *savedItemsService) AddProductItemToSavedItems(userID uint64, dto *dtos.AddProductItemToSavedItemsRequestDTO) (*models.SavedItems, error) {
	savedItems, err := s.GetOrCreateSavedItems(userID)
	if err != nil {
		return nil, err
	}

	productItem, err := s.productRepository.GetProductItemByID(dto.ProductItemID)
	if productItem == nil && err != nil {
		return nil, errors.New("product item not found")
	}

	// Check if item already exists in the savedItems.
	var SavedItemsList *models.SavedItemsList
	for _, item := range savedItems.SavedItemsList {
		if item.ProductItemID == dto.ProductItemID {
			SavedItemsList = &item
			break
		}
	}

	if SavedItemsList != nil {
		// Update quantity if item exists.
		SavedItemsList.Quantity += dto.Quantity
		if err := s.savedItemsRepository.UpdateSavedItem(SavedItemsList); err != nil {
			return nil, errors.New("failed to update savedItems item")
		}
	} else {
		// Add new item to savedItems.
		SavedItemsList = &models.SavedItemsList{
			SavedItemsID: savedItems.ID,
			ProductItemID:  dto.ProductItemID,
			Quantity:       dto.Quantity,
		}
		if err := s.savedItemsRepository.AddSavedItem(SavedItemsList); err != nil {
			return nil, errors.New("failed to add item to savedItems")
		}
	}

	// Retrieve the updated savedItems with the new or updated item.
	return s.savedItemsRepository.GetSavedItemsByUserID(userID)
}

// UpdateProductItemInSavedItems handles updating the quantity of a product item in the user's savedItems.
func (s *savedItemsService) UpdateProductItemInSavedItems(userID uint64, savedItemID uint64, dto *dtos.UpdateProductItemInSavedItemsRequestDTO) (*models.SavedItems, error) {
	savedItems, err := s.savedItemsRepository.GetSavedItemsByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to retrieve shopping savedItems")
	}

	SavedItemsList, err := s.savedItemsRepository.GetSavedItemByItemID(savedItemID)
	if err != nil {
		return nil, errors.New("savedItems item not found")
	}

	// Verify the savedItems item belongs to the user's savedItems.
	if SavedItemsList.SavedItemsID != savedItems.ID {
		return nil, errors.New("savedItems item does not belong to this user")
	}

	SavedItemsList.Quantity = dto.Quantity
	if err := s.savedItemsRepository.UpdateSavedItem(SavedItemsList); err != nil {
		return nil, errors.New("failed to update savedItems item quantity")
	}

	return s.savedItemsRepository.GetSavedItemsByUserID(userID)
}

// RemoveProductItemFromSavedItems handles removing an item from the user's savedItems.
func (s *savedItemsService) RemoveProductItemFromSavedItems(userID uint64, savedItemID uint64) error {
	savedItems, err := s.savedItemsRepository.GetSavedItemsByUserID(userID)
	if err != nil {
		return errors.New("failed to retrieve shopping savedItems")
	}

	SavedItemsList, err := s.savedItemsRepository.GetSavedItemByItemID(savedItemID)
	if err != nil {
		return errors.New("savedItems item not found")
	}

	// Verify the savedItems item belongs to the user's savedItems.
	if SavedItemsList.SavedItemsID != savedItems.ID {
		return errors.New("savedItems item does not belong to this user")
	}

	return s.savedItemsRepository.DeleteSavedItem(savedItemID)
}
