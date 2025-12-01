package dtos

// AI Prediction DTOs

// SkinColorTonePredictionResponseDTO represents the response from skin color tone prediction
// API Response: {"skin_tone": "Light Spring"}
type SkinColorTonePredictionResponseDTO struct {
	SkinTone             string   `json:"skin_tone"`
	ColorRecommendations []string `json:"color_recommendations,omitempty"`
}

// WomanBodyScanPredictionResponseDTO represents the response from woman body scan prediction
// API Response: {"predicted_class": "hourglass", "confidence": 0.4587}
type WomanBodyScanPredictionResponseDTO struct {
	PredictedClass       string   `json:"predicted_class"`
	Confidence           float64  `json:"confidence"`
	StyleRecommendations []string `json:"style_recommendations,omitempty"`
}

// MenBodyScanPredictionResponseDTO represents the response from men's body scan prediction
// API Response: {"predicted_class": "endomorf", "confidence": 0.9764}
type MenBodyScanPredictionResponseDTO struct {
	PredictedClass       string   `json:"predicted_class"`
	Confidence           float64  `json:"confidence"`
	StyleRecommendations []string `json:"style_recommendations,omitempty"`
}

// AIErrorResponseDTO represents error response from AI API
type AIErrorResponseDTO struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}
