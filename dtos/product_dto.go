package dtos

import (
	"flicknfit_backend/models"
	"time"
)

// ProductResponseDTO digunakan untuk menampilkan data produk publik.
type ProductResponseDTO struct {
	ID          uint64    `json:"id"`
	BrandID     uint64    `json:"brand_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Discount    float64   `json:"discount"`
	Rating      float64   `json:"rating"`
	Reviewer    int       `json:"reviewer"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToProductResponseDTO mengonversi model Product menjadi ProductResponseDTO.
func ToProductResponseDTO(product *models.Product) ProductResponseDTO {
	return ProductResponseDTO{
		ID:          product.ID,
		BrandID:     product.BrandID,
		Name:        product.Name,
		Description: product.Description,
		Discount:    product.Discount,
		Rating:      product.Rating,
		Reviewer:    product.Reviewer,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

// ToProductResponseDTOs mengonversi slice model Product menjadi slice ProductResponseDTO.
func ToProductResponseDTOs(products []*models.Product) []ProductResponseDTO {
	result := make([]ProductResponseDTO, len(products))
	for i, p := range products {
		result[i] = ToProductResponseDTO(p)
	}
	return result
}

// AdminProductCreateRequestDTO digunakan untuk membuat produk baru oleh admin.
type AdminProductCreateRequestDTO struct {
	BrandID     uint64  `json:"brand_id" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Discount    float64 `json:"discount" validate:"min=0,max=1"`
}

// AdminProductUpdateRequestDTO digunakan untuk memperbarui produk oleh admin.
type AdminProductUpdateRequestDTO struct {
	BrandID     uint64  `json:"brand_id" validate:"omitempty"`
	Name        string  `json:"name" validate:"omitempty"`
	Description string  `json:"description" validate:"omitempty"`
	Discount    float64 `json:"discount" validate:"omitempty,min=0,max=1"`
	Rating      float64 `json:"rating" validate:"omitempty,min=0,max=5"`
	Reviewer    int     `json:"reviewer" validate:"omitempty,min=0"`
}

// AdminProductResponseDTO merepresentasikan data produk yang dikembalikan ke admin.
type AdminProductResponseDTO struct {
	ID          uint64    `json:"id"`
	BrandID     uint64    `json:"brand_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Discount    float64   `json:"discount"`
	Rating      float64   `json:"rating"`
	Reviewer    int       `json:"reviewer"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToAdminProductResponseDTO mengonversi model Product menjadi AdminProductResponseDTO.
func ToAdminProductResponseDTO(product *models.Product) AdminProductResponseDTO {
	return AdminProductResponseDTO{
		ID:          product.ID,
		BrandID:     product.BrandID,
		Name:        product.Name,
		Description: product.Description,
		Discount:    product.Discount,
		Rating:      product.Rating,
		Reviewer:    product.Reviewer,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

// ToAdminProductResponseDTOs mengonversi slice model Product menjadi slice AdminProductResponseDTO.
func ToAdminProductResponseDTOs(products []*models.Product) []AdminProductResponseDTO {
	result := make([]AdminProductResponseDTO, len(products))
	for i, p := range products {
		result[i] = ToAdminProductResponseDTO(p)
	}
	return result
}

// ReviewCreateDTO digunakan untuk membuat review produk baru.
type ReviewCreateDTO struct {
	ProductID  uint64 `json:"product_id"`
	UserID     uint64 `json:"user_id"`
	Rating     int    `json:"rating" validate:"required,min=1,max=5"`
	ReviewText string `json:"review_text"`
}

// ReviewResponseDTO represents the review data returned to a public user.
type ReviewResponseDTO struct {
	ID         uint64 `json:"id"`
	ProductID  uint64 `json:"product_id"`
	Rating     int    `json:"rating"`
	ReviewText string `json:"review_text"`
}

// ToReviewResponseDTO converts a Review model to a ReviewResponseDTO.
func ToReviewResponseDTO(review *models.Review) ReviewResponseDTO {
	return ReviewResponseDTO{
		ID:         review.ID,
		ProductID:  review.ProductID,
		Rating:     review.Rating,
		ReviewText: review.ReviewText,
	}
}

// AdminReviewCreateRequestDTO digunakan untuk membuat review oleh admin.
type AdminReviewCreateRequestDTO struct {
	Rating     int    `json:"rating" validate:"required,min=1,max=5"`
	ReviewText string `json:"review_text"`
}

// AdminReviewUpdateRequestDTO digunakan untuk memperbarui review oleh admin.
type AdminReviewUpdateRequestDTO struct {
	Rating     int    `json:"rating" validate:"omitempty,min=1,max=5"`
	ReviewText string `json:"review_text" validate:"omitempty"`
}

// AdminReviewResponseDTO represents the review data returned to an admin.
type AdminReviewResponseDTO struct {
	ID         uint64 `json:"id"`
	ProductID  uint64 `json:"product_id"`
	Rating     int    `json:"rating"`
	ReviewText string `json:"review_text"`
}

// ProductItemDTO represents a product item.
type ProductItemDTO struct {
	ID        uint64 `json:"id"`
	ProductID uint64 `json:"product_id"`
	SKU       string `json:"sku"`
	Price     int    `json:"price"`
	Stock     int    `json:"stock"`
	PhotoURL  string `json:"photo_url"`
}

func ToAdminReviewResponseDTO(review models.Review) AdminReviewResponseDTO {
	return AdminReviewResponseDTO{
		ID:         review.ID,
		ProductID:  review.ProductID,
		Rating:     review.Rating,
		ReviewText: review.ReviewText,
	}
}

func ToAdminReviewResponseDTOs(reviews []*models.Review) []AdminReviewResponseDTO {
	result := make([]AdminReviewResponseDTO, 0, len(reviews))
	for _, r := range reviews {
		result = append(result, ToAdminReviewResponseDTO(*r))
	}
	return result
}

func ToProductItemDTO(item models.ProductItem) ProductItemDTO {
	return ProductItemDTO{
		ID:        item.ID,
		ProductID: item.ProductID,
		SKU:       item.SKU,
		Price:     item.Price,
		Stock:     item.Stock,
		PhotoURL:  item.PhotoURL,
	}
}
