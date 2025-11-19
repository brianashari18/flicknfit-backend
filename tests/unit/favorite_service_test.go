package unit

import (
	"errors"
	"flicknfit_backend/dtos"
	"flicknfit_backend/models"
	"flicknfit_backend/services"
	"flicknfit_backend/tests/mocks"
	"flicknfit_backend/tests/testhelpers"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFavoriteService_GetUserFavorites(t *testing.T) {
	t.Run("should return user favorites successfully", func(t *testing.T) {
		// Arrange
		mockFavoriteRepo := new(mocks.MockFavoriteRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewFavoriteService(mockFavoriteRepo, mockProductRepo)

		userID := uint64(1)
		expectedFavorites := []models.Favorite{testhelpers.CreateTestFavorite()}

		mockFavoriteRepo.On("GetUserFavorites", userID).Return(expectedFavorites, nil)

		// Act
		result, err := service.GetUserFavorites(userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedFavorites, result)
		mockFavoriteRepo.AssertExpectations(t)
	})

	t.Run("should handle repository error", func(t *testing.T) {
		// Arrange
		mockFavoriteRepo := new(mocks.MockFavoriteRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewFavoriteService(mockFavoriteRepo, mockProductRepo)

		userID := uint64(1)
		expectedError := errors.New("database error")

		mockFavoriteRepo.On("GetUserFavorites", userID).Return([]models.Favorite{}, expectedError)

		// Act
		result, err := service.GetUserFavorites(userID)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
		mockFavoriteRepo.AssertExpectations(t)
	})
}

func TestFavoriteService_AddFavorite(t *testing.T) {
	t.Run("should add favorite successfully", func(t *testing.T) {
		// Arrange
		mockFavoriteRepo := new(mocks.MockFavoriteRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewFavoriteService(mockFavoriteRepo, mockProductRepo)

		userID := uint64(1)
		productItem := testhelpers.CreateTestProductItem()
		dto := &dtos.AddFavoriteDTO{ProductItemID: productItem.ID}

		mockProductRepo.On("GetProductItemByID", dto.ProductItemID).Return(&productItem, nil)
		mockFavoriteRepo.On("IsFavorite", userID, dto.ProductItemID).Return(false, nil)
		mockFavoriteRepo.On("AddFavorite", mock.AnythingOfType("*models.Favorite")).Return(nil)

		// Act
		err := service.AddFavorite(userID, dto)

		// Assert
		assert.NoError(t, err)
		mockFavoriteRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when product item not found", func(t *testing.T) {
		// Arrange
		mockFavoriteRepo := new(mocks.MockFavoriteRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewFavoriteService(mockFavoriteRepo, mockProductRepo)

		userID := uint64(1)
		dto := &dtos.AddFavoriteDTO{ProductItemID: 999}

		mockProductRepo.On("GetProductItemByID", dto.ProductItemID).Return(nil, errors.New("not found"))

		// Act
		err := service.AddFavorite(userID, dto)

		// Assert
		assert.Error(t, err)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when product already favorited", func(t *testing.T) {
		// Arrange
		mockFavoriteRepo := new(mocks.MockFavoriteRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewFavoriteService(mockFavoriteRepo, mockProductRepo)

		userID := uint64(1)
		productItem := testhelpers.CreateTestProductItem()
		dto := &dtos.AddFavoriteDTO{ProductItemID: productItem.ID}

		mockProductRepo.On("GetProductItemByID", dto.ProductItemID).Return(&productItem, nil)
		mockFavoriteRepo.On("IsFavorite", userID, dto.ProductItemID).Return(true, nil)

		// Act
		err := service.AddFavorite(userID, dto)

		// Assert
		assert.Error(t, err)
		mockFavoriteRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestFavoriteService_RemoveFavorite(t *testing.T) {
	t.Run("should remove favorite successfully", func(t *testing.T) {
		// Arrange
		mockFavoriteRepo := new(mocks.MockFavoriteRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewFavoriteService(mockFavoriteRepo, mockProductRepo)

		userID := uint64(1)
		productItemID := uint64(1)

		mockFavoriteRepo.On("IsFavorite", userID, productItemID).Return(true, nil)
		mockFavoriteRepo.On("RemoveFavorite", userID, productItemID).Return(nil)

		// Act
		err := service.RemoveFavorite(userID, productItemID)

		// Assert
		assert.NoError(t, err)
		mockFavoriteRepo.AssertExpectations(t)
	})

	t.Run("should return error when favorite not found", func(t *testing.T) {
		// Arrange
		mockFavoriteRepo := new(mocks.MockFavoriteRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewFavoriteService(mockFavoriteRepo, mockProductRepo)

		userID := uint64(1)
		productItemID := uint64(999)

		mockFavoriteRepo.On("IsFavorite", userID, productItemID).Return(false, nil)

		// Act
		err := service.RemoveFavorite(userID, productItemID)

		// Assert
		assert.Error(t, err)
		mockFavoriteRepo.AssertExpectations(t)
	})
}

func TestFavoriteService_ToggleFavorite(t *testing.T) {
	t.Run("should add favorite when not favorited", func(t *testing.T) {
		// Arrange
		mockFavoriteRepo := new(mocks.MockFavoriteRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewFavoriteService(mockFavoriteRepo, mockProductRepo)

		userID := uint64(1)
		productItem := testhelpers.CreateTestProductItem()
		dto := &dtos.AddFavoriteDTO{ProductItemID: productItem.ID}

		mockProductRepo.On("GetProductItemByID", dto.ProductItemID).Return(&productItem, nil)
		mockFavoriteRepo.On("IsFavorite", userID, dto.ProductItemID).Return(false, nil)
		mockFavoriteRepo.On("AddFavorite", mock.AnythingOfType("*models.Favorite")).Return(nil)

		// Act
		result, err := service.ToggleFavorite(userID, dto)

		// Assert
		assert.NoError(t, err)
		assert.True(t, result.IsFavorited)
		assert.Equal(t, "Added to favorites", result.Message)
		assert.Equal(t, dto.ProductItemID, result.ProductItemID)
		mockFavoriteRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should remove favorite when already favorited", func(t *testing.T) {
		// Arrange
		mockFavoriteRepo := new(mocks.MockFavoriteRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewFavoriteService(mockFavoriteRepo, mockProductRepo)

		userID := uint64(1)
		productItem := testhelpers.CreateTestProductItem()
		dto := &dtos.AddFavoriteDTO{ProductItemID: productItem.ID}

		mockProductRepo.On("GetProductItemByID", dto.ProductItemID).Return(&productItem, nil)
		mockFavoriteRepo.On("IsFavorite", userID, dto.ProductItemID).Return(true, nil)
		mockFavoriteRepo.On("RemoveFavorite", userID, dto.ProductItemID).Return(nil)

		// Act
		result, err := service.ToggleFavorite(userID, dto)

		// Assert
		assert.NoError(t, err)
		assert.False(t, result.IsFavorited)
		assert.Equal(t, "Removed from favorites", result.Message)
		assert.Equal(t, dto.ProductItemID, result.ProductItemID)
		mockFavoriteRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})
}
