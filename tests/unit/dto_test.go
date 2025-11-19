package unit

import (
	"flicknfit_backend/dtos"
	"flicknfit_backend/models"
	"flicknfit_backend/tests/testhelpers"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToFavoriteResponseDTO(t *testing.T) {
	t.Run("should convert favorite model to response DTO", func(t *testing.T) {
		// Arrange
		favorite := testhelpers.CreateTestFavorite()

		// Act
		result := dtos.ToFavoriteResponseDTO(favorite)

		// Assert
		assert.Equal(t, favorite.ID, result.ID)
		assert.Equal(t, favorite.ProductItemID, result.ProductItemID)
		assert.Equal(t, favorite.CreatedAt, result.CreatedAt)
		assert.Equal(t, favorite.ProductItem.ID, result.ProductItem.ID)
		assert.Equal(t, favorite.ProductItem.SKU, result.ProductItem.SKU)
		assert.Equal(t, favorite.ProductItem.Price, result.ProductItem.Price)
	})
}

func TestToFavoriteResponseDTOs(t *testing.T) {
	t.Run("should convert slice of favorites to response DTOs", func(t *testing.T) {
		// Arrange
		favorites := []models.Favorite{
			testhelpers.CreateTestFavorite(),
		}

		secondFavorite := testhelpers.CreateTestFavorite()
		secondFavorite.ID = 2
		secondFavorite.ProductItemID = 2
		favorites = append(favorites, secondFavorite)

		// Act
		result := dtos.ToFavoriteResponseDTOs(favorites)

		// Assert
		assert.Len(t, result, 2)
		assert.Equal(t, favorites[0].ID, result[0].ID)
		assert.Equal(t, favorites[1].ID, result[1].ID)
	})

	t.Run("should return empty slice for empty input", func(t *testing.T) {
		// Arrange
		favorites := []models.Favorite{}

		// Act
		result := dtos.ToFavoriteResponseDTOs(favorites)

		// Assert
		assert.Empty(t, result)
	})
}

func TestToProductItemResponseDTO(t *testing.T) {
	t.Run("should convert product item model to response DTO", func(t *testing.T) {
		// Arrange
		productItem := testhelpers.CreateTestProductItem()

		// Act
		result := dtos.ToProductItemResponseDTO(productItem)

		// Assert
		assert.Equal(t, productItem.ID, result.ID)
		assert.Equal(t, productItem.SKU, result.SKU)
		assert.Equal(t, productItem.Price, result.Price)
		assert.Equal(t, productItem.Stock, result.Stock)
		assert.Equal(t, productItem.PhotoURL, result.PhotoURL)
		assert.Equal(t, productItem.Product.ID, result.Product.ID)
		assert.Equal(t, productItem.Product.Name, result.Product.Name)
	})
}

func TestToReviewResponseDTO(t *testing.T) {
	t.Run("should convert review model to response DTO", func(t *testing.T) {
		// Arrange
		review := testhelpers.CreateTestReview()

		// Act
		result := dtos.ToReviewResponseDTO(&review)

		// Assert
		assert.Equal(t, review.ID, result.ID)
		assert.Equal(t, review.UserID, result.UserID)
		assert.Equal(t, review.ProductID, result.ProductID)
		assert.Equal(t, review.Rating, result.Rating)
		assert.Equal(t, review.ReviewText, result.ReviewText)
		assert.WithinDuration(t, review.CreatedAt, result.CreatedAt, time.Second)
	})
}

func TestToReviewResponseDTOs(t *testing.T) {
	t.Run("should convert slice of reviews to response DTOs", func(t *testing.T) {
		// Arrange
		reviews := []models.Review{
			testhelpers.CreateTestReview(),
		}

		secondReview := testhelpers.CreateTestReview()
		secondReview.ID = 2
		secondReview.Rating = 4
		reviews = append(reviews, secondReview)

		// Act
		result := dtos.ToReviewResponseDTOs(reviews)

		// Assert
		assert.Len(t, result, 2)
		assert.Equal(t, reviews[0].ID, result[0].ID)
		assert.Equal(t, reviews[0].Rating, result[0].Rating)
		assert.Equal(t, reviews[1].ID, result[1].ID)
		assert.Equal(t, reviews[1].Rating, result[1].Rating)
	})
}

func TestFavoriteToggleResponseDTO(t *testing.T) {
	t.Run("should create toggle response for adding favorite", func(t *testing.T) {
		// Arrange & Act
		response := dtos.FavoriteToggleResponseDTO{
			ProductItemID: 1,
			IsFavorited:   true,
			Message:       "Added to favorites",
		}

		// Assert
		assert.Equal(t, uint64(1), response.ProductItemID)
		assert.True(t, response.IsFavorited)
		assert.Equal(t, "Added to favorites", response.Message)
	})

	t.Run("should create toggle response for removing favorite", func(t *testing.T) {
		// Arrange & Act
		response := dtos.FavoriteToggleResponseDTO{
			ProductItemID: 1,
			IsFavorited:   false,
			Message:       "Removed from favorites",
		}

		// Assert
		assert.Equal(t, uint64(1), response.ProductItemID)
		assert.False(t, response.IsFavorited)
		assert.Equal(t, "Removed from favorites", response.Message)
	})
}
