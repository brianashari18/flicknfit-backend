package services

import (
	"flicknfit_backend/dtos"
	"flicknfit_backend/errors"
	"flicknfit_backend/models"
	"flicknfit_backend/repositories"
	"strings"
)

// WardrobeService defines business logic for wardrobe management
type WardrobeService interface {
	GetUserWardrobe(userID uint64) (map[string][]models.UserWardrobe, error)
	GetWardrobeByCategory(userID uint64, category string) ([]models.UserWardrobe, error)
	CreateWardrobeItem(userID uint64, dto *dtos.CreateWardrobeItemDTO) error
	UpdateWardrobeItem(userID uint64, itemID uint64, dto *dtos.UpdateWardrobeItemDTO) error
	DeleteWardrobeItem(userID, itemID uint64) error
	GetWardrobeCategories(userID uint64) ([]string, error)
}

// wardrobeService implements WardrobeService interface
type wardrobeService struct {
	wardrobeRepo repositories.WardrobeRepository
}

// NewWardrobeService creates a new wardrobe service
func NewWardrobeService(wardrobeRepo repositories.WardrobeRepository) WardrobeService {
	return &wardrobeService{
		wardrobeRepo: wardrobeRepo,
	}
}

// GetUserWardrobe retrieves user's wardrobe organized by category
func (s *wardrobeService) GetUserWardrobe(userID uint64) (map[string][]models.UserWardrobe, error) {
	items, err := s.wardrobeRepo.GetUserWardrobe(userID)
	if err != nil {
		return nil, errors.NewDatabaseError("get user wardrobe", err)
	}

	// Group items by category
	wardrobe := make(map[string][]models.UserWardrobe)
	for _, item := range items {
		wardrobe[item.Category] = append(wardrobe[item.Category], item)
	}

	return wardrobe, nil
}

// GetWardrobeByCategory retrieves wardrobe items by specific category
func (s *wardrobeService) GetWardrobeByCategory(userID uint64, category string) ([]models.UserWardrobe, error) {
	// Validate category
	if !s.isValidCategory(category) {
		return nil, errors.NewValidationError("Invalid wardrobe category")
	}

	items, err := s.wardrobeRepo.GetWardrobeByCategory(userID, category)
	if err != nil {
		return nil, errors.NewDatabaseError("get wardrobe by category", err)
	}

	return items, nil
}

// CreateWardrobeItem creates a new wardrobe item
func (s *wardrobeService) CreateWardrobeItem(userID uint64, dto *dtos.CreateWardrobeItemDTO) error {
	// Validate category
	if !s.isValidCategory(dto.Category) {
		return errors.NewValidationError("Invalid wardrobe category")
	}

	item := &models.UserWardrobe{
		UserID:   userID,
		Category: strings.ToLower(dto.Category),
		ImageURL: dto.ImageURL,
	}

	err := s.wardrobeRepo.CreateWardrobeItem(item)
	if err != nil {
		return errors.NewDatabaseError("create wardrobe item", err)
	}

	return nil
}

// UpdateWardrobeItem updates an existing wardrobe item
func (s *wardrobeService) UpdateWardrobeItem(userID uint64, itemID uint64, dto *dtos.UpdateWardrobeItemDTO) error {
	// Get existing item
	item, err := s.wardrobeRepo.GetWardrobeItemByID(itemID)
	if err != nil {
		return errors.NewNotFoundError("Wardrobe item")
	}

	// Check ownership
	if item.UserID != userID {
		return errors.NewAuthorizationError("Cannot update another user's wardrobe item")
	}

	// Update fields
	if dto.Category != "" {
		if !s.isValidCategory(dto.Category) {
			return errors.NewValidationError("Invalid wardrobe category")
		}
		item.Category = strings.ToLower(dto.Category)
	}
	if dto.ImageURL != "" {
		item.ImageURL = dto.ImageURL
	}

	err = s.wardrobeRepo.UpdateWardrobeItem(item)
	if err != nil {
		return errors.NewDatabaseError("update wardrobe item", err)
	}

	return nil
}

// DeleteWardrobeItem deletes a wardrobe item
func (s *wardrobeService) DeleteWardrobeItem(userID, itemID uint64) error {
	// Get existing item to check ownership
	item, err := s.wardrobeRepo.GetWardrobeItemByID(itemID)
	if err != nil {
		return errors.NewNotFoundError("Wardrobe item")
	}

	// Check ownership
	if item.UserID != userID {
		return errors.NewAuthorizationError("Cannot delete another user's wardrobe item")
	}

	err = s.wardrobeRepo.DeleteWardrobeItem(userID, itemID)
	if err != nil {
		return errors.NewDatabaseError("delete wardrobe item", err)
	}

	return nil
}

// GetWardrobeCategories retrieves unique categories in user's wardrobe
func (s *wardrobeService) GetWardrobeCategories(userID uint64) ([]string, error) {
	categories, err := s.wardrobeRepo.GetWardrobeCategories(userID)
	if err != nil {
		return nil, errors.NewDatabaseError("get wardrobe categories", err)
	}

	return categories, nil
}

// isValidCategory validates wardrobe categories
func (s *wardrobeService) isValidCategory(category string) bool {
	validCategories := []string{
		"top", "bottom", "outerwear", "footwear", "accessories",
		"dress", "activewear", "underwear", "sleepwear", "formal",
	}

	category = strings.ToLower(category)
	for _, valid := range validCategories {
		if category == valid {
			return true
		}
	}
	return false
}
