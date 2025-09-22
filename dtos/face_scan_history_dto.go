package dtos

// AddFaceScanDTO is used for creating a new face scan history entry.
type AddFaceScanDTO struct {
	UserID   uint64 `json:"user_id"`
	SkinTone string `json:"skin_tone"`
	ImageURL string `json:"image_url"`
}
