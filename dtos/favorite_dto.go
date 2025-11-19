package dtos

import (
	"flicknfit_backend/models"
	"time"
)

// AddFavoriteDTO is used for adding a product item to a user's favorites.
type AddFavoriteDTO struct {
	ProductItemID uint64 `json:"product_item_id" validate:"required"`
}

// FavoriteResponseDTO represents a favorite in response
type FavoriteResponseDTO struct {
	ID            uint64                 `json:"id"`
	ProductItemID uint64                 `json:"product_item_id"`
	ProductItem   ProductItemResponseDTO `json:"product_item"`
	CreatedAt     time.Time              `json:"created_at"`
}

// FavoriteToggleResponseDTO represents toggle favorite response
type FavoriteToggleResponseDTO struct {
	ProductItemID uint64 `json:"product_item_id"`
	IsFavorited   bool   `json:"is_favorited"`
	Message       string `json:"message"`
}

// ProductItemResponseDTO represents product item in response
type ProductItemResponseDTO struct {
	ID       uint64                  `json:"id"`
	SKU      string                  `json:"sku"`
	Price    int                     `json:"price"`
	Stock    int                     `json:"stock"`
	PhotoURL string                  `json:"photo_url"`
	Product  ProductBasicResponseDTO `json:"product"`
}

// ProductBasicResponseDTO represents basic product info
type ProductBasicResponseDTO struct {
	ID          uint64  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	BrandID     uint64  `json:"brand_id"`
	Rating      float64 `json:"rating"`
	Reviewer    int     `json:"reviewer"`
}

// ToFavoriteResponseDTO converts a favorite model to response DTO
func ToFavoriteResponseDTO(favorite models.Favorite) FavoriteResponseDTO {
	return FavoriteResponseDTO{
		ID:            favorite.ID,
		ProductItemID: favorite.ProductItemID,
		ProductItem:   ToProductItemResponseDTO(favorite.ProductItem),
		CreatedAt:     favorite.CreatedAt,
	}
}

// ToFavoriteResponseDTOs converts slice of favorite models to response DTOs
func ToFavoriteResponseDTOs(favorites []models.Favorite) []FavoriteResponseDTO {
	var result []FavoriteResponseDTO
	for _, favorite := range favorites {
		result = append(result, ToFavoriteResponseDTO(favorite))
	}
	return result
}

// ToProductItemResponseDTO converts product item model to response DTO
func ToProductItemResponseDTO(item models.ProductItem) ProductItemResponseDTO {
	return ProductItemResponseDTO{
		ID:       item.ID,
		SKU:      item.SKU,
		Price:    item.Price,
		Stock:    item.Stock,
		PhotoURL: item.PhotoURL,
		Product:  ToProductBasicResponseDTO(item.Product),
	}
}

// ToProductBasicResponseDTO converts product model to basic response DTO
func ToProductBasicResponseDTO(product models.Product) ProductBasicResponseDTO {
	return ProductBasicResponseDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		BrandID:     product.BrandID,
		Rating:      product.Rating,
		Reviewer:    product.Reviewer,
	}
}
