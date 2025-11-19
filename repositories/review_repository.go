package repositories

import (
	"flicknfit_backend/models"

	"gorm.io/gorm"
)

// ReviewRepository defines methods for review data access
type ReviewRepository interface {
	GetProductReviews(productID uint64, limit, offset int) ([]models.Review, error)
	GetReviewByID(id uint64) (*models.Review, error)
	CreateReview(review *models.Review) error
	UpdateReview(review *models.Review) error
	DeleteReview(id uint64) error
	GetUserReviews(userID uint64) ([]models.Review, error)
	GetProductReviewStats(productID uint64) (map[string]interface{}, error)
	HasUserReviewedProduct(userID, productID uint64) (bool, error)
}

// reviewRepository implements ReviewRepository interface
type reviewRepository struct {
	db *gorm.DB
}

// NewReviewRepository creates a new review repository
func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

// GetProductReviews retrieves reviews for a specific product
func (r *reviewRepository) GetProductReviews(productID uint64, limit, offset int) ([]models.Review, error) {
	var reviews []models.Review
	query := r.db.Preload("User").Where("product_id = ?", productID)

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

// GetReviewByID retrieves a review by ID
func (r *reviewRepository) GetReviewByID(id uint64) (*models.Review, error) {
	var review models.Review
	err := r.db.Preload("User").Preload("Product").First(&review, id).Error
	return &review, err
}

// CreateReview creates a new review
func (r *reviewRepository) CreateReview(review *models.Review) error {
	return r.db.Create(review).Error
}

// UpdateReview updates an existing review
func (r *reviewRepository) UpdateReview(review *models.Review) error {
	return r.db.Save(review).Error
}

// DeleteReview deletes a review
func (r *reviewRepository) DeleteReview(id uint64) error {
	return r.db.Delete(&models.Review{}, id).Error
}

// GetUserReviews retrieves all reviews by a user
func (r *reviewRepository) GetUserReviews(userID uint64) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Preload("Product").Where("user_id = ?", userID).
		Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

// GetProductReviewStats calculates review statistics for a product
func (r *reviewRepository) GetProductReviewStats(productID uint64) (map[string]interface{}, error) {
	var result struct {
		TotalReviews  int     `json:"total_reviews"`
		AverageRating float64 `json:"average_rating"`
	}

	err := r.db.Model(&models.Review{}).
		Select("COUNT(*) as total_reviews, AVG(rating) as average_rating").
		Where("product_id = ?", productID).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	// Get rating distribution
	var ratingDist []struct {
		Rating int `json:"rating"`
		Count  int `json:"count"`
	}

	err = r.db.Model(&models.Review{}).
		Select("rating, COUNT(*) as count").
		Where("product_id = ?", productID).
		Group("rating").
		Order("rating DESC").
		Scan(&ratingDist).Error

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_reviews":       result.TotalReviews,
		"average_rating":      result.AverageRating,
		"rating_distribution": ratingDist,
	}, nil
}

// HasUserReviewedProduct checks if user has already reviewed a product
func (r *reviewRepository) HasUserReviewedProduct(userID, productID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&models.Review{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error
	return count > 0, err
}
