package dtos

// AddStyleRecommendationDTO is used for adding a new style recommendation.
type AddStyleRecommendationDTO struct {
	BodyScanHistoriesID uint64 `json:"body_scan_history_id"`
	Style               string `json:"style"`
}
