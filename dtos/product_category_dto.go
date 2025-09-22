package dtos

// ProductCategoryCreateDTO is used for creating a new product category association.
type ProductCategoryCreateDTO struct {
	ProductID uint64 `json:"product_id" binding:"required"`
	Category  string `json:"category" binding:"required"`
}
