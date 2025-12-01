package services

import (
	"bytes"
	"encoding/json"
	"flicknfit_backend/config"
	"flicknfit_backend/dtos"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

// AIService defines the interface for AI prediction operations
type AIService interface {
	PredictSkinColorTone(file multipart.File, filename string) (*dtos.SkinColorTonePredictionResponseDTO, error)
	PredictWomanBodyScan(file multipart.File, filename string) (*dtos.WomanBodyScanPredictionResponseDTO, error)
	PredictMenBodyScan(file multipart.File, filename string) (*dtos.MenBodyScanPredictionResponseDTO, error)
}

// aiService implements AIService
type aiService struct {
	config     *config.Config
	httpClient *http.Client
	llmChain   *LLMChain
}

// NewAIService creates a new AI service instance
func NewAIService(cfg *config.Config) AIService {
	if cfg == nil {
		panic("config cannot be nil")
	}
	if cfg.AIApiURL == "" {
		panic("AI_API_URL is required but not set in config")
	}

	// Initialize LLM providers in priority order: Groq -> Gemini -> Telkom
	var providers []LLMProvider

	// Add Groq if configured
	if cfg.GroqAPIKey != "" {
		model := cfg.GroqModel
		if model == "" {
			model = "llama-3.3-70b-versatile" // Default Groq model
		}
		providers = append(providers, NewGroqProvider(cfg.GroqAPIKey, model))
		log.Printf("[AIService] Groq provider initialized with model: %s", model)
	}

	// Add Gemini if configured
	if cfg.GeminiAPIKey != "" {
		model := cfg.GeminiModel
		if model == "" {
			model = "gemini-2.0-flash-exp" // Default Gemini model
		}
		providers = append(providers, NewGeminiProvider(cfg.GeminiAPIKey, model))
		log.Printf("[AIService] Gemini provider initialized with model: %s", model)
	}

	// Add Telkom if configured
	if cfg.TelkomLLMAPIKey != "" {
		model := cfg.TelkomModel
		if model == "" {
			model = "Qwen2.5-Coder-32B-Instruct" // Default Telkom model
		}
		providers = append(providers, NewTelkomProvider(cfg.TelkomLLMAPIKey, model))
		log.Printf("[AIService] Telkom provider initialized with model: %s", model)
	}

	var llmChain *LLMChain
	if len(providers) > 0 {
		llmChain = NewLLMChain(providers...)
		log.Printf("[AIService] LLM chain initialized with %d providers", len(providers))
	} else {
		log.Printf("[AIService] WARNING: No LLM providers configured, recommendations will be disabled")
	}

	return &aiService{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		llmChain: llmChain,
	}
}

// PredictSkinColorTone calls the skin color tone prediction API
func (s *aiService) PredictSkinColorTone(file multipart.File, filename string) (*dtos.SkinColorTonePredictionResponseDTO, error) {
	log.Printf("DEBUG: AI Service PredictSkinColorTone called with filename: %s", filename)

	// Nil checks
	if s == nil {
		log.Printf("ERROR: aiService is nil")
		return nil, fmt.Errorf("aiService is nil")
	}
	if s.config == nil {
		log.Printf("ERROR: config is nil")
		return nil, fmt.Errorf("config is nil")
	}
	if file == nil {
		log.Printf("ERROR: file parameter is nil")
		return nil, fmt.Errorf("file is nil")
	}
	if filename == "" {
		log.Printf("ERROR: filename is empty")
		return nil, fmt.Errorf("filename is empty")
	}

	endpoint := fmt.Sprintf("%s/predict/sct", s.config.AIApiURL)
	log.Printf("DEBUG: AI API endpoint: %s", endpoint)

	log.Printf("DEBUG: Creating multipart form...")
	body, contentType, err := s.createMultipartForm(file, filename, "file")
	if err != nil {
		log.Printf("ERROR: Failed to create multipart form: %v", err)
		return nil, fmt.Errorf("failed to create multipart form: %w", err)
	}

	log.Printf("DEBUG: Making HTTP request...")
	resp, err := s.makeRequest("POST", endpoint, body, contentType)
	if err != nil {
		log.Printf("ERROR: Failed to make request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("DEBUG: Got response with status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR: Non-OK status code: %d", resp.StatusCode)
		return nil, s.handleErrorResponse(resp)
	}

	log.Printf("DEBUG: Decoding response...")
	var result dtos.SkinColorTonePredictionResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("ERROR: Failed to decode response: %v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("DEBUG: AI Service success - Result: %+v", result)

	// Enrich with LLM color recommendations if available
	if s.llmChain != nil && result.SkinTone != "" {
		log.Printf("DEBUG: Getting color recommendations for skin tone: %s", result.SkinTone)
		colors, err := s.llmChain.GenerateColorRecommendations(result.SkinTone)
		if err != nil {
			log.Printf("WARNING: Failed to get color recommendations: %v", err)
			// Don't fail the whole request, just skip recommendations
		} else {
			result.ColorRecommendations = colors
			log.Printf("DEBUG: Added %d color recommendations", len(colors))
		}
	}

	return &result, nil
}

// PredictWomanBodyScan calls the woman body scan prediction API
func (s *aiService) PredictWomanBodyScan(file multipart.File, filename string) (*dtos.WomanBodyScanPredictionResponseDTO, error) {
	endpoint := fmt.Sprintf("%s/wbs/predict", s.config.AIApiURL)

	body, contentType, err := s.createMultipartForm(file, filename, "file")
	if err != nil {
		return nil, fmt.Errorf("failed to create multipart form: %w", err)
	}

	resp, err := s.makeRequest("POST", endpoint, body, contentType)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, s.handleErrorResponse(resp)
	}

	var result dtos.WomanBodyScanPredictionResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Enrich with LLM style recommendations if available
	if s.llmChain != nil && result.PredictedClass != "" {
		log.Printf("DEBUG: Getting style recommendations for body type: %s", result.PredictedClass)
		styles, err := s.llmChain.GenerateStyleRecommendations(result.PredictedClass)
		if err != nil {
			log.Printf("WARNING: Failed to get style recommendations: %v", err)
			// Don't fail the whole request, just skip recommendations
		} else {
			result.StyleRecommendations = styles
			log.Printf("DEBUG: Added %d style recommendations", len(styles))
		}
	}

	return &result, nil
}

// PredictMenBodyScan calls the men's body scan prediction API
func (s *aiService) PredictMenBodyScan(file multipart.File, filename string) (*dtos.MenBodyScanPredictionResponseDTO, error) {
	endpoint := fmt.Sprintf("%s/mbs/predict", s.config.AIApiURL)

	body, contentType, err := s.createMultipartForm(file, filename, "file")
	if err != nil {
		return nil, fmt.Errorf("failed to create multipart form: %w", err)
	}

	resp, err := s.makeRequest("POST", endpoint, body, contentType)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, s.handleErrorResponse(resp)
	}

	var result dtos.MenBodyScanPredictionResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Enrich with LLM style recommendations if available
	if s.llmChain != nil && result.PredictedClass != "" {
		log.Printf("DEBUG: Getting style recommendations for body type: %s", result.PredictedClass)
		styles, err := s.llmChain.GenerateStyleRecommendations(result.PredictedClass)
		if err != nil {
			log.Printf("WARNING: Failed to get style recommendations: %v", err)
			// Don't fail the whole request, just skip recommendations
		} else {
			result.StyleRecommendations = styles
			log.Printf("DEBUG: Added %d style recommendations", len(styles))
		}
	}

	return &result, nil
}

// createMultipartForm creates a multipart form with the uploaded file
func (s *aiService) createMultipartForm(file multipart.File, filename, fieldName string) (io.Reader, string, error) {
	// Reset file position to beginning
	if _, err := file.Seek(0, 0); err != nil {
		return nil, "", fmt.Errorf("failed to seek file: %w", err)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Create form file field
	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file content to form
	if _, err := io.Copy(part, file); err != nil {
		return nil, "", fmt.Errorf("failed to copy file content: %w", err)
	}

	// Close writer to finalize the form
	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	return &body, writer.FormDataContentType(), nil
}

// makeRequest makes HTTP request to AI API
func (s *aiService) makeRequest(method, url string, body io.Reader, contentType string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return resp, nil
}

// handleErrorResponse handles error responses from AI API
func (s *aiService) handleErrorResponse(resp *http.Response) error {
	var errorResp dtos.AIErrorResponseDTO

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("AI API error (status %d): failed to read error response", resp.StatusCode)
	}

	if err := json.Unmarshal(body, &errorResp); err != nil {
		return fmt.Errorf("AI API error (status %d): %s", resp.StatusCode, string(body))
	}

	return fmt.Errorf("AI API error (status %d): %s - %s", resp.StatusCode, errorResp.Error, errorResp.Message)
}
