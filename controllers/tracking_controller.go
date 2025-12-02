package controllers

import (
	"flicknfit_backend/repositories"
	"flicknfit_backend/services"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// TrackingController handles product click tracking endpoints
type TrackingController interface {
	TrackProductClick(c *fiber.Ctx) error
}

type trackingController struct {
	trackingService services.TrackingService
	productService  services.ProductService
	brandRepo       repositories.BrandRepository
}

// NewTrackingController creates a new tracking controller
func NewTrackingController(
	trackingService services.TrackingService,
	productService services.ProductService,
	brandRepo repositories.BrandRepository,
) TrackingController {
	return &trackingController{
		trackingService: trackingService,
		productService:  productService,
		brandRepo:       brandRepo,
	}
}

// TrackProductClick tracks a product click and redirects to the appropriate platform
// @Summary Track product click and redirect
// @Description Records a click event for analytics and redirects user to brand store (WhatsApp/Instagram/Tokopedia/etc)
// @Tags Tracking
// @Produce html
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 302 "Redirect to brand store"
// @Failure 404 {string} string "Product not found"
// @Failure 500 {string} string "Internal server error"
// @Router /track/click/product/{id} [get]
func (ctrl *trackingController) TrackProductClick(c *fiber.Ctx) error {
	// Get product ID from URL
	productID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		log.Printf("[TrackingController] Invalid product ID: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}

	// Get user ID from context (optional - may be anonymous)
	var userID uint64
	if userIDValue := c.Locals("userID"); userIDValue != nil {
		userID = userIDValue.(uint64)
	}

	// Get product details
	product, err := ctrl.productService.GetProductPublicByID(productID)
	if err != nil {
		log.Printf("[TrackingController] Product not found: %v", err)
		return c.Status(fiber.StatusNotFound).SendString("Product not found")
	}

	// Get brand details
	brand, err := ctrl.brandRepo.GetBrandByID(product.BrandID)
	if err != nil {
		log.Printf("[TrackingController] Brand not found: %v", err)
		return c.Status(fiber.StatusNotFound).SendString("Brand not found")
	}

	// Track the click (best effort - don't fail if tracking fails)
	// Support both authenticated and anonymous tracking
	err = ctrl.trackingService.TrackClick(
		userID, // 0 for anonymous users
		productID,
		c.IP(),
		c.Get("User-Agent"),
	)
	if err != nil {
		log.Printf("[TrackingController] Failed to track click: %v", err)
		// Continue with redirect even if tracking fails
	} else {
		if userID > 0 {
			log.Printf("[TrackingController] Click tracked: user=%d, product=%d, brand=%d", userID, productID, product.BrandID)
		} else {
			log.Printf("[TrackingController] Anonymous click tracked: product=%d, brand=%d", productID, product.BrandID)
		}
	}

	// Determine redirect URL based on available platforms
	redirectURL := ctrl.trackingService.GetRedirectURL(product, brand)

	log.Printf("[TrackingController] Redirecting to: %s", redirectURL)

	// Redirect to brand store
	return c.Redirect(redirectURL, fiber.StatusFound) // 302 redirect
}
