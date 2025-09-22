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

// ShoppingCartController defines the HTTP handlers for shopping cart operations using Fiber.
type ShoppingCartController interface {
	GetUserCart(c *fiber.Ctx) error
	AddProductItemToCart(c *fiber.Ctx) error
	UpdateProductItemInCart(c *fiber.Ctx) error
	RemoveProductItemFromCart(c *fiber.Ctx) error
}

// shoppingCartController is the implementation of ShoppingCartController.
type shoppingCartController struct {
	shoppingCartService services.ShoppingCartService
	validator           *validator.Validate
}

// NewShoppingCartController creates and returns a new instance of ShoppingCartController.
func NewShoppingCartController(shoppingCartService services.ShoppingCartService, validator *validator.Validate) ShoppingCartController {
	return &shoppingCartController{
		shoppingCartService: shoppingCartService,
		validator:           validator,
	}
}

// GetUserCart handles fetching the user's shopping cart.
func (ctrl *shoppingCartController) GetUserCart(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to get user ID from token", nil)
	}

	cart, err := ctrl.shoppingCartService.GetCartItems(userID)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to retrieve shopping cart", nil)
	}

	response := dtos.ToShoppingCartDTO(cart)
	return utils.SendResponse(c, http.StatusOK, "Shopping cart retrieved successfully", response)
}

// AddProductItemToCart handles adding an item to the shopping cart.
func (ctrl *shoppingCartController) AddProductItemToCart(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to get user ID from token", nil)
	}

	var dto dtos.AddProductItemToCartRequestDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}
	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
	}

	cart, err := ctrl.shoppingCartService.AddProductItemToCart(userID, &dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to add item to cart: "+err.Error(), nil)
	}

	response := dtos.ToShoppingCartDTO(cart)
	return utils.SendResponse(c, http.StatusOK, "Product item added to cart successfully", response)
}

// UpdateProductItemInCart handles updating the quantity of a product item in the cart.
func (ctrl *shoppingCartController) UpdateProductItemInCart(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to get user ID from token", nil)
	}

	cartItemID, err := strconv.ParseUint(c.Params("itemId"), 10, 64)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid cart item ID", nil)
	}

	var dto dtos.UpdateProductItemInCartRequestDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}
	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
	}

	cart, err := ctrl.shoppingCartService.UpdateProductItemInCart(userID, cartItemID, &dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusUnauthorized, "Failed to update cart item: "+err.Error(), nil)
	}

	response := dtos.ToShoppingCartDTO(cart)
	return utils.SendResponse(c, http.StatusOK, "Cart item updated successfully", response)
}

// RemoveProductItemFromCart handles removing an item from the shopping cart.
func (ctrl *shoppingCartController) RemoveProductItemFromCart(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to get user ID from token", nil)
	}

	cartItemID, err := strconv.ParseUint(c.Params("itemId"), 10, 64)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid cart item ID", nil)
	}

	if err := ctrl.shoppingCartService.RemoveProductItemFromCart(userID, cartItemID); err != nil {
		return utils.SendResponse(c, http.StatusUnauthorized, "Failed to remove item from cart: "+err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusOK, "Cart item removed successfully", nil)
}
