package dtos

import (
	"flicknfit_backend/models"
	"time"
)

// CreateWardrobeItemDTO for creating a new wardrobe item
type CreateWardrobeItemDTO struct {
	Category string `json:"category" validate:"required,min=2,max=50"`
	ImageURL string `json:"image_url" validate:"required,url"`
}

// UpdateWardrobeItemDTO for updating a wardrobe item
type UpdateWardrobeItemDTO struct {
	Category string `json:"category,omitempty" validate:"omitempty,min=2,max=50"`
	ImageURL string `json:"image_url,omitempty" validate:"omitempty,url"`
}

// WardrobeItemResponseDTO represents a wardrobe item in response
type WardrobeItemResponseDTO struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Category  string    `json:"category"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WardrobeResponseDTO represents organized wardrobe response
type WardrobeResponseDTO struct {
	Categories map[string][]WardrobeItemResponseDTO `json:"categories"`
	Summary    WardrobeSummaryDTO                   `json:"summary"`
}

// WardrobeSummaryDTO represents wardrobe summary
type WardrobeSummaryDTO struct {
	TotalItems      int      `json:"total_items"`
	CategoriesCount int      `json:"categories_count"`
	Categories      []string `json:"categories"`
}

// ToWardrobeItemResponseDTO converts a wardrobe model to response DTO
func ToWardrobeItemResponseDTO(item models.UserWardrobe) WardrobeItemResponseDTO {
	return WardrobeItemResponseDTO{
		ID:        item.ID,
		UserID:    item.UserID,
		Category:  item.Category,
		ImageURL:  item.ImageURL,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

// ToWardrobeItemResponseDTOs converts slice of wardrobe models to response DTOs
func ToWardrobeItemResponseDTOs(items []models.UserWardrobe) []WardrobeItemResponseDTO {
	var result []WardrobeItemResponseDTO
	for _, item := range items {
		result = append(result, ToWardrobeItemResponseDTO(item))
	}
	return result
}

// ToWardrobeResponseDTO converts organized wardrobe to response DTO
func ToWardrobeResponseDTO(wardrobe map[string][]models.UserWardrobe) WardrobeResponseDTO {
	categories := make(map[string][]WardrobeItemResponseDTO)
	totalItems := 0
	categoryNames := make([]string, 0, len(wardrobe))

	for category, items := range wardrobe {
		categories[category] = ToWardrobeItemResponseDTOs(items)
		totalItems += len(items)
		categoryNames = append(categoryNames, category)
	}

	return WardrobeResponseDTO{
		Categories: categories,
		Summary: WardrobeSummaryDTO{
			TotalItems:      totalItems,
			CategoriesCount: len(wardrobe),
			Categories:      categoryNames,
		},
	}
}

// Legacy DTO for backward compatibility
type AddWardrobeDTO struct {
	UserID   uint64 `json:"user_id"`
	Category string `json:"category"`
	ImageURL string `json:"image_url"`
}
