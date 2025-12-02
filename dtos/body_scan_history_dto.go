package dtos

import "time"

// AddBodyScanDTO is used for creating a new body scan history entry.
type AddBodyScanDTO struct {
	UserID    uint64 `json:"user_id"`
	BodyShape string `json:"body_shape"`
	ImageURL  string `json:"image_url"`
}

// BodyScanHistoryDTO represents body scan history response
type BodyScanHistoryDTO struct {
	ID                   uint64    `json:"id"`
	ScanName             string    `json:"scan_name"`
	BodyType             string    `json:"body_type"`
	Gender               string    `json:"gender"`
	ImageURL             string    `json:"image_url"` // Signed URL from Supabase
	StyleRecommendations []string  `json:"style_recommendations"`
	Confidence           float64   `json:"confidence"`
	CreatedAt            time.Time `json:"created_at"`
}

// BodyScanHistoryListDTO represents list of body scan history
type BodyScanHistoryListDTO struct {
	Histories []BodyScanHistoryDTO `json:"histories"`
	Total     int64                `json:"total"`
}
