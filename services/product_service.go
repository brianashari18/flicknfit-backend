package services

import (
	"errors"
	"flicknfit_backend/dtos"
	"flicknfit_backend/models"
	"flicknfit_backend/repositories"
)

// ProductService mendefinisikan antarmuka untuk logika bisnis terkait produk.
type ProductService interface {
	// Metode untuk admin
	AdminCreateProduct(dto *dtos.AdminProductCreateRequestDTO) (*models.Product, error)
	AdminUpdateProduct(id uint64, dto *dtos.AdminProductUpdateRequestDTO) (*models.Product, error)
	AdminDeleteProduct(id uint64) error
	AdminGetProductByID(id uint64) (*models.Product, error)
	AdminGetAllProducts() ([]*models.Product, error)
	AdminCreateReview(productID uint64, dto *dtos.AdminReviewCreateRequestDTO) (*models.Review, error)
	AdminGetAllReviewsByProductID(productID uint64) ([]*models.Review, error)
	AdminUpdateReview(reviewID uint64, dto *dtos.AdminReviewUpdateRequestDTO) (*models.Review, error)
	AdminDeleteReview(reviewID uint64) error

	// Metode untuk user biasa
	GetProductPublicByID(id uint64) (*models.Product, error)
	GetAllProductsPublic() ([]*models.Product, error)
}

// productService adalah implementasi dari ProductService.
type productService struct {
	productRepository repositories.ProductRepository
}

// NewProductService membuat dan mengembalikan instance baru dari ProductService.
func NewProductService(productRepo repositories.ProductRepository) ProductService {
	return &productService{
		productRepository: productRepo,
	}
}

// AdminCreateProduct mengimplementasikan logika untuk membuat produk baru.
func (s *productService) AdminCreateProduct(dto *dtos.AdminProductCreateRequestDTO) (*models.Product, error) {
	product := &models.Product{
		Name:        dto.Name,
		Description: dto.Description,
		Discount:    dto.Discount,
		Rating:      dto.Rating,
		Reviewer:    dto.Reviewer,
	}
	if err := s.productRepository.CreateProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}

// AdminUpdateProduct mengimplementasikan logika untuk memperbarui produk.
func (s *productService) AdminUpdateProduct(id uint64, dto *dtos.AdminProductUpdateRequestDTO) (*models.Product, error) {
	product, err := s.productRepository.GetProductByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if dto.Name != "" {
		product.Name = dto.Name
	}
	if dto.Description != "" {
		product.Description = dto.Description
	}
	if dto.Discount != 0 {
		product.Discount = dto.Discount
	}
	if dto.Rating != 0 {
		product.Rating = dto.Rating
	}

	if err := s.productRepository.UpdateProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}

// AdminDeleteProduct mengimplementasikan logika untuk menghapus produk.
func (s *productService) AdminDeleteProduct(id uint64) error {
	return s.productRepository.DeleteProduct(id)
}

// AdminGetProductByID mengimplementasikan logika untuk mendapatkan produk berdasarkan ID untuk admin.
func (s *productService) AdminGetProductByID(id uint64) (*models.Product, error) {
	return s.productRepository.GetProductByID(id)
}

// AdminGetAllProducts mengimplementasikan logika untuk mendapatkan semua produk untuk admin.
func (s *productService) AdminGetAllProducts() ([]*models.Product, error) {
	return s.productRepository.GetAllProducts()
}

// GetProductPublicByID mengimplementasikan logika untuk mendapatkan produk berdasarkan ID untuk user biasa.
func (s *productService) GetProductPublicByID(id uint64) (*models.Product, error) {
	return s.productRepository.GetProductPublicByID(id)
}

// GetAllProductsPublic mengimplementasikan logika untuk mendapatkan semua produk untuk user biasa.
func (s *productService) GetAllProductsPublic() ([]*models.Product, error) {
	return s.productRepository.GetAllProductsPublic()
}

// AdminCreateReview mengimplementasikan logika untuk membuat review baru.
func (s *productService) AdminCreateReview(productID uint64, dto *dtos.AdminReviewCreateRequestDTO) (*models.Review, error) {
	review := &models.Review{
		ProductID:  productID,
		Rating:     dto.Rating,
		ReviewText: dto.ReviewText,
	}
	if err := s.productRepository.CreateReview(review); err != nil {
		return nil, err
	}
	return review, nil
}

// AdminGetAllReviewsByProductID mengimplementasikan logika untuk mendapatkan semua review berdasarkan product ID.
func (s *productService) AdminGetAllReviewsByProductID(productID uint64) ([]*models.Review, error) {
	return s.productRepository.GetReviewsByProductID(productID)
}

// AdminUpdateReview mengimplementasikan logika untuk memperbarui review.
func (s *productService) AdminUpdateReview(reviewID uint64, dto *dtos.AdminReviewUpdateRequestDTO) (*models.Review, error) {
	review, err := s.productRepository.GetReviewByID(reviewID)
	if err != nil {
		return nil, errors.New("review not found")
	}

	if dto.Rating != 0 {
		review.Rating = dto.Rating
	}
	if dto.ReviewText != "" {
		review.ReviewText = dto.ReviewText
	}

	if err := s.productRepository.UpdateReview(review); err != nil {
		return nil, err
	}
	return review, nil
}

// AdminDeleteReview mengimplementasikan logika untuk menghapus review.
func (s *productService) AdminDeleteReview(reviewID uint64) error {
	return s.productRepository.DeleteReview(reviewID)
}
