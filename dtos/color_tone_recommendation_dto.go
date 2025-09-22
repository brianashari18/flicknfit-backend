package dtos

// AddColorToneRecommendationDTO is used for adding a new color tone recommendation.
type AddColorToneRecommendationDTO struct {
	FaceScanHistoryID uint64 `json:"face_scan_history_id"`
	ColorTone         string `json:"color_tone"`
}
