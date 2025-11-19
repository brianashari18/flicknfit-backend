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

func TestReviewService_CreateReview(t *testing.T) {
	t.Run("should create review successfully", func(t *testing.T) {
		// Arrange
		mockReviewRepo := new(mocks.MockReviewRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewReviewService(mockReviewRepo, mockProductRepo)

		userID := uint64(1)
		product := testhelpers.CreateTestProduct()
		dto := &dtos.CreateReviewDTO{
			ProductID:  product.ID,
			Rating:     5,
			ReviewText: "Great product!",
		}

		mockProductRepo.On("GetProductByID", dto.ProductID).Return(&product, nil)
		mockReviewRepo.On("HasUserReviewedProduct", userID, dto.ProductID).Return(false, nil)
		mockReviewRepo.On("CreateReview", mock.AnythingOfType("*models.Review")).Return(nil)

		// Act
		err := service.CreateReview(userID, dto)

		// Assert
		assert.NoError(t, err)
		mockReviewRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when product not found", func(t *testing.T) {
		// Arrange
		mockReviewRepo := new(mocks.MockReviewRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewReviewService(mockReviewRepo, mockProductRepo)

		userID := uint64(1)
		dto := &dtos.CreateReviewDTO{
			ProductID:  999,
			Rating:     5,
			ReviewText: "Great product!",
		}

		mockProductRepo.On("GetProductByID", dto.ProductID).Return(nil, errors.New("not found"))

		// Act
		err := service.CreateReview(userID, dto)

		// Assert
		assert.Error(t, err)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("should return error when user already reviewed product", func(t *testing.T) {
		// Arrange
		mockReviewRepo := new(mocks.MockReviewRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewReviewService(mockReviewRepo, mockProductRepo)

		userID := uint64(1)
		product := testhelpers.CreateTestProduct()
		dto := &dtos.CreateReviewDTO{
			ProductID:  product.ID,
			Rating:     5,
			ReviewText: "Great product!",
		}

		mockProductRepo.On("GetProductByID", dto.ProductID).Return(&product, nil)
		mockReviewRepo.On("HasUserReviewedProduct", userID, dto.ProductID).Return(true, nil)

		// Act
		err := service.CreateReview(userID, dto)

		// Assert
		assert.Error(t, err)
		mockReviewRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestReviewService_GetProductReviews(t *testing.T) {
	t.Run("should get product reviews successfully", func(t *testing.T) {
		// Arrange
		mockReviewRepo := new(mocks.MockReviewRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewReviewService(mockReviewRepo, mockProductRepo)

		productID := uint64(1)
		page := 1
		limit := 10
		expectedReviews := []models.Review{testhelpers.CreateTestReview()}
		expectedCount := 1

		mockProductRepo.On("GetProductByID", productID).Return(&models.Product{ID: productID}, nil)
		mockReviewRepo.On("GetProductReviews", productID, limit, 0).Return(expectedReviews, nil) // offset = (page-1)*limit = (1-1)*10 = 0
		mockReviewRepo.On("GetProductReviewStats", productID).Return(map[string]interface{}{"average_rating": 4.5, "total_reviews": 1}, nil)

		// Act
		result, err := service.GetProductReviews(productID, page, limit)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, len(expectedReviews), len(result.Reviews))
		assert.Equal(t, expectedCount, result.Pagination.Total)
		mockReviewRepo.AssertExpectations(t)
	})

	t.Run("should handle repository error", func(t *testing.T) {
		// Arrange
		mockReviewRepo := new(mocks.MockReviewRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewReviewService(mockReviewRepo, mockProductRepo)

		productID := uint64(1)
		page := 1
		limit := 10
		expectedError := errors.New("database error")

		mockProductRepo.On("GetProductByID", productID).Return(&models.Product{ID: productID}, nil)
		mockReviewRepo.On("GetProductReviews", productID, limit, 0).Return([]models.Review{}, expectedError)

		// Act
		result, err := service.GetProductReviews(productID, page, limit)
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		mockReviewRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestReviewService_UpdateReview(t *testing.T) {
	t.Run("should update review successfully", func(t *testing.T) {
		// Arrange
		mockReviewRepo := new(mocks.MockReviewRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewReviewService(mockReviewRepo, mockProductRepo)

		userID := uint64(1)
		reviewID := uint64(1)
		existingReview := testhelpers.CreateTestReview()
		existingReview.UserID = userID
		dto := &dtos.UpdateReviewDTO{
			Rating:     4,
			ReviewText: "Updated review text",
		}

		mockReviewRepo.On("GetReviewByID", reviewID).Return(&existingReview, nil)
		mockReviewRepo.On("UpdateReview", mock.AnythingOfType("*models.Review")).Return(nil)

		// Act
		err := service.UpdateReview(userID, reviewID, dto)

		// Assert
		assert.NoError(t, err)
		mockReviewRepo.AssertExpectations(t)
	})

	t.Run("should return error when review not found", func(t *testing.T) {
		// Arrange
		mockReviewRepo := new(mocks.MockReviewRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewReviewService(mockReviewRepo, mockProductRepo)

		userID := uint64(1)
		reviewID := uint64(999)
		dto := &dtos.UpdateReviewDTO{
			Rating:     4,
			ReviewText: "Updated review text",
		}

		mockReviewRepo.On("GetReviewByID", reviewID).Return(nil, errors.New("not found"))

		// Act
		err := service.UpdateReview(userID, reviewID, dto)

		// Assert
		assert.Error(t, err)
		mockReviewRepo.AssertExpectations(t)
	})

	t.Run("should return error when user not authorized", func(t *testing.T) {
		// Arrange
		mockReviewRepo := new(mocks.MockReviewRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewReviewService(mockReviewRepo, mockProductRepo)

		userID := uint64(1)
		reviewID := uint64(1)
		existingReview := testhelpers.CreateTestReview()
		existingReview.UserID = uint64(2) // Different user
		dto := &dtos.UpdateReviewDTO{
			Rating:     4,
			ReviewText: "Updated review text",
		}

		mockReviewRepo.On("GetReviewByID", reviewID).Return(&existingReview, nil)

		// Act
		err := service.UpdateReview(userID, reviewID, dto)

		// Assert
		assert.Error(t, err)
		mockReviewRepo.AssertExpectations(t)
	})
}

func TestReviewService_DeleteReview(t *testing.T) {
	t.Run("should delete review successfully", func(t *testing.T) {
		// Arrange
		mockReviewRepo := new(mocks.MockReviewRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewReviewService(mockReviewRepo, mockProductRepo)

		userID := uint64(1)
		reviewID := uint64(1)
		existingReview := testhelpers.CreateTestReview()
		existingReview.UserID = userID

		mockReviewRepo.On("GetReviewByID", reviewID).Return(&existingReview, nil)
		mockReviewRepo.On("DeleteReview", reviewID).Return(nil)

		// Act
		err := service.DeleteReview(userID, reviewID)

		// Assert
		assert.NoError(t, err)
		mockReviewRepo.AssertExpectations(t)
	})

	t.Run("should return error when user not authorized", func(t *testing.T) {
		// Arrange
		mockReviewRepo := new(mocks.MockReviewRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		service := services.NewReviewService(mockReviewRepo, mockProductRepo)

		userID := uint64(1)
		reviewID := uint64(1)
		existingReview := testhelpers.CreateTestReview()
		existingReview.UserID = uint64(2) // Different user

		mockReviewRepo.On("GetReviewByID", reviewID).Return(&existingReview, nil)

		// Act
		err := service.DeleteReview(userID, reviewID)

		// Assert
		assert.Error(t, err)
		mockReviewRepo.AssertExpectations(t)
	})
}
