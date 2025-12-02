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
	return fmt.Sprintf(`You are a professional fashion color analyst specializing in seasonal color analysis. Based on the skin tone "%s", provide exactly 15 color recommendations in hexadecimal format.

IMPORTANT GUIDELINES:
- For WARM tones (Warm Spring, Warm Autumn): Use earthy, rich colors with yellow/golden undertones (terracotta, olive, camel, rust, golden yellows, warm browns)
- For COOL tones (Cool Summer, Cool Winter): Use colors with blue/purple undertones (navy, emerald, plum, icy blues, cool grays)
- For LIGHT tones (Light Spring, Light Summer): Use soft, delicate colors (pastels, light neutrals)
- For DARK/DEEP tones (Dark Winter, Dark Autumn): Use deep, saturated colors (burgundy, forest green, charcoal)
- AVOID pink/pastel colors for Warm Autumn/Spring (use terracotta, coral, peach instead)
- AVOID warm colors for Cool Summer/Winter (use berry tones, cool reds instead)

Return your response as a JSON object with this exact format:
{
  "colors": ["#hexcode1", "#hexcode2", "#hexcode3", ...]
}

Provide exactly 15 hex color codes with the # prefix. Ensure all colors match the undertone and depth of the specified skin tone category.`, skinTone)
}

// GetStylePrompt generates prompt untuk style recommendations
func GetStylePrompt(bodyType string) string {
	return fmt.Sprintf(`You are a professional fashion stylist specializing in body type analysis. Based on the body type "%s", provide 5-7 clothing category recommendations that would be most flattering.

AVAILABLE CLOTHING CATEGORIES (you MUST choose from these):
- T-Shirts
- Shirts
- Pants
- Jackets
- Cardigans
- Skirts
- Tank Tops
- Sweaters
- Jeans
- Crop Tops
- Dress

GUIDELINES:
- Recommend specific categories that complement the body type
- Consider proportions, silhouettes, and flattering fits
- For example: "Hourglass" might benefit from "Dress", "Jeans", "Crop Tops"
- For "Pear shape" might benefit from "Jackets", "Shirts", "Cardigans"

Return your response as a JSON object with this exact format:
{
  "styles": ["Category1", "Category2", "Category3", ...]
}

Provide 5 to 7 category names from the list above. Use the exact category names as listed.`, bodyType)
}
