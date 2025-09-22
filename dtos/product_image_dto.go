package dtos

// ProductImageCreateDTO is used for adding a new image to a product.
type ProductImageCreateDTO struct {
	ProductID uint64 `json:"product_id" binding:"required"`
	ImageURL  string `json:"image_url" binding:"required"`
}
