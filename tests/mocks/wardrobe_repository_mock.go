package mocks

import (
	"flicknfit_backend/models"

	"github.com/stretchr/testify/mock"
)

type MockWardrobeRepository struct {
	mock.Mock
}

func (m *MockWardrobeRepository) GetUserWardrobe(userID uint64) ([]models.UserWardrobe, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.UserWardrobe), args.Error(1)
}

func (m *MockWardrobeRepository) GetWardrobeByCategory(userID uint64, category string) ([]models.UserWardrobe, error) {
	args := m.Called(userID, category)
	return args.Get(0).([]models.UserWardrobe), args.Error(1)
}

func (m *MockWardrobeRepository) CreateWardrobeItem(item *models.UserWardrobe) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockWardrobeRepository) UpdateWardrobeItem(item *models.UserWardrobe) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockWardrobeRepository) DeleteWardrobeItem(userID, itemID uint64) error {
	args := m.Called(userID, itemID)
	return args.Error(0)
}

func (m *MockWardrobeRepository) GetWardrobeItemByID(id uint64) (*models.UserWardrobe, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserWardrobe), args.Error(1)
}

func (m *MockWardrobeRepository) GetWardrobeCategories(userID uint64) ([]string, error) {
	args := m.Called(userID)
	return args.Get(0).([]string), args.Error(1)
}
