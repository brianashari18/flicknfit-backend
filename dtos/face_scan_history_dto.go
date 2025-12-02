package dtos

import "time"

// AddFaceScanDTO is used for creating a new face scan history entry.
type AddFaceScanDTO struct {
	UserID   uint64 `json:"user_id"`
	SkinTone string `json:"skin_tone"`
	ImageURL string `json:"image_url"`
}

// FaceScanHistoryDTO represents face scan history response
type FaceScanHistoryDTO struct {
	ID                   uint64    `json:"id"`
	ScanName             string    `json:"scan_name"`
	SkinTone             string    `json:"skin_tone"`
	ImageURL             string    `json:"image_url"` // Signed URL from Supabase
	ColorRecommendations []string  `json:"color_recommendations"`
	Confidence           *float64  `json:"confidence,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

// FaceScanHistoryListDTO represents list of face scan history
type FaceScanHistoryListDTO struct {
	Histories []FaceScanHistoryDTO `json:"histories"`
	Total     int64                `json:"total"`
}
