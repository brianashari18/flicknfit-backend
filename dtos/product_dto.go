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

// ProductPublicResponseDTO mewakili data produk lengkap untuk pengguna publik, termasuk item produk dan ulasan.
type ProductPublicResponseDTO struct {
	ID           uint64           `json:"id"`
	BrandID      uint64           `json:"brand_id"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Discount     float64          `json:"discount"`
	Rating       float64          `json:"rating"`
	Reviewer     int              `json:"reviewer"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	ProductItems []ProductItemDTO `json:"product_items"`
	Reviews      []ReviewDTO      `json:"reviews"`
}

// AdminProductDetailsDTO mewakili data produk lengkap untuk admin, termasuk item produk, kategori, dan ulasan.
type AdminProductDetailsDTO struct {
	ID                uint64               `json:"id"`
	BrandID           uint64               `json:"brand_id"`
	Name              string               `json:"name"`
	Description       string               `json:"description"`
	Discount          float64              `json:"discount"`
	Rating            float64              `json:"rating"`
	Reviewer          int                  `json:"reviewer"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
	ProductItems      []ProductItemDTO     `json:"product_items"`
	ProductCategories []ProductCategoryDTO `json:"product_categories"`
	Reviews           []ReviewDTO          `json:"reviews"`
}

// ToProductPublicResponseDTO mengonversi model Product menjadi ProductPublicResponseDTO.
func ToProductPublicResponseDTO(product *models.Product) ProductPublicResponseDTO {
	productItems := make([]ProductItemDTO, len(product.ProductItems))
	for i, item := range product.ProductItems {
		productItems[i] = ToProductItemDTO(item)
	}

	reviews := make([]ReviewDTO, len(product.Reviews))
	for i, review := range product.Reviews {
		reviews[i] = ToReviewDTO(&review)
	}

	return ProductPublicResponseDTO{
		ID:           product.ID,
		BrandID:      product.BrandID,
		Name:         product.Name,
		Description:  product.Description,
		Discount:     product.Discount,
		Rating:       product.Rating,
		Reviewer:     product.Reviewer,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
		ProductItems: productItems,
		Reviews:      reviews,
	}
}

// ToAdminProductDetailsDTO mengonversi model Product menjadi AdminProductDetailsDTO.
func ToAdminProductDetailsDTO(product *models.Product) AdminProductDetailsDTO {
	productItems := make([]ProductItemDTO, len(product.ProductItems))
	for i, item := range product.ProductItems {
		productItems[i] = ToProductItemDTO(item)
	}
	productCategories := make([]ProductCategoryDTO, len(product.ProductCategories))
	for i, category := range product.ProductCategories {
		productCategories[i] = ToProductCategoryDTO(&category)
	}
	reviews := make([]ReviewDTO, len(product.Reviews))
	for i, review := range product.Reviews {
		reviews[i] = ToReviewDTO(&review)
	}

	return AdminProductDetailsDTO{
		ID:                product.ID,
		BrandID:           product.BrandID,
		Name:              product.Name,
		Description:       product.Description,
		Discount:          product.Discount,
		Rating:            product.Rating,
		Reviewer:          product.Reviewer,
		CreatedAt:         product.CreatedAt,
		UpdatedAt:         product.UpdatedAt,
		ProductItems:      productItems,
		ProductCategories: productCategories,
		Reviews:           reviews,
	}
}

// ReviewCreateDTO digunakan untuk membuat review produk baru.
type ReviewCreateDTO struct {
	ProductID  uint64 `json:"product_id"`
	UserID     uint64 `json:"user_id"`
	Rating     int    `json:"rating" validate:"required,min=1,max=5"`
	ReviewText string `json:"review_text"`
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

// ProductFilterRequestDTO is used for filtering products.
type ProductFilterRequestDTO struct {
	Name      string  `query:"name"`
	MinPrice  float64 `query:"min_price"`
	MaxPrice  float64 `query:"max_price"`
	BrandName string  `query:"brand_name"`
	Category  string  `query:"category"`
	MinRating float64 `query:"min_rating"`
}

// ReviewDTO represents a product review.
type ReviewDTO struct {
	ID         uint64 `json:"id"`
	ProductID  uint64 `json:"product_id"`
	Rating     int    `json:"rating"`
	ReviewText string `json:"review_text"`
}

// ProductCategoryDTO represents a product category.
type ProductCategoryDTO struct {
	ID        uint64 `json:"id"`
	ProductID uint64 `json:"product_id"`
	Category  string `json:"category"`
}

// ToReviewDTO mengonversi model Review menjadi ReviewDTO.
func ToReviewDTO(review *models.Review) ReviewDTO {
	return ReviewDTO{
		ID:         review.ID,
		ProductID:  review.ProductID,
		Rating:     review.Rating,
		ReviewText: review.ReviewText,
	}
}

// ToProductCategoryDTO mengonversi model ProductCategory menjadi ProductCategoryDTO.
func ToProductCategoryDTO(category *models.ProductCategory) ProductCategoryDTO {
	return ProductCategoryDTO{
		ID:        category.ID,
		ProductID: category.ProductID,
		Category:  category.Category,
	}
}
