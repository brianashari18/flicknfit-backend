package services

import (
	"flicknfit_backend/dtos"
	"flicknfit_backend/errors"
	"flicknfit_backend/models"
	"flicknfit_backend/repositories"
)

// FavoriteService defines business logic for favorites
type FavoriteService interface {
	GetUserFavorites(userID uint64) ([]models.Favorite, error)
	AddFavorite(userID uint64, dto *dtos.AddFavoriteDTO) error
	RemoveFavorite(userID, productItemID uint64) error
	ToggleFavorite(userID uint64, dto *dtos.AddFavoriteDTO) (*dtos.FavoriteToggleResponseDTO, error)
}

// favoriteService implements FavoriteService interface
type favoriteService struct {
	favoriteRepo repositories.FavoriteRepository
	productRepo  repositories.ProductRepository
}

// NewFavoriteService creates a new favorite service
func NewFavoriteService(favoriteRepo repositories.FavoriteRepository, productRepo repositories.ProductRepository) FavoriteService {
	return &favoriteService{
		favoriteRepo: favoriteRepo,
		productRepo:  productRepo,
	}
}

// GetUserFavorites retrieves all user's favorites
func (s *favoriteService) GetUserFavorites(userID uint64) ([]models.Favorite, error) {
	favorites, err := s.favoriteRepo.GetUserFavorites(userID)
	if err != nil {
		return nil, errors.NewDatabaseError("get user favorites", err)
	}
	return favorites, nil
}

// AddFavorite adds a product to user's favorites
func (s *favoriteService) AddFavorite(userID uint64, dto *dtos.AddFavoriteDTO) error {
	// Check if product item exists
	_, err := s.productRepo.GetProductItemByID(dto.ProductItemID)
	if err != nil {
		return errors.NewNotFoundError("Product item")
	}
	// Check if already favorited
	isFav, err := s.favoriteRepo.IsFavorite(userID, dto.ProductItemID)
	if err != nil {
		return errors.NewDatabaseError("check favorite status", err)
	}
	if isFav {
		return errors.NewConflictError("Favorite")
	}

	favorite := &models.Favorite{
		UserID:        userID,
		ProductItemID: dto.ProductItemID,
	}

	err = s.favoriteRepo.AddFavorite(favorite)
	if err != nil {
		return errors.NewDatabaseError("add favorite", err)
	}

	return nil
}

// RemoveFavorite removes a product from user's favorites
func (s *favoriteService) RemoveFavorite(userID, productItemID uint64) error {
	// Check if favorite exists
	isFav, err := s.favoriteRepo.IsFavorite(userID, productItemID)
	if err != nil {
		return errors.NewDatabaseError("check favorite status", err)
	}
	if !isFav {
		return errors.NewNotFoundError("Favorite")
	}

	err = s.favoriteRepo.RemoveFavorite(userID, productItemID)
	if err != nil {
		return errors.NewDatabaseError("remove favorite", err)
	}

	return nil
}

// ToggleFavorite toggles favorite status for a product
func (s *favoriteService) ToggleFavorite(userID uint64, dto *dtos.AddFavoriteDTO) (*dtos.FavoriteToggleResponseDTO, error) {
	// Check if product item exists
	_, err := s.productRepo.GetProductItemByID(dto.ProductItemID)
	if err != nil {
		return nil, errors.NewNotFoundError("Product item")
	}

	// Check current favorite status
	isFav, err := s.favoriteRepo.IsFavorite(userID, dto.ProductItemID)
	if err != nil {
		return nil, errors.NewDatabaseError("check favorite status", err)
	}

	var message string
	var isFavorited bool

	if isFav {
		// Remove from favorites
		err = s.favoriteRepo.RemoveFavorite(userID, dto.ProductItemID)
		if err != nil {
			return nil, errors.NewDatabaseError("remove favorite", err)
		}
		message = "Removed from favorites"
		isFavorited = false
	} else {
		// Add to favorites
		favorite := &models.Favorite{
			UserID:        userID,
			ProductItemID: dto.ProductItemID,
		}
		err = s.favoriteRepo.AddFavorite(favorite)
		if err != nil {
			return nil, errors.NewDatabaseError("add favorite", err)
		}
		message = "Added to favorites"
		isFavorited = true
	}

	response := &dtos.FavoriteToggleResponseDTO{
		ProductItemID: dto.ProductItemID,
		IsFavorited:   isFavorited,
		Message:       message,
	}

	return response, nil
}
