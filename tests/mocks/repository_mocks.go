package mocks

import (
	"flicknfit_backend/dtos"
	"flicknfit_backend/models"

	"github.com/stretchr/testify/mock"
)

// MockFavoriteRepository is a mock implementation of FavoriteRepository
type MockFavoriteRepository struct {
	mock.Mock
}

func (m *MockFavoriteRepository) GetUserFavorites(userID uint64) ([]models.Favorite, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Favorite), args.Error(1)
}

func (m *MockFavoriteRepository) AddFavorite(favorite *models.Favorite) error {
	args := m.Called(favorite)
	return args.Error(0)
}

func (m *MockFavoriteRepository) RemoveFavorite(userID, productItemID uint64) error {
	args := m.Called(userID, productItemID)
	return args.Error(0)
}

func (m *MockFavoriteRepository) IsFavorite(userID, productItemID uint64) (bool, error) {
	args := m.Called(userID, productItemID)
	return args.Bool(0), args.Error(1)
}

func (m *MockFavoriteRepository) GetFavoriteByID(id uint64) (*models.Favorite, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Favorite), args.Error(1)
}

// MockProductRepository is a mock implementation of ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetProductItemByID(id uint64) (*models.ProductItem, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProductItem), args.Error(1)
}

func (m *MockProductRepository) GetProductByID(id uint64) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) CreateProduct(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) UpdateProduct(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) DeleteProduct(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepository) GetAllProducts() ([]*models.Product, error) {
	args := m.Called()
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *MockProductRepository) CreateReview(review *models.Review) error {
	args := m.Called(review)
	return args.Error(0)
}

func (m *MockProductRepository) GetReviewsByProductID(productID uint64) ([]*models.Review, error) {
	args := m.Called(productID)
	return args.Get(0).([]*models.Review), args.Error(1)
}

func (m *MockProductRepository) UpdateReview(review *models.Review) error {
	args := m.Called(review)
	return args.Error(0)
}

func (m *MockProductRepository) DeleteReview(reviewID uint64) error {
	args := m.Called(reviewID)
	return args.Error(0)
}

func (m *MockProductRepository) GetReviewByID(reviewID uint64) (*models.Review, error) {
	args := m.Called(reviewID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Review), args.Error(1)
}

func (m *MockProductRepository) GetProductPublicByID(id uint64) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) GetAllProductsPublic() ([]*models.Product, error) {
	args := m.Called()
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *MockProductRepository) SearchProducts(query string) ([]*models.Product, error) {
	args := m.Called(query)
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *MockProductRepository) GetAllProductsPublicWithFilter(filter *dtos.ProductFilterRequestDTO) ([]*models.Product, error) {
	args := m.Called(filter)
	return args.Get(0).([]*models.Product), args.Error(1)
}
