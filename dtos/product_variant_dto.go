package dtos

// ProductVariantCreateDTO is used for creating a new product variant.
type ProductVariantCreateDTO struct {
	ProductID uint64 `json:"product_id" binding:"required"`
	SKU       string `json:"sku" binding:"required"`
	Size      string `json:"size" binding:"required"`
	Color     string `json:"color" binding:"required"`
	Stock     int    `json:"stock" binding:"required,min=0"`
	Price     int    `json:"price" binding:"required,min=0"`
}

// ProductVariantUpdateDTO is used for updating an existing product variant.
type ProductVariantUpdateDTO struct {
	SKU   string `json:"sku"`
	Size  string `json:"size"`
	Color string `json:"color"`
	Stock int    `json:"stock"`
	Price int    `json:"price"`
}
