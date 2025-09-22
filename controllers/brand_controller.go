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

// BrandController defines the HTTP handlers for brand-related operations using Fiber.
type BrandController interface {
	// Admin Routes
	AdminCreateBrand(c *fiber.Ctx) error
	AdminUpdateBrand(c *fiber.Ctx) error
	AdminDeleteBrand(c *fiber.Ctx) error

	// Public Routes
	GetAllBrands(c *fiber.Ctx) error
	GetBrandByID(c *fiber.Ctx) error
}

// brandController is the implementation of BrandController.
type brandController struct {
	service   services.BrandService
	validator validator.Validate
}

// NewBrandController creates and returns a new instance of BrandController.
func NewBrandController(service services.BrandService, validator *validator.Validate) BrandController {
	return &brandController{
		service:   service,
		validator: *validator,
	}
}

// AdminCreateBrand handles the creation of a new brand by an admin.
func (ctrl *brandController) AdminCreateBrand(c *fiber.Ctx) error {
	var dto dtos.BrandCreateRequestDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}

	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
	}

	brand, err := ctrl.service.AdminCreateBrand(&dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusCreated, "Brand created successfully", dtos.ToBrandResponseDTO(*brand))
}

// AdminUpdateBrand handles the update of an existing brand by an admin.
func (ctrl *brandController) AdminUpdateBrand(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid brand ID", nil)
	}

	var dto dtos.BrandUpdateRequestDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}

	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
	}

	brand, err := ctrl.service.AdminUpdateBrand(id, &dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusOK, "Brand updated successfully", dtos.ToBrandResponseDTO(*brand))
}

// AdminDeleteBrand handles the deletion of a brand by an admin.
func (ctrl *brandController) AdminDeleteBrand(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid brand ID", nil)
	}

	if err := ctrl.service.AdminDeleteBrand(id); err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusOK, "Brand deleted successfully", nil)
}

// GetAllBrands retrieves all brands for public display.
func (ctrl *brandController) GetAllBrands(c *fiber.Ctx) error {
	brands, err := ctrl.service.GetAllBrands()
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusOK, "Brands retrieved successfully", brands)
}

// GetBrandByID retrieves a single brand by its ID for public display.
func (ctrl *brandController) GetBrandByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid brand ID", nil)
	}

	brand, err := ctrl.service.GetBrandByID(id)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusOK, "Brand retrieved successfully", brand)
}
