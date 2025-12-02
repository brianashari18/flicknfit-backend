package controllers

import (
	"flicknfit_backend/services"
	"flicknfit_backend/utils"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// AIController defines the HTTP handlers for AI prediction operations using Fiber.
type AIController interface {
	PredictSkinColorTone(c *fiber.Ctx) error
	PredictWomanBodyScan(c *fiber.Ctx) error
	PredictMenBodyScan(c *fiber.Ctx) error
}

// aiController is the implementation of AIController.
type aiController struct {
	aiService          services.AIService
	scanHistoryService services.ScanHistoryService
}

// NewAIController creates and returns a new instance of AIController.
func NewAIController(aiService services.AIService, scanHistoryService services.ScanHistoryService) AIController {
	if aiService == nil {
		panic("aiService cannot be nil")
	}
	return &aiController{
		aiService:          aiService,
		scanHistoryService: scanHistoryService,
	}
}

// PredictSkinColorTone handles skin color tone prediction from uploaded image.
// @Summary Predict skin color tone
// @Description Analyze uploaded image to predict user's skin color tone
// @Tags AI Predictions
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Image file for skin color tone analysis"
// @Success 200 {object} utils.Response "Skin color tone predicted successfully"
// @Failure 400 {object} utils.Response "Invalid file or request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error or AI API error"
// @Router /ai/predict/skin-color-tone [post]
func (ctrl *aiController) PredictSkinColorTone(c *fiber.Ctx) error {
	// Debug logging
	log.Printf("DEBUG: PredictSkinColorTone called")

	// Check for nil controller
	if ctrl == nil {
		log.Printf("ERROR: Controller is nil")
		return utils.SendResponse(c, http.StatusInternalServerError, "Controller not available", nil)
	}

	// Check for nil service (safety check)
	if ctrl.aiService == nil {
		log.Printf("ERROR: AI service is nil")
		return utils.SendResponse(c, http.StatusInternalServerError, "AI service not available", nil)
	}

	log.Printf("DEBUG: Getting uploaded file...")

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("ERROR: Failed to get form file: %v", err)
		return utils.SendResponse(c, http.StatusBadRequest, "No file uploaded or invalid file", nil)
	}

	// Check if file is nil
	if file == nil {
		log.Printf("ERROR: File is nil")
		return utils.SendResponse(c, http.StatusBadRequest, "File is nil", nil)
	}

	log.Printf("DEBUG: File received - Name: %s, Size: %d", file.Filename, file.Size)

	// Check if file.Header is nil
	if file.Header == nil {
		log.Printf("ERROR: File header is nil")
		return utils.SendResponse(c, http.StatusBadRequest, "File header is nil", nil)
	}

	// Validate file type (basic image validation)
	contentType := file.Header.Get("Content-Type")
	log.Printf("DEBUG: Content-Type: %s", contentType)

	if !isValidImageType(contentType) {
		log.Printf("ERROR: Invalid content type: %s", contentType)
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid file type. Only image files are allowed", nil)
	}

	log.Printf("DEBUG: Opening file...")

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		log.Printf("ERROR: Failed to open file: %v", err)
		return utils.SendResponse(c, http.StatusBadRequest, "Failed to open uploaded file: "+err.Error(), nil)
	}
	if src == nil {
		log.Printf("ERROR: File source is nil")
		return utils.SendResponse(c, http.StatusBadRequest, "File source is nil", nil)
	}
	defer src.Close()

	log.Printf("DEBUG: Calling AI service...")

	// Call AI service with error recovery
	result, err := ctrl.aiService.PredictSkinColorTone(src, file.Filename)
	if err != nil {
		log.Printf("ERROR: AI service error: %v", err)
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to predict skin color tone: "+err.Error(), nil)
	}

	if result == nil {
		log.Printf("ERROR: AI service returned nil result")
		return utils.SendResponse(c, http.StatusInternalServerError, "AI service returned nil result", nil)
	}

	log.Printf("DEBUG: Success - Result: %+v", result)

	// Auto-save to history if user is authenticated and scanHistoryService is available
	if ctrl.scanHistoryService != nil {
		if userID, ok := c.Locals("userID").(uint64); ok && userID > 0 {
			// Reopen file for upload (src was consumed by PredictSkinColorTone)
			srcForUpload, err := file.Open()
			if err == nil {
				defer srcForUpload.Close()
				if err := ctrl.scanHistoryService.SaveFaceScanHistory(userID, srcForUpload, file.Filename, result); err != nil {
					log.Printf("WARNING: Failed to save face scan history: %v", err)
					// Don't fail the request, just log the warning
				} else {
					log.Printf("DEBUG: Face scan history saved for user %d", userID)
				}
			}
		}
	}

	return utils.SendResponse(c, http.StatusOK, "Skin color tone predicted successfully", result)
}

