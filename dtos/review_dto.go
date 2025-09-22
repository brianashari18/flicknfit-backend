package dtos

// ReviewCreateDTO is used for creating a new product review.
type ReviewCreateDTO struct {
	ProductID  uint64 `json:"product_id" binding:"required"`
	UserID     uint64 `json:"user_id" binding:"required"`
	Rating     int    `json:"rating" binding:"required,min=1,max=5"`
	ReviewText string `json:"review_text"`
}
