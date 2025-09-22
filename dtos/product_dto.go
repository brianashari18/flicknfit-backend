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
	for i, product := range products {
		result[i] = ToProductResponseDTO(product)
	}
	return result
}

// AdminProductCreateRequestDTO digunakan untuk permintaan pembuatan produk oleh admin.
type AdminProductCreateRequestDTO struct {
	BrandID     uint64  `json:"brand_id" validate:"required"`
	Name        string  `json:"name" validate:"required,min=3,max=255"`
	Description string  `json:"description" validate:"required,min=10"`
	Discount    float64 `json:"discount" validate:"required,gte=0,lte=1"`
	Rating      float64 `json:"rating" validate:"omitempty,gte=0,lte=5"`
	Reviewer    int     `json:"reviewer" validate:"omitempty,min=0"`
}

// AdminProductUpdateRequestDTO digunakan untuk permintaan pembaruan produk oleh admin.
type AdminProductUpdateRequestDTO struct {
	BrandID     uint64  `json:"brand_id" validate:"omitempty"`
	Name        string  `json:"name" validate:"omitempty,min=3,max=255"`
	Description string  `json:"description" validate:"omitempty,min=10"`
	Discount    float64 `json:"discount" validate:"omitempty,gte=0,lte=1"`
	Rating      float64 `json:"rating" validate:"omitempty,gte=0,lte=5"`
	Reviewer    int     `json:"reviewer" validate:"omitempty,min=0"`
}

// AdminProductResponseDTO digunakan untuk tanggapan API produk admin.
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

// AdminReviewCreateRequestDTO is used for creating a new review by an admin.
type AdminReviewCreateRequestDTO struct {
	Rating     int    `json:"rating" validate:"required,min=1,max=5"`
	ReviewText string `json:"review_text" validate:"required"`
}

// AdminReviewUpdateRequestDTO is used for updating an existing review by an admin.
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
