package dtos

// AddWardrobeDTO is used for adding an item to the user's wardrobe.
type AddWardrobeDTO struct {
	UserID   uint64 `json:"user_id"`
	Category string `json:"category"`
	ImageURL string `json:"image_url"`
}
