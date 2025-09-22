package dtos

// AddBodyScanDTO is used for creating a new body scan history entry.
type AddBodyScanDTO struct {
	UserID    uint64 `json:"user_id"`
	BodyShape string `json:"body_shape"`
	ImageURL  string `json:"image_url"`
}
