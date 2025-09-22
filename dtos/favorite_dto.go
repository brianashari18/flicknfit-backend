package dtos

// AddFavoriteDTO is used for adding a product to a user's favorites.
type AddFavoriteDTO struct {
	UserID           uint64 `json:"user_id"`
	ProductVariantID uint64 `json:"product_variant_id"`
}
