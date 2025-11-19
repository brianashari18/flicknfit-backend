package controllers

import (
	"flicknfit_backend/dtos"
	"flicknfit_backend/services"
	"flicknfit_backend/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// FavoriteController defines HTTP handlers for favorite operations
type FavoriteController interface {
	GetUserFavorites(c *fiber.Ctx) error
	AddFavorite(c *fiber.Ctx) error
	RemoveFavorite(c *fiber.Ctx) error
	ToggleFavorite(c *fiber.Ctx) error
}

// favoriteController implements FavoriteController interface
type favoriteController struct {
	service   services.FavoriteService
	validator validator.Validate
}

// NewFavoriteController creates a new favorite controller
func NewFavoriteController(service services.FavoriteService, validator *validator.Validate) FavoriteController {
	return &favoriteController{
		service:   service,
		validator: *validator,
	}
}

// GetUserFavorites retrieves all user's favorites
// @Summary Get user's favorite products
// @Description Retrieve all products marked as favorites by the authenticated user
// @Tags Favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]dtos.FavoriteResponseDTO} "Favorites retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /favorites [get]
func (ctrl *favoriteController) GetUserFavorites(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusUnauthorized, err.Error(), nil)
	}

	favorites, err := ctrl.service.GetUserFavorites(userID)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to retrieve favorites", nil)
	}

	response := dtos.ToFavoriteResponseDTOs(favorites)
	return utils.SendResponse(c, http.StatusOK, "Favorites retrieved successfully", response)
}

// AddFavorite adds a product to user's favorites
// @Summary Add product to favorites
// @Description Add a product item to user's favorite list
// @Tags Favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param favorite body dtos.AddFavoriteDTO true "Product item to add to favorites"
// @Success 201 {object} utils.Response "Product added to favorites successfully"
// @Failure 400 {object} utils.Response "Invalid request body or validation failed"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Product item not found"
// @Failure 409 {object} utils.Response "Product already in favorites"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /favorites [post]
func (ctrl *favoriteController) AddFavorite(c *fiber.Ctx) error {
	var dto dtos.AddFavoriteDTO
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

	err = ctrl.service.AddFavorite(userID, &dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusCreated, "Added to favorites successfully", nil)
}

// RemoveFavorite removes a product from user's favorites
// @Summary Remove product from favorites
// @Description Remove a product item from user's favorite list
// @Tags Favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param productId path int true "Product Item ID to remove from favorites"
// @Success 200 {object} utils.Response "Removed from favorites successfully"
// @Failure 400 {object} utils.Response "Invalid product item ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Product not in favorites"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /favorites/{productId} [delete]
func (ctrl *favoriteController) RemoveFavorite(c *fiber.Ctx) error {
	productItemID, err := utils.GetUintParam(c, "productId")
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid product item ID", nil)
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusUnauthorized, err.Error(), nil)
	}

	err = ctrl.service.RemoveFavorite(userID, productItemID)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusOK, "Removed from favorites successfully", nil)
}

// ToggleFavorite toggles favorite status for a product
func (ctrl *favoriteController) ToggleFavorite(c *fiber.Ctx) error {
	var dto dtos.AddFavoriteDTO
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

	response, err := ctrl.service.ToggleFavorite(userID, &dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusOK, response.Message, response)
}