// PredictWomanBodyScan handles woman body scanning from uploaded image.
// @Summary Scan woman's body measurements
// @Description Analyze uploaded image to predict body type and measurements for women
// @Tags AI Predictions
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Image file for body scanning"
// @Success 200 {object} utils.Response "Body measurements predicted successfully"
// @Failure 400 {object} utils.Response "Invalid file or request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error or AI API error"
// @Router /ai/predict/woman-body-scan [post]
func (ctrl *aiController) PredictWomanBodyScan(c *fiber.Ctx) error {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "No file uploaded or invalid file", nil)
	}

	// Validate file type
	if !isValidImageType(file.Header.Get("Content-Type")) {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid file type. Only image files are allowed", nil)
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Failed to open uploaded file", nil)
	}
	defer src.Close()

	// Call AI service
	result, err := ctrl.aiService.PredictWomanBodyScan(src, file.Filename)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to scan woman body: "+err.Error(), nil)
	}

	// Auto-save to history if user is authenticated and scanHistoryService is available
	if ctrl.scanHistoryService != nil {
		if userID, ok := c.Locals("userID").(uint64); ok && userID > 0 {
			// Reopen file for upload
			srcForUpload, err := file.Open()
			if err == nil {
				defer srcForUpload.Close()
				if err := ctrl.scanHistoryService.SaveBodyScanHistory(userID, srcForUpload, file.Filename, "woman", result); err != nil {
					log.Printf("WARNING: Failed to save body scan history: %v", err)
					// Don't fail the request, just log the warning
				} else {
					log.Printf("DEBUG: Woman body scan history saved for user %d", userID)
				}
			}
		}
	}

	return utils.SendResponse(c, http.StatusOK, "Body measurements predicted successfully", result)
}

// PredictMenBodyScan handles men's body scanning from uploaded image.
// @Summary Scan men's body measurements
// @Description Analyze uploaded image to predict body type and measurements for men
// @Tags AI Predictions
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Image file for body scanning"
// @Success 200 {object} utils.Response "Body measurements predicted successfully"
// @Failure 400 {object} utils.Response "Invalid file or request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error or AI API error"
// @Router /ai/predict/men-body-scan [post]
func (ctrl *aiController) PredictMenBodyScan(c *fiber.Ctx) error {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "No file uploaded or invalid file", nil)
	}

	// Validate file type
	if !isValidImageType(file.Header.Get("Content-Type")) {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid file type. Only image files are allowed", nil)
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Failed to open uploaded file", nil)
	}
	defer src.Close()

	// Call AI service
	result, err := ctrl.aiService.PredictMenBodyScan(src, file.Filename)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to scan body measurements: "+err.Error(), nil)
	}

	// Auto-save to history if user is authenticated and scanHistoryService is available
	if ctrl.scanHistoryService != nil {
		if userID, ok := c.Locals("userID").(uint64); ok && userID > 0 {
			// Reopen file for upload
			srcForUpload, err := file.Open()
			if err == nil {
				defer srcForUpload.Close()
				if err := ctrl.scanHistoryService.SaveBodyScanHistory(userID, srcForUpload, file.Filename, "man", result); err != nil {
					log.Printf("WARNING: Failed to save body scan history: %v", err)
					// Don't fail the request, just log the warning
				} else {
					log.Printf("DEBUG: Men body scan history saved for user %d", userID)
				}
			}
		}
	}

	return utils.SendResponse(c, http.StatusOK, "Body measurements predicted successfully", result)
}

// isValidImageType checks if the file type is a valid image format
func isValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	return validTypes[contentType]
}
