package services

import (
	"fmt"
	"log"
)

// LLMChain tries multiple LLM providers in order with fallback
type LLMChain struct {
	providers []LLMProvider
}

// NewLLMChain creates a new LLM chain with ordered providers
func NewLLMChain(providers ...LLMProvider) *LLMChain {
	return &LLMChain{
		providers: providers,
	}
}

// GenerateColorRecommendations tries each provider until success
func (c *LLMChain) GenerateColorRecommendations(skinTone string) ([]string, error) {
	var lastErr error

	for _, provider := range c.providers {
		log.Printf("[LLMChain] Trying %s for color recommendations (skinTone: %s)", provider.GetName(), skinTone)

		colors, err := provider.GenerateColorRecommendations(skinTone)
		if err != nil {
			log.Printf("[LLMChain] %s failed: %v", provider.GetName(), err)
			lastErr = err
			continue
		}

		log.Printf("[LLMChain] %s succeeded, got %d colors", provider.GetName(), len(colors))
		return colors, nil
	}

	return nil, fmt.Errorf("all providers failed, last error: %w", lastErr)
}

// GenerateStyleRecommendations tries each provider until success
func (c *LLMChain) GenerateStyleRecommendations(bodyType string) ([]string, error) {
	var lastErr error

	for _, provider := range c.providers {
		log.Printf("[LLMChain] Trying %s for style recommendations (bodyType: %s)", provider.GetName(), bodyType)

		styles, err := provider.GenerateStyleRecommendations(bodyType)
		if err != nil {
			log.Printf("[LLMChain] %s failed: %v", provider.GetName(), err)
			lastErr = err
			continue
		}

		log.Printf("[LLMChain] %s succeeded, got %d styles", provider.GetName(), len(styles))
		return styles, nil
	}

	return nil, fmt.Errorf("all providers failed, last error: %w", lastErr)
}
