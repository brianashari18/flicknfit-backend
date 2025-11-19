package mocks

import (
	"flicknfit_backend/models"

	"github.com/stretchr/testify/mock"
)

// MockReviewRepository is a mock implementation of ReviewRepository
type MockReviewRepository struct {
	mock.Mock
}

func (m *MockReviewRepository) GetProductReviews(productID uint64, limit, offset int) ([]models.Review, error) {
	args := m.Called(productID, limit, offset)
	return args.Get(0).([]models.Review), args.Error(1)
}

func (m *MockReviewRepository) CreateReview(review *models.Review) error {
	args := m.Called(review)
	return args.Error(0)
}

func (m *MockReviewRepository) UpdateReview(review *models.Review) error {
	args := m.Called(review)
	return args.Error(0)
}

func (m *MockReviewRepository) DeleteReview(reviewID uint64) error {
	args := m.Called(reviewID)
	return args.Error(0)
}

func (m *MockReviewRepository) GetReviewByID(reviewID uint64) (*models.Review, error) {
	args := m.Called(reviewID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Review), args.Error(1)
}

func (m *MockReviewRepository) GetUserReviews(userID uint64) ([]models.Review, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Review), args.Error(1)
}

func (m *MockReviewRepository) GetProductReviewStats(productID uint64) (map[string]interface{}, error) {
	args := m.Called(productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockReviewRepository) HasUserReviewedProduct(userID, productID uint64) (bool, error) {
	args := m.Called(userID, productID)
	return args.Bool(0), args.Error(1)
}
