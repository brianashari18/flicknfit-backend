package dtos

import (
	"flicknfit_backend/models"
	"time"
)

// BrandCreateRequestDTO is used for creating a new brand by an admin.
type BrandCreateRequestDTO struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"required,min=1,max=255"`
	LogoURL     string `json:"logo_url" validate:"omitempty,url"`
	WebsiteURL  string `json:"website_url" validate:"omitempty,url"`
}

// BrandUpdateRequestDTO is used for updating an existing brand by an admin.
type BrandUpdateRequestDTO struct {
	Name        string `json:"name" validate:"omitempty,min=1,max=100"`
	Description string `json:"description" validate:"omitempty,min=1,max=255"`
	LogoURL     string `json:"logo_url" validate:"omitempty,url"`
	WebsiteURL  string `json:"website_url" validate:"omitempty,url"`
}

// BrandResponseDTO represents the brand data returned to any user.
type BrandResponseDTO struct {
	ID            uint64    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Rating        float64   `json:"rating"`
	Reviewer      uint      `json:"reviewer"`
	LogoURL       string    `json:"logo_url"`
	WebsiteURL    string    `json:"website_url"`
	TotalProducts uint      `json:"total_products"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ToBrandResponseDTO converts a models.Brand to a BrandResponseDTO.
func ToBrandResponseDTO(brand models.Brand) BrandResponseDTO {
	return BrandResponseDTO{
		ID:            brand.ID,
		Name:          brand.Name,
		Description:   brand.Description,
		Rating:        brand.Rating,
		Reviewer:      brand.Reviewer,
		LogoURL:       brand.LogoURL,
		WebsiteURL:    brand.WebsiteURL,
		TotalProducts: brand.TotalProducts,
		CreatedAt:     brand.CreatedAt,
		UpdatedAt:     brand.UpdatedAt,
	}
}

// ToBrandResponseDTOs converts a slice of models.Brand to a slice of BrandResponseDTO.
func ToBrandResponseDTOs(brands []models.Brand) []BrandResponseDTO {
	result := make([]BrandResponseDTO, 0, len(brands))
	for _, brand := range brands {
		result = append(result, ToBrandResponseDTO(brand))
	}
	return result
}
