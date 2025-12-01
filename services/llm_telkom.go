package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// TelkomLLMProvider implements LLMProvider for Telkom AI
type TelkomLLMProvider struct {
	APIURL string
	APIKey string
	Model  string
	client *http.Client
}

// NewTelkomLLMProvider creates a new Telkom LLM provider
func NewTelkomLLMProvider(apiURL, apiKey, model string) *TelkomLLMProvider {
	return &TelkomLLMProvider{
		APIURL: apiURL,
		APIKey: apiKey,
		Model:  model,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

// GetName returns provider name
func (t *TelkomLLMProvider) GetName() string {
	return "Telkom AI"
}

// telkomRequest adalah struktur request untuk Telkom LLM
type telkomRequest struct {
	Model       string       `json:"model"`
	Messages    []LLMMessage `json:"messages"`
	MaxTokens   int          `json:"max_tokens"`
	Temperature float64      `json:"temperature"`
	Stream      bool         `json:"stream"`
}

// telkomResponse adalah struktur response dari Telkom LLM
type telkomResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// invoke calls Telkom LLM API
func (t *TelkomLLMProvider) invoke(systemPrompt, userPrompt string) (string, error) {
	messages := []LLMMessage{
		{
			Role:    "system",
			Content: "You are a helpful assistant for fashion recommendations. You MUST respond ONLY with a valid JSON object, do NOT include any markdown code blocks or surrounding text.",
		},
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: userPrompt,
		},
	}

	payload := telkomRequest{
		Model:       t.Model,
		Messages:    messages,
		MaxTokens:   2000,
		Temperature: 0.3,
		Stream:      false,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", t.APIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", t.APIKey)

	resp, err := t.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var telkomResp telkomResponse
	if err := json.Unmarshal(body, &telkomResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if telkomResp.Error != nil {
		return "", fmt.Errorf("API error: %s", telkomResp.Error.Message)
	}

	if len(telkomResp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	content := telkomResp.Choices[0].Message.Content

	// Remove markdown code blocks if present
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	return content, nil
}

// GenerateColorRecommendations generates color recommendations based on skin tone
func (t *TelkomLLMProvider) GenerateColorRecommendations(skinTone string) ([]string, error) {
	prompt := GetColorPrompt(skinTone)
	content, err := t.invoke("You are a fashion color expert.", prompt)
	if err != nil {
		return nil, err
	}

	llmResp, err := ParseLLMResponse(content)
	if err != nil {
		return nil, err
	}

	if len(llmResp.Colors) == 0 {
		return nil, fmt.Errorf("no colors returned from LLM")
	}

	return llmResp.Colors, nil
}

// GenerateStyleRecommendations generates style recommendations based on body type
func (t *TelkomLLMProvider) GenerateStyleRecommendations(bodyType string) ([]string, error) {
	prompt := GetStylePrompt(bodyType)
	content, err := t.invoke("You are a fashion style expert.", prompt)
	if err != nil {
		return nil, err
	}

	llmResp, err := ParseLLMResponse(content)
	if err != nil {
		return nil, err
	}

	if len(llmResp.Styles) == 0 {
		return nil, fmt.Errorf("no styles returned from LLM")
	}

	return llmResp.Styles, nil
}
