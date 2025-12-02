package controllers

import (
	"flicknfit_backend/dtos"
	"flicknfit_backend/services"
	"flicknfit_backend/utils"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// SavedItemsController defines the HTTP handlers for saved items operations using Fiber.
type SavedItemsController interface {
	GetUserSavedItems(c *fiber.Ctx) error
	AddProductItemToSavedItems(c *fiber.Ctx) error
	UpdateProductItemInSavedItems(c *fiber.Ctx) error
	RemoveProductItemFromSavedItems(c *fiber.Ctx) error
}

// SavedItemsController is the implementation of SavedItemsController.
type savedItemsController struct {
	SavedItemsService services.SavedItemsService
	validator         *validator.Validate
}

// NewSavedItemsController creates and returns a new instance of SavedItemsController.
func NewSavedItemsController(SavedItemsService services.SavedItemsService, validator *validator.Validate) SavedItemsController {
	return &savedItemsController{
		SavedItemsService: SavedItemsService,
		validator:         validator,
	}
}

// GetUserSavedItems handles fetching the user's saved items.
// @Summary Get user's saved items
// @Description Retrieve all items in the authenticated user's saved items
// @Tags Saved items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dtos.SavedItemsDTO} "saved items retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /savedItems [get]
func (ctrl *savedItemsController) GetUserSavedItems(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to get user ID from token", nil)
	}

	savedItems, err := ctrl.SavedItemsService.GetSavedItems(userID)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to retrieve saved items", nil)
	}

	response := dtos.ToSavedItemsDTO(savedItems)
	return utils.SendResponse(c, http.StatusOK, "saved items retrieved successfully", response)
}

// AddProductItemToSavedItems handles adding an item to the saved items.
// @Summary Add product to savedItems
// @Description Add a product item to the authenticated user's saved items
// @Tags Saved items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item body dtos.AddProductItemToSavedItemsRequestDTO true "Product item to add"
// @Success 200 {object} utils.Response{data=dtos.SavedItemsDTO} "Product item added to savedItems successfully"
// @Failure 400 {object} utils.Response "Invalid request body or validation failed"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /savedItems [post]
func (ctrl *savedItemsController) AddProductItemToSavedItems(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to get user ID from token", nil)
	}

	var dto dtos.AddProductItemToSavedItemsRequestDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}
	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
	}

	savedItems, err := ctrl.SavedItemsService.AddProductItemToSavedItems(userID, &dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to add item to savedItems: "+err.Error(), nil)
	}

	response := dtos.ToSavedItemsDTO(savedItems)
	return utils.SendResponse(c, http.StatusOK, "Product item added to savedItems successfully", response)
}

// UpdateProductItemInSavedItems handles updating the quantity of a product item in the savedItems.
// @Summary Update savedItems item
// @Description Update quantity of a product item in the saved items
// @Tags Saved items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param itemId path string true "savedItems Item ID"
// @Param item body dtos.UpdateProductItemInSavedItemsRequestDTO true "Updated item data"
// @Success 200 {object} utils.Response{data=dtos.SavedItemsDTO} "savedItems item updated successfully"
// @Failure 400 {object} utils.Response "Invalid request body or item ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "savedItems item not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /savedItems/{itemId} [put]
func (ctrl *savedItemsController) UpdateProductItemInSavedItems(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to get user ID from token", nil)
	}

	savedItemsItemID, err := strconv.ParseUint(c.Params("itemId"), 10, 64)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid savedItems item ID", nil)
	}

	var dto dtos.UpdateProductItemInSavedItemsRequestDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}
	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
	}

	savedItems, err := ctrl.SavedItemsService.UpdateProductItemInSavedItems(userID, savedItemsItemID, &dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusUnauthorized, "Failed to update savedItems item: "+err.Error(), nil)
	}

	response := dtos.ToSavedItemsDTO(savedItems)
	return utils.SendResponse(c, http.StatusOK, "savedItems item updated successfully", response)
}

// RemoveProductItemFromSavedItems handles removing an item from the saved items.
// @Summary Remove item from savedItems
// @Description Remove a product item from the saved items
// @Tags Saved items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param itemId path string true "savedItems Item ID"
// @Success 200 {object} utils.Response "savedItems item removed successfully"
// @Failure 400 {object} utils.Response "Invalid savedItems item ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "savedItems item not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /savedItems/{itemId} [delete]
func (ctrl *savedItemsController) RemoveProductItemFromSavedItems(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to get user ID from token", nil)
	}

	savedItemsItemID, err := strconv.ParseUint(c.Params("itemId"), 10, 64)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid savedItems item ID", nil)
	}

	if err := ctrl.SavedItemsService.RemoveProductItemFromSavedItems(userID, savedItemsItemID); err != nil {
		return utils.SendResponse(c, http.StatusUnauthorized, "Failed to remove item from savedItems: "+err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusOK, "savedItems item removed successfully", nil)
}
