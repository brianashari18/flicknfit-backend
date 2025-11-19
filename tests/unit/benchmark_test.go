package unit

import (
	"flicknfit_backend/dtos"
	"flicknfit_backend/models"
	"flicknfit_backend/services"
	"flicknfit_backend/tests/mocks"
	"flicknfit_backend/tests/testhelpers"
	"testing"

	"github.com/stretchr/testify/mock"
)

// BenchmarkFavoriteService_GetUserFavorites benchmarks the GetUserFavorites method
func BenchmarkFavoriteService_GetUserFavorites(b *testing.B) {
	// Setup
	mockFavoriteRepo := new(mocks.MockFavoriteRepository)
	mockProductRepo := new(mocks.MockProductRepository)
	service := services.NewFavoriteService(mockFavoriteRepo, mockProductRepo)

	userID := uint64(1)
	favorites := make([]models.Favorite, 100) // Create 100 favorites for benchmarking
	for i := 0; i < 100; i++ {
		favorite := testhelpers.CreateTestFavorite()
		favorite.ID = uint64(i + 1)
		favorites[i] = favorite
	}

	mockFavoriteRepo.On("GetUserFavorites", userID).Return(favorites, nil)

	// Reset timer and run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetUserFavorites(userID)
	}
}

// BenchmarkFavoriteService_ToggleFavorite benchmarks the ToggleFavorite method
func BenchmarkFavoriteService_ToggleFavorite(b *testing.B) {
	// Setup
	mockFavoriteRepo := new(mocks.MockFavoriteRepository)
	mockProductRepo := new(mocks.MockProductRepository)
	service := services.NewFavoriteService(mockFavoriteRepo, mockProductRepo)

	userID := uint64(1)
	productItem := testhelpers.CreateTestProductItem()
	dto := &dtos.AddFavoriteDTO{ProductItemID: productItem.ID}

	mockProductRepo.On("GetProductItemByID", dto.ProductItemID).Return(&productItem, nil)
	mockFavoriteRepo.On("IsFavorite", userID, dto.ProductItemID).Return(false, nil)
	mockFavoriteRepo.On("AddFavorite", mock.AnythingOfType("*models.Favorite")).Return(nil)

	// Reset timer and run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ToggleFavorite(userID, dto)
	}
}

// BenchmarkReviewService_GetProductReviews benchmarks the GetProductReviews method
func BenchmarkReviewService_GetProductReviews(b *testing.B) {
	// Setup
	mockReviewRepo := new(mocks.MockReviewRepository)
	mockProductRepo := new(mocks.MockProductRepository)
	service := services.NewReviewService(mockReviewRepo, mockProductRepo)

	productID := uint64(1)
	page := 1
	limit := 10

	reviews := make([]models.Review, limit)
	for i := 0; i < limit; i++ {
		review := testhelpers.CreateTestReview()
		review.ID = uint64(i + 1)
		reviews[i] = review
	}

	mockReviewRepo.On("GetReviewsByProductID", productID, page, limit).Return(reviews, nil)
	mockReviewRepo.On("CountReviewsByProductID", productID).Return(100, nil)

	// Reset timer and run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetProductReviews(productID, page, limit)
	}
}

// BenchmarkDTOConversion benchmarks DTO conversion functions
func BenchmarkDTOConversion_ToFavoriteResponseDTO(b *testing.B) {
	favorite := testhelpers.CreateTestFavorite()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dtos.ToFavoriteResponseDTO(favorite)
	}
}

func BenchmarkDTOConversion_ToFavoriteResponseDTOs(b *testing.B) {
	favorites := make([]models.Favorite, 100)
	for i := 0; i < 100; i++ {
		favorite := testhelpers.CreateTestFavorite()
		favorite.ID = uint64(i + 1)
		favorites[i] = favorite
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dtos.ToFavoriteResponseDTOs(favorites)
	}
}

func BenchmarkDTOConversion_ToReviewResponseDTO(b *testing.B) {
	review := testhelpers.CreateTestReview()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dtos.ToReviewResponseDTO(&review)
	}
}

func BenchmarkDTOConversion_ToReviewResponseDTOs(b *testing.B) {
	reviews := make([]models.Review, 100)
	for i := 0; i < 100; i++ {
		review := testhelpers.CreateTestReview()
		review.ID = uint64(i + 1)
		reviews[i] = review
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dtos.ToReviewResponseDTOs(reviews)
	}
}
