package repositories

import (
	"flicknfit_backend/models"

	"gorm.io/gorm"
)

// BrandRepository defines the interface for data access operations on the Brand model.
//
//go:generate mockery --name BrandRepository
type BrandRepository interface {
	CreateBrand(brand *models.Brand) error
	GetAllBrands() ([]models.Brand, error)
	GetBrandByID(id uint64) (*models.Brand, error)
	UpdateBrand(brand *models.Brand) error
	DeleteBrand(brand *models.Brand) error
}

// brandRepository is the implementation of BrandRepository.
type brandRepository struct {
	BaseRepository
}

// NewBrandRepository creates and returns a new instance of BrandRepository.
func NewBrandRepository(db *gorm.DB) BrandRepository {
	return &brandRepository{BaseRepository{DB: db}}
}

// CreateBrand creates a new brand record in the database.
func (r *brandRepository) CreateBrand(brand *models.Brand) error {
	return r.DB.Create(brand).Error
}

// GetAllBrands retrieves all brand records from the database.
func (r *brandRepository) GetAllBrands() ([]models.Brand, error) {
	var brands []models.Brand
	if err := r.DB.Find(&brands).Error; err != nil {
		return nil, err
	}
	return brands, nil
}

// GetBrandByID retrieves a brand record by its ID.
func (r *brandRepository) GetBrandByID(id uint64) (*models.Brand, error) {
	var brand models.Brand
	if err := r.DB.First(&brand, id).Error; err != nil {
		return nil, err
	}
	return &brand, nil
}

// UpdateBrand updates an existing brand record in the database.
func (r *brandRepository) UpdateBrand(brand *models.Brand) error {
	return r.DB.Save(brand).Error
}

// DeleteBrand deletes a brand record from the database.
func (r *brandRepository) DeleteBrand(brand *models.Brand) error {
	return r.DB.Delete(brand).Error
}
