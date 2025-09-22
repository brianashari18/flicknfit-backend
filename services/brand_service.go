package services

import (
	"errors"
	"flicknfit_backend/dtos"
	"flicknfit_backend/models"
	"flicknfit_backend/repositories"
)

// BrandService defines the interface for business logic related to brands.
type BrandService interface {
	// Admin
	AdminCreateBrand(dto *dtos.BrandCreateRequestDTO) (*models.Brand, error)
	AdminUpdateBrand(id uint64, dto *dtos.BrandUpdateRequestDTO) (*models.Brand, error)
	AdminDeleteBrand(id uint64) error

	// Public
	GetAllBrands() ([]dtos.BrandResponseDTO, error)
	GetBrandByID(id uint64) (*dtos.BrandResponseDTO, error)
}

// brandService is the implementation of BrandService.
type brandService struct {
	brandRepository repositories.BrandRepository
}

// NewBrandService creates and returns a new instance of BrandService.
func NewBrandService(brandRepository repositories.BrandRepository) BrandService {
	return &brandService{
		brandRepository: brandRepository,
	}
}

// AdminCreateBrand handles the creation of a new brand by an admin.
func (s *brandService) AdminCreateBrand(dto *dtos.BrandCreateRequestDTO) (*models.Brand, error) {
	brand := models.Brand{
		Name:        dto.Name,
		Description: dto.Description,
		LogoURL:     dto.LogoURL,
		WebsiteURL:  dto.WebsiteURL,
	}

	if err := s.brandRepository.CreateBrand(&brand); err != nil {
		return nil, errors.New("failed to create brand")
	}

	return &brand, nil
}

// AdminUpdateBrand handles the update of an existing brand by an admin.
func (s *brandService) AdminUpdateBrand(id uint64, dto *dtos.BrandUpdateRequestDTO) (*models.Brand, error) {
	brand, err := s.brandRepository.GetBrandByID(id)
	if err != nil {
		return nil, errors.New("brand not found")
	}

	if dto.Name != "" {
		brand.Name = dto.Name
	}
	if dto.Description != "" {
		brand.Description = dto.Description
	}
	if dto.LogoURL != "" {
		brand.LogoURL = dto.LogoURL
	}
	if dto.WebsiteURL != "" {
		brand.WebsiteURL = dto.WebsiteURL
	}

	if err := s.brandRepository.UpdateBrand(brand); err != nil {
		return nil, errors.New("failed to update brand")
	}

	return brand, nil
}

// AdminDeleteBrand handles the deletion of a brand by an admin.
func (s *brandService) AdminDeleteBrand(id uint64) error {
	brand, err := s.brandRepository.GetBrandByID(id)
	if err != nil {
		return errors.New("brand not found")
	}

	if err := s.brandRepository.DeleteBrand(brand); err != nil {
		return errors.New("failed to delete brand")
	}

	return nil
}

// GetAllBrands retrieves all brands for public display.
func (s *brandService) GetAllBrands() ([]dtos.BrandResponseDTO, error) {
	brands, err := s.brandRepository.GetAllBrands()
	if err != nil {
		return nil, errors.New("failed to retrieve brands")
	}

	return dtos.ToBrandResponseDTOs(brands), nil
}

// GetBrandByID retrieves a single brand by its ID for public display.
func (s *brandService) GetBrandByID(id uint64) (*dtos.BrandResponseDTO, error) {
	brand, err := s.brandRepository.GetBrandByID(id)
	if err != nil {
		return nil, errors.New("brand not found")
	}

	response := dtos.ToBrandResponseDTO(*brand)
	return &response, nil
}
