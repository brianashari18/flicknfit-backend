package services

import (
	"flicknfit_backend/models"
	"flicknfit_backend/repositories"
	"fmt"
	"net/url"
	"strings"
)

// TrackingService handles product click tracking and redirect logic
type TrackingService interface {
	TrackClick(userID, productID uint64, ipAddress, userAgent string) error
	GetRedirectURL(product *models.Product, brand *models.Brand) string
	GenerateWhatsAppLink(phoneNumber, message string) string
	GetClickStats(productID uint64) (int64, error)
	GetBrandClickStats(brandID uint64) (int64, error)
}

type trackingService struct {
	clickRepo   repositories.ProductClickRepository
	productRepo repositories.ProductRepository
	brandRepo   repositories.BrandRepository
}

// NewTrackingService creates a new tracking service
func NewTrackingService(
	clickRepo repositories.ProductClickRepository,
	productRepo repositories.ProductRepository,
	brandRepo repositories.BrandRepository,
) TrackingService {
	return &trackingService{
		clickRepo:   clickRepo,
		productRepo: productRepo,
		brandRepo:   brandRepo,
	}
}

// TrackClick records a product click event
func (s *trackingService) TrackClick(userID, productID uint64, ipAddress, userAgent string) error {
	// Get product to obtain brand_id
	product, err := s.productRepo.GetProductByID(productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	// Handle anonymous users (userID = 0) by setting nil
	var userIDPtr *uint64
	if userID > 0 {
		userIDPtr = &userID
	}

	click := &models.ProductClick{
		UserID:    userIDPtr,
		ProductID: productID,
		BrandID:   product.BrandID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	return s.clickRepo.Create(click)
}

// GetRedirectURL determines the best redirect URL based on available platform links
func (s *trackingService) GetRedirectURL(product *models.Product, brand *models.Brand) string {
	// Priority 1: Product-specific WhatsApp (if template exists)
	if product.WhatsAppTemplate != "" && brand.WhatsAppNumber != "" {
		message := strings.ReplaceAll(product.WhatsAppTemplate, "{product_name}", product.Name)
		message = strings.ReplaceAll(message, "{brand_name}", brand.Name)
		return s.GenerateWhatsAppLink(brand.WhatsAppNumber, message)
	}

	// Priority 2: Product-specific Tokopedia/Shopee/Instagram
	if product.TokopediaProductURL != "" {
		return product.TokopediaProductURL
	}
	if product.ShopeeProductURL != "" {
		return product.ShopeeProductURL
	}
	if product.InstagramProductURL != "" {
		return product.InstagramProductURL
	}

	// Priority 3: Brand-level WhatsApp (with default message)
	if brand.WhatsAppNumber != "" {
		defaultMessage := fmt.Sprintf(
			"Halo! Saya tertarik dengan produk %s dari FlickNFit. Apakah masih tersedia?",
			product.Name,
		)
		return s.GenerateWhatsAppLink(brand.WhatsAppNumber, defaultMessage)
	}

	// Priority 4: Brand-level Tokopedia/Shopee/Instagram
	if brand.TokopediaURL != "" {
		return brand.TokopediaURL
	}
	if brand.ShopeeURL != "" {
		return brand.ShopeeURL
	}
	if brand.InstagramURL != "" {
		return brand.InstagramURL
	}

	// Priority 5: Generic product URL
	if product.BrandProductURL != "" {
		return product.BrandProductURL
	}

	// Fallback: Brand website
	if brand.WebsiteURL != "" {
		return brand.WebsiteURL
	}

	// Last resort: FlickNFit homepage
	return "https://flicknfit.com"
}

// GenerateWhatsAppLink creates a WhatsApp deep link with pre-filled message
func (s *trackingService) GenerateWhatsAppLink(phoneNumber, message string) string {
	// Remove all non-numeric characters from phone number
	cleanPhone := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phoneNumber)

	// Ensure phone number starts with country code (62 for Indonesia)
	if !strings.HasPrefix(cleanPhone, "62") {
		// Remove leading 0 if exists
		cleanPhone = strings.TrimPrefix(cleanPhone, "0")
		cleanPhone = "62" + cleanPhone
	}

	// URL encode the message
	encodedMessage := url.QueryEscape(message)

	return fmt.Sprintf("https://wa.me/%s?text=%s", cleanPhone, encodedMessage)
}

// GetClickStats returns total clicks for a product
func (s *trackingService) GetClickStats(productID uint64) (int64, error) {
	return s.clickRepo.CountByProductID(productID)
}

// GetBrandClickStats returns total clicks for a brand
func (s *trackingService) GetBrandClickStats(brandID uint64) (int64, error) {
	return s.clickRepo.CountByBrandID(brandID)
}
