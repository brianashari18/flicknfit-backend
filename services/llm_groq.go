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

// GroqProvider implements LLMProvider for Groq
type GroqProvider struct {
	APIKey string
	Model  string
	client *http.Client
}

// NewGroqProvider creates a new Groq provider
func NewGroqProvider(apiKey, model string) *GroqProvider {
	return &GroqProvider{
		APIKey: apiKey,
		Model:  model,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

// GetName returns provider name
func (g *GroqProvider) GetName() string {
	return "Groq"
}

// groqRequest adalah struktur request untuk Groq API (OpenAI-compatible)
type groqRequest struct {
	Model       string       `json:"model"`
	Messages    []LLMMessage `json:"messages"`
	MaxTokens   int          `json:"max_tokens"`
	Temperature float64      `json:"temperature"`
}

// groqResponse adalah struktur response dari Groq API
type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// invoke calls Groq API
func (g *GroqProvider) invoke(systemPrompt, userPrompt string) (string, error) {
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

	payload := groqRequest{
		Model:       g.Model,
		Messages:    messages,
		MaxTokens:   2000,
		Temperature: 0.3,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+g.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
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

	var groqResp groqResponse
	if err := json.Unmarshal(body, &groqResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if groqResp.Error != nil {
		return "", fmt.Errorf("API error: %s", groqResp.Error.Message)
	}

	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	content := groqResp.Choices[0].Message.Content

	// Remove markdown code blocks if present
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	return content, nil
}

// GenerateColorRecommendations generates color recommendations based on skin tone
func (g *GroqProvider) GenerateColorRecommendations(skinTone string) ([]string, error) {
	prompt := GetColorPrompt(skinTone)
	content, err := g.invoke("You are a fashion color expert.", prompt)
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
func (g *GroqProvider) GenerateStyleRecommendations(bodyType string) ([]string, error) {
	prompt := GetStylePrompt(bodyType)
	content, err := g.invoke("You are a fashion style expert.", prompt)
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
