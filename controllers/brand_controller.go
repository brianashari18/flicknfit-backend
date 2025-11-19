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
// @Summary Create a new brand (Admin only)
// @Description Create a new brand in the system
// @Tags Admin - Brand Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param brand body dtos.BrandCreateRequestDTO true "Brand creation data"
// @Success 201 {object} utils.Response{data=dtos.BrandResponseDTO} "Brand created successfully"
// @Failure 400 {object} utils.Response "Invalid request body or validation failed"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/brands [post]
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
// @Summary Update brand (Admin only)
// @Description Update an existing brand's information
// @Tags Admin - Brand Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Brand ID"
// @Param brand body dtos.BrandUpdateRequestDTO true "Brand update data"
// @Success 200 {object} utils.Response{data=dtos.BrandResponseDTO} "Brand updated successfully"
// @Failure 400 {object} utils.Response "Invalid request body or brand ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 404 {object} utils.Response "Brand not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/brands/{id} [put]
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
// @Summary Delete brand (Admin only)
// @Description Permanently delete a brand from the system
// @Tags Admin - Brand Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Brand ID"
// @Success 200 {object} utils.Response "Brand deleted successfully"
// @Failure 400 {object} utils.Response "Invalid brand ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 404 {object} utils.Response "Brand not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/brands/{id} [delete]
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
// @Summary Get all brands
// @Description Retrieve a list of all available brands
// @Tags Brands
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]dtos.BrandResponseDTO} "Brands retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /brands [get]
func (ctrl *brandController) GetAllBrands(c *fiber.Ctx) error {
	brands, err := ctrl.service.GetAllBrands()
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusOK, "Brands retrieved successfully", brands)
}

// GetBrandByID retrieves a single brand by its ID for public display.
// @Summary Get brand by ID
// @Description Retrieve detailed information about a specific brand
// @Tags Brands
// @Accept json
// @Produce json
// @Param id path string true "Brand ID"
// @Success 200 {object} utils.Response{data=dtos.BrandResponseDTO} "Brand retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid brand ID"
// @Failure 404 {object} utils.Response "Brand not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /brands/{id} [get]
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
