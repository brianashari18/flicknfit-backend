package services

import (
	"encoding/json"
	"fmt"
)

// LLMProvider adalah interface untuk semua LLM provider
type LLMProvider interface {
	GenerateColorRecommendations(skinTone string) ([]string, error)
	GenerateStyleRecommendations(bodyType string) ([]string, error)
	GetName() string
}

// LLMRequest adalah struktur request umum untuk LLM
type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMResponse adalah struktur response umum dari LLM
type LLMResponse struct {
	Colors []string `json:"colors,omitempty"`
	Styles []string `json:"styles,omitempty"`
}

// ParseLLMResponse mencoba extract JSON dari response LLM
func ParseLLMResponse(content string) (*LLMResponse, error) {
	var response LLMResponse
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}
	return &response, nil
}

// GetColorPrompt generates prompt untuk color recommendations
func GetColorPrompt(skinTone string) string {
	return fmt.Sprintf(`You are a fashion color expert. Based on the skin tone "%s", provide exactly 15 color recommendations that would look best on this person.

Return your response as a JSON object with this exact format:
{
  "colors": ["color1", "color2", "color3", ...]
}

The colors should be specific color names (e.g., "Burgundy", "Sage Green", "Navy Blue") not generic terms. Provide exactly 15 colors.`, skinTone)
}

// GetStylePrompt generates prompt untuk style recommendations
func GetStylePrompt(bodyType string) string {
	return fmt.Sprintf(`You are a fashion style expert. Based on the body type "%s", provide 5-7 fashion style recommendations that would be most flattering for this body type.

Return your response as a JSON object with this exact format:
{
  "styles": ["style1", "style2", "style3", ...]
}

The styles should be specific style names (e.g., "A-Line Dresses", "Tailored Blazers", "High-Waisted Pants"). Provide 5 to 7 styles.`, bodyType)
}
