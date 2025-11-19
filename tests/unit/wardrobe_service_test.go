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

func TestWardrobeService_GetUserWardrobe(t *testing.T) {
	t.Run("should return user wardrobe grouped by category", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockWardrobeRepository)
		service := services.NewWardrobeService(mockRepo)

		userID := uint64(1)
		wardrobeItems := []models.UserWardrobe{
			testhelpers.CreateTestUserWardrobe(),
		}

		mockRepo.On("GetUserWardrobe", userID).Return(wardrobeItems, nil)

		// Act
		result, err := service.GetUserWardrobe(userID)

		// Assert
		assert.NoError(t, err)
		assert.Contains(t, result, "Shirts")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should handle repository error", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockWardrobeRepository)
		service := services.NewWardrobeService(mockRepo)

		userID := uint64(1)
		expectedError := errors.New("database error")

		mockRepo.On("GetUserWardrobe", userID).Return([]models.UserWardrobe{}, expectedError)

		// Act
		result, err := service.GetUserWardrobe(userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestWardrobeService_CreateWardrobeItem(t *testing.T) {
	t.Run("should create wardrobe item successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockWardrobeRepository)
		service := services.NewWardrobeService(mockRepo)

		userID := uint64(1)
		dto := &dtos.CreateWardrobeItemDTO{
			Category: "Shirts",
			ImageURL: "https://example.com/shirt.jpg",
		}

		mockRepo.On("CreateWardrobeItem", mock.AnythingOfType("*models.UserWardrobe")).Return(nil)

		// Act
		err := service.CreateWardrobeItem(userID, dto)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestWardrobeService_DeleteWardrobeItem(t *testing.T) {
	t.Run("should delete wardrobe item successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.MockWardrobeRepository)
		service := services.NewWardrobeService(mockRepo)

		userID := uint64(1)
		itemID := uint64(1)
		existingItem := testhelpers.CreateTestUserWardrobe()
		existingItem.UserID = userID

		mockRepo.On("GetWardrobeItemByID", itemID).Return(&existingItem, nil)
		mockRepo.On("DeleteWardrobeItem", userID, itemID).Return(nil)

		// Act
		err := service.DeleteWardrobeItem(userID, itemID)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
