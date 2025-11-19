package dtos

import (
	"flicknfit_backend/models"
	"time"
)

// CreateReviewDTO for creating a new review
type CreateReviewDTO struct {
	ProductID  uint64 `json:"product_id" validate:"required"`
	Rating     int    `json:"rating" validate:"required,min=1,max=5"`
	ReviewText string `json:"review_text" validate:"required,min=10,max=1000"`
}

// UpdateReviewDTO for updating an existing review
type UpdateReviewDTO struct {
	Rating     int    `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	ReviewText string `json:"review_text,omitempty" validate:"omitempty,min=10,max=1000"`
}

// ReviewResponseDTO represents a review in response
type ReviewResponseDTO struct {
	ID         uint64               `json:"id"`
	UserID     uint64               `json:"user_id"`
	ProductID  uint64               `json:"product_id"`
	Rating     int                  `json:"rating"`
	ReviewText string               `json:"review_text"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
	User       UserBasicResponseDTO `json:"user"`
}

// ReviewListResponseDTO represents paginated review list
type ReviewListResponseDTO struct {
	Reviews    []ReviewResponseDTO    `json:"reviews"`
	Statistics map[string]interface{} `json:"statistics"`
	Pagination PaginationDTO          `json:"pagination"`
}

// ReviewStatsResponseDTO represents review statistics
type ReviewStatsResponseDTO struct {
	AverageRating float64        `json:"average_rating"`
	TotalReviews  int            `json:"total_reviews"`
	RatingCounts  map[string]int `json:"rating_counts"`
}

// UserBasicResponseDTO represents basic user info
type UserBasicResponseDTO struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
}

// PaginationDTO represents pagination info
type PaginationDTO struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

// ToReviewResponseDTO converts a review model to response DTO
func ToReviewResponseDTO(review *models.Review) ReviewResponseDTO {
	return ReviewResponseDTO{
		ID:         review.ID,
		UserID:     review.UserID,
		ProductID:  review.ProductID,
		Rating:     review.Rating,
		ReviewText: review.ReviewText,
		CreatedAt:  review.CreatedAt,
		UpdatedAt:  review.UpdatedAt,
		User:       ToUserBasicResponseDTO(review.User),
	}
}

// ToReviewResponseDTOs converts slice of review models to response DTOs
func ToReviewResponseDTOs(reviews []models.Review) []ReviewResponseDTO {
	var result []ReviewResponseDTO
	for _, review := range reviews {
		result = append(result, ToReviewResponseDTO(&review))
	}
	return result
}

// ToUserBasicResponseDTO converts user model to basic response DTO
func ToUserBasicResponseDTO(user models.User) UserBasicResponseDTO {
	return UserBasicResponseDTO{
		ID:       user.ID,
		Username: user.Username,
		FullName: user.Username, // Use username as full name since FullName field doesn't exist
	}
}
