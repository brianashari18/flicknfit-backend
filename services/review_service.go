package services

import (
	"flicknfit_backend/constants"
	"flicknfit_backend/dtos"
	"flicknfit_backend/errors"
	"flicknfit_backend/models"
	"flicknfit_backend/repositories"
)

// ReviewService defines business logic for reviews
type ReviewService interface {
	GetProductReviews(productID uint64, page, limit int) (*dtos.ReviewListResponseDTO, error)
	CreateReview(userID uint64, dto *dtos.CreateReviewDTO) error
	UpdateReview(userID, reviewID uint64, dto *dtos.UpdateReviewDTO) error
	DeleteReview(userID, reviewID uint64) error
	GetUserReviews(userID uint64) ([]models.Review, error)
	GetProductReviewStats(productID uint64) (map[string]interface{}, error)
}

// reviewService implements ReviewService interface
type reviewService struct {
	reviewRepo  repositories.ReviewRepository
	productRepo repositories.ProductRepository
}

// NewReviewService creates a new review service
func NewReviewService(reviewRepo repositories.ReviewRepository, productRepo repositories.ProductRepository) ReviewService {
	return &reviewService{
		reviewRepo:  reviewRepo,
		productRepo: productRepo,
	}
}

// GetProductReviews retrieves reviews for a product with pagination
func (s *reviewService) GetProductReviews(productID uint64, page, limit int) (*dtos.ReviewListResponseDTO, error) {
	// Validate product exists
	_, err := s.productRepo.GetProductByID(productID)
	if err != nil {
		return nil, errors.NewNotFoundError("Product")
	}

	// Set default pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > constants.MaxPageSize {
		limit = constants.DefaultPageSize
	}

	offset := (page - 1) * limit

	// Get reviews
	reviews, err := s.reviewRepo.GetProductReviews(productID, limit, offset)
	if err != nil {
		return nil, errors.NewDatabaseError("get product reviews", err)
	}

	// Get review stats
	stats, err := s.reviewRepo.GetProductReviewStats(productID)
	if err != nil {
		return nil, errors.NewDatabaseError("get review stats", err)
	}

	response := &dtos.ReviewListResponseDTO{
		Reviews:    dtos.ToReviewResponseDTOs(reviews),
		Statistics: stats,
		Pagination: dtos.PaginationDTO{
			Page:  page,
			Limit: limit,
			Total: len(reviews),
		},
	}

	return response, nil
}

// CreateReview creates a new review
func (s *reviewService) CreateReview(userID uint64, dto *dtos.CreateReviewDTO) error {
	// Validate product exists
	_, err := s.productRepo.GetProductByID(dto.ProductID)
	if err != nil {
		return errors.NewNotFoundError("Product")
	}

	// Check if user has already reviewed this product
	hasReviewed, err := s.reviewRepo.HasUserReviewedProduct(userID, dto.ProductID)
	if err != nil {
		return errors.NewDatabaseError("check existing review", err)
	}
	if hasReviewed {
		return errors.NewConflictError("Review already exists for this product")
	}

	review := &models.Review{
		UserID:     userID,
		ProductID:  dto.ProductID,
		Rating:     dto.Rating,
		ReviewText: dto.ReviewText,
	}

	err = s.reviewRepo.CreateReview(review)
	if err != nil {
		return errors.NewDatabaseError("create review", err)
	}

	return nil
}

// UpdateReview updates an existing review
func (s *reviewService) UpdateReview(userID, reviewID uint64, dto *dtos.UpdateReviewDTO) error {
	// Get existing review
	review, err := s.reviewRepo.GetReviewByID(reviewID)
	if err != nil {
		return errors.NewNotFoundError("Review")
	}

	// Check if user owns the review
	if review.UserID != userID {
		return errors.NewAuthorizationError("Cannot update another user's review")
	}

	// Update fields
	if dto.Rating != 0 {
		review.Rating = dto.Rating
	}
	if dto.ReviewText != "" {
		review.ReviewText = dto.ReviewText
	}

	err = s.reviewRepo.UpdateReview(review)
	if err != nil {
		return errors.NewDatabaseError("update review", err)
	}

	return nil
}

// DeleteReview deletes a review
func (s *reviewService) DeleteReview(userID, reviewID uint64) error {
	// Get existing review
	review, err := s.reviewRepo.GetReviewByID(reviewID)
	if err != nil {
		return errors.NewNotFoundError("Review")
	}

	// Check if user owns the review
	if review.UserID != userID {
		return errors.NewAuthorizationError("Cannot delete another user's review")
	}

	err = s.reviewRepo.DeleteReview(reviewID)
	if err != nil {
		return errors.NewDatabaseError("delete review", err)
	}

	return nil
}

// GetUserReviews retrieves all reviews by a user
func (s *reviewService) GetUserReviews(userID uint64) ([]models.Review, error) {
	reviews, err := s.reviewRepo.GetUserReviews(userID)
	if err != nil {
		return nil, errors.NewDatabaseError("get user reviews", err)
	}
	return reviews, nil
}

// GetProductReviewStats gets review statistics for a product
func (s *reviewService) GetProductReviewStats(productID uint64) (map[string]interface{}, error) {
	// Validate product exists
	_, err := s.productRepo.GetProductByID(productID)
	if err != nil {
		return nil, errors.NewNotFoundError("Product")
	}

	stats, err := s.reviewRepo.GetProductReviewStats(productID)
	if err != nil {
		return nil, errors.NewDatabaseError("get review stats", err)
	}

	return stats, nil
}
