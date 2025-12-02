package controllers

import (
	"flicknfit_backend/services"
	"flicknfit_backend/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ScanHistoryController handles scan history endpoints
type ScanHistoryController interface {
	GetFaceScanHistories(c *fiber.Ctx) error
	DeleteFaceScanHistory(c *fiber.Ctx) error
	GetBodyScanHistories(c *fiber.Ctx) error
	DeleteBodyScanHistory(c *fiber.Ctx) error
	GetScanImage(c *fiber.Ctx) error
}

type scanHistoryController struct {
	scanHistoryService services.ScanHistoryService
	storageService     services.SupabaseStorageService
}

// NewScanHistoryController creates a new scan history controller
func NewScanHistoryController(
	scanHistoryService services.ScanHistoryService,
	storageService services.SupabaseStorageService,
) ScanHistoryController {
	return &scanHistoryController{
		scanHistoryService: scanHistoryService,
		storageService:     storageService,
	}
}

// GetFaceScanHistories retrieves all face scan histories for authenticated user
// @Summary Get face scan histories
// @Description Get all face scan histories for the authenticated user
// @Tags Scan History
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dtos.FaceScanHistoryListDTO}
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /history/face [get]
func (ctrl *scanHistoryController) GetFaceScanHistories(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint64)

	histories, err := ctrl.scanHistoryService.GetFaceScanHistories(userID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve face scan histories")
	}

	return utils.SendSuccess(c, "Face scan histories retrieved successfully", histories)
}

// DeleteFaceScanHistory deletes a face scan history
// @Summary Delete face scan history
// @Description Delete a specific face scan history by ID
// @Tags Scan History
// @Produce json
// @Security BearerAuth
// @Param id path int true "Face Scan History ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response "Invalid ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - not owner"
// @Failure 404 {object} utils.Response "History not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /history/face/{id} [delete]
func (ctrl *scanHistoryController) DeleteFaceScanHistory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint64)

	historyID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid history ID")
	}

	if err := ctrl.scanHistoryService.DeleteFaceScanHistory(userID, historyID); err != nil {
		if err.Error() == "history not found: record not found" {
			return utils.SendError(c, fiber.StatusNotFound, "Face scan history not found")
		}
		if err.Error() == "unauthorized: scan does not belong to user" {
			return utils.SendError(c, fiber.StatusForbidden, "You don't have permission to delete this scan")
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to delete face scan history")
	}

	return utils.SendSuccess(c, "Face scan history deleted successfully", nil)
}

// GetBodyScanHistories retrieves all body scan histories for authenticated user
// @Summary Get body scan histories
// @Description Get all body scan histories for the authenticated user
// @Tags Scan History
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dtos.BodyScanHistoryListDTO}
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /history/body [get]
func (ctrl *scanHistoryController) GetBodyScanHistories(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint64)

	histories, err := ctrl.scanHistoryService.GetBodyScanHistories(userID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve body scan histories")
	}

	return utils.SendSuccess(c, "Body scan histories retrieved successfully", histories)
}

// DeleteBodyScanHistory deletes a body scan history
// @Summary Delete body scan history
// @Description Delete a specific body scan history by ID
// @Tags Scan History
// @Produce json
// @Security BearerAuth
// @Param id path int true "Body Scan History ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response "Invalid ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - not owner"
// @Failure 404 {object} utils.Response "History not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /history/body/{id} [delete]
func (ctrl *scanHistoryController) DeleteBodyScanHistory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint64)

	historyID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid history ID")
	}

	if err := ctrl.scanHistoryService.DeleteBodyScanHistory(userID, historyID); err != nil {
		if err.Error() == "history not found: record not found" {
			return utils.SendError(c, fiber.StatusNotFound, "Body scan history not found")
		}
		if err.Error() == "unauthorized: scan does not belong to user" {
			return utils.SendError(c, fiber.StatusForbidden, "You don't have permission to delete this scan")
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to delete body scan history")
	}

	return utils.SendSuccess(c, "Body scan history deleted successfully", nil)
}

// GetScanImage retrieves and serves a decrypted scan image
// @Summary Get scan image
// @Description Download and decrypt a scan image by file path
// @Tags Scan History
// @Produce image/jpeg
// @Security BearerAuth
// @Param path query string true "File path in storage"
// @Success 200 {file} binary "Image file"
// @Failure 400 {object} utils.Response "Missing file path"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "File not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /history/image [get]
func (ctrl *scanHistoryController) GetScanImage(c *fiber.Ctx) error {
	filePath := c.Query("path")
	if filePath == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "File path is required")
	}

	// Download and decrypt file
	imageData, err := ctrl.storageService.DownloadAndDecryptFile(filePath)
	if err != nil {
		if err.Error() == "file not found" || err.Error() == "failed to download file: not found" {
			return utils.SendError(c, fiber.StatusNotFound, "Image not found")
		}
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve image")
	}

	// Set content type and return image
	c.Set("Content-Type", "image/jpeg")
	return c.Send(imageData)
}
