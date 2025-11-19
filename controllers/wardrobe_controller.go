package controllers

import (
	"flicknfit_backend/dtos"
	"flicknfit_backend/services"
	"flicknfit_backend/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// WardrobeController defines HTTP handlers for wardrobe operations
type WardrobeController interface {
	GetUserWardrobe(c *fiber.Ctx) error
	GetWardrobeByCategory(c *fiber.Ctx) error
	CreateWardrobeItem(c *fiber.Ctx) error
	UpdateWardrobeItem(c *fiber.Ctx) error
	DeleteWardrobeItem(c *fiber.Ctx) error
	GetWardrobeCategories(c *fiber.Ctx) error
}

// wardrobeController implements WardrobeController interface
type wardrobeController struct {
	service   services.WardrobeService
	validator validator.Validate
}

// NewWardrobeController creates a new wardrobe controller
func NewWardrobeController(service services.WardrobeService, validator *validator.Validate) WardrobeController {
	return &wardrobeController{
		service:   service,
		validator: *validator,
	}
}

// GetUserWardrobe retrieves user's complete wardrobe organized by category
// @Summary Get user's wardrobe
// @Description Retrieve user's complete wardrobe organized by clothing categories
// @Tags Wardrobe
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dtos.WardrobeResponseDTO} "Wardrobe retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /wardrobe [get]
func (ctrl *wardrobeController) GetUserWardrobe(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusUnauthorized, err.Error(), nil)
	}

	wardrobe, err := ctrl.service.GetUserWardrobe(userID)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	response := dtos.ToWardrobeResponseDTO(wardrobe)
	return utils.SendResponse(c, http.StatusOK, "Wardrobe retrieved successfully", response)
}

// GetWardrobeByCategory retrieves wardrobe items by specific category
// @Summary Get wardrobe items by category
// @Description Retrieve wardrobe items filtered by specific clothing category
// @Tags Wardrobe
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category path string true "Clothing category"
// @Success 200 {object} utils.Response{data=[]dtos.WardrobeItemResponseDTO} "Wardrobe items retrieved successfully"
// @Failure 400 {object} utils.Response "Category is required"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /wardrobe/category/{category} [get]
func (ctrl *wardrobeController) GetWardrobeByCategory(c *fiber.Ctx) error {
	category := c.Params("category")
	if category == "" {
		return utils.SendResponse(c, http.StatusBadRequest, "Category is required", nil)
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusUnauthorized, err.Error(), nil)
	}

	items, err := ctrl.service.GetWardrobeByCategory(userID, category)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	response := dtos.ToWardrobeItemResponseDTOs(items)
	return utils.SendResponse(c, http.StatusOK, "Wardrobe items retrieved successfully", response)
}

// CreateWardrobeItem creates a new wardrobe item
// @Summary Add item to wardrobe
// @Description Add a new clothing item to user's wardrobe
// @Tags Wardrobe
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item body dtos.CreateWardrobeItemDTO true "Wardrobe item data"
// @Success 201 {object} utils.Response "Item added to wardrobe successfully"
// @Failure 400 {object} utils.Response "Invalid request body or validation failed"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /wardrobe [post]
func (ctrl *wardrobeController) CreateWardrobeItem(c *fiber.Ctx) error {
	var dto dtos.CreateWardrobeItemDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}

	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), nil)
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusUnauthorized, err.Error(), nil)
	}

	err = ctrl.service.CreateWardrobeItem(userID, &dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusCreated, "Wardrobe item created successfully", nil)
}

// UpdateWardrobeItem updates an existing wardrobe item
// @Summary Update wardrobe item
// @Description Update an existing clothing item in user's wardrobe
// @Tags Wardrobe
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param itemId path string true "Wardrobe item ID"
// @Param item body dtos.UpdateWardrobeItemDTO true "Updated wardrobe item data"
// @Success 200 {object} utils.Response "Wardrobe item updated successfully"
// @Failure 400 {object} utils.Response "Invalid request body or item ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Wardrobe item not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /wardrobe/{itemId} [put]
func (ctrl *wardrobeController) UpdateWardrobeItem(c *fiber.Ctx) error {
	itemID, err := utils.GetUintParam(c, "itemId")
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid item ID", nil)
	}

	var dto dtos.UpdateWardrobeItemDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}

	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), nil)
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusUnauthorized, err.Error(), nil)
	}

	err = ctrl.service.UpdateWardrobeItem(userID, itemID, &dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusOK, "Wardrobe item updated successfully", nil)
}

// DeleteWardrobeItem deletes a wardrobe item
// @Summary Delete wardrobe item
// @Description Remove a clothing item from user's wardrobe
// @Tags Wardrobe
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param itemId path string true "Wardrobe item ID"
// @Success 200 {object} utils.Response "Wardrobe item deleted successfully"
// @Failure 400 {object} utils.Response "Invalid item ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Wardrobe item not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /wardrobe/{itemId} [delete]
func (ctrl *wardrobeController) DeleteWardrobeItem(c *fiber.Ctx) error {
	itemID, err := utils.GetUintParam(c, "itemId")
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid item ID", nil)
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusUnauthorized, err.Error(), nil)
	}

	err = ctrl.service.DeleteWardrobeItem(userID, itemID)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusOK, "Wardrobe item deleted successfully", nil)
}

// GetWardrobeCategories retrieves unique categories in user's wardrobe
// @Summary Get wardrobe categories
// @Description Retrieve all unique clothing categories in user's wardrobe
// @Tags Wardrobe
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]string} "Wardrobe categories retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /wardrobe/categories [get]
func (ctrl *wardrobeController) GetWardrobeCategories(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusUnauthorized, err.Error(), nil)
	}

	categories, err := ctrl.service.GetWardrobeCategories(userID)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusOK, "Wardrobe categories retrieved successfully", categories)
}
