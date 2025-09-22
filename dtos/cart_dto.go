package dtos

// AddToCartDTO is used for adding a product to the cart.
type AddToCartDTO struct {
	UserID           uint64 `json:"user_id"`
	ProductVariantID uint64 `json:"product_variant_id"`
	Quantity         int    `json:"quantity"`
}
