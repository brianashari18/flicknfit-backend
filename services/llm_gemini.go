package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiProvider implements LLMProvider for Google Gemini
type GeminiProvider struct {
	APIKey string
	Model  string
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(apiKey, model string) *GeminiProvider {
	return &GeminiProvider{
		APIKey: apiKey,
		Model:  model,
	}
}

// GetName returns provider name
func (g *GeminiProvider) GetName() string {
	return "Gemini"
}

// invoke calls Gemini API
func (g *GeminiProvider) invoke(systemPrompt, userPrompt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := genai.NewClient(ctx, option.WithAPIKey(g.APIKey))
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(g.Model)

	// Set generation config
	model.GenerationConfig = genai.GenerationConfig{
		Temperature:     genai.Ptr(float32(0.3)),
		MaxOutputTokens: genai.Ptr(int32(2000)),
	}

	// Set system instruction
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text("You are a helpful assistant for fashion recommendations. You MUST respond ONLY with a valid JSON object, do NOT include any markdown code blocks or surrounding text."),
			genai.Text(systemPrompt),
		},
	}

	// Generate content
	resp, err := model.GenerateContent(ctx, genai.Text(userPrompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	if len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response content")
	}

	content := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

	// Remove markdown code blocks if present
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	return content, nil
}

// GenerateColorRecommendations generates color recommendations based on skin tone
func (g *GeminiProvider) GenerateColorRecommendations(skinTone string) ([]string, error) {
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
func (g *GeminiProvider) GenerateStyleRecommendations(bodyType string) ([]string, error) {
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
