package controllers

import (
	"flicknfit_backend/dtos"
	"flicknfit_backend/services"
	"flicknfit_backend/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ProductController defines the HTTP handlers for product-related operations using Fiber.
type ProductController interface {
	// Admin Routes
	AdminCreateProduct(c *fiber.Ctx) error
	AdminUpdateProduct(c *fiber.Ctx) error
	AdminDeleteProduct(c *fiber.Ctx) error
	AdminGetProductByID(c *fiber.Ctx) error
	AdminGetAllProducts(c *fiber.Ctx) error

	// Public Routes
	GetProductPublicByID(c *fiber.Ctx) error
	GetAllProductsPublic(c *fiber.Ctx) error
}

// productController is the implementation of ProductController.
type productController struct {
	productService services.ProductService
	validator      *validator.Validate
}

// NewProductController creates and returns a new instance of ProductController.
func NewProductController(productService services.ProductService, validator *validator.Validate) ProductController {
	return &productController{
		productService: productService,
		validator:      validator,
	}
}

// AdminCreateProduct handles the creation of a new product by an admin.
func (ctrl *productController) AdminCreateProduct(c *fiber.Ctx) error {
	var dto dtos.AdminProductCreateRequestDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body", nil)
	}

	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed", err)
	}

	product, err := ctrl.productService.AdminCreateProduct(&dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to create product", nil)
	}

	response := dtos.ToAdminProductResponseDTO(product)
	return utils.SendResponse(c, http.StatusCreated, "Product created successfully", response)
}

// AdminUpdateProduct handles the update of a product by an admin.
func (ctrl *productController) AdminUpdateProduct(c *fiber.Ctx) error {
	id, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid product ID", nil)
	}

	var dto dtos.AdminProductUpdateRequestDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body", nil)
	}

	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed", err)
	}

	updatedProduct, err := ctrl.productService.AdminUpdateProduct(id, &dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusNotFound, "Product not found or failed to update", nil)
	}

	response := dtos.ToAdminProductResponseDTO(updatedProduct)
	return utils.SendResponse(c, http.StatusOK, "Product updated successfully", response)
}

// AdminDeleteProduct handles the deletion of a product by an admin.
func (ctrl *productController) AdminDeleteProduct(c *fiber.Ctx) error {
	id, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid product ID", nil)
	}

	if err := ctrl.productService.AdminDeleteProduct(id); err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to delete product", nil)
	}

	return utils.SendResponse(c, http.StatusOK, "Product deleted successfully", nil)
}

// AdminGetProductByID handles fetching a single product by ID for admin.
func (ctrl *productController) AdminGetProductByID(c *fiber.Ctx) error {
	id, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid product ID", nil)
	}

	product, err := ctrl.productService.AdminGetProductByID(id)
	if err != nil {
		return utils.SendResponse(c, http.StatusNotFound, "Product not found", nil)
	}

	response := dtos.ToAdminProductResponseDTO(product)
	return utils.SendResponse(c, http.StatusOK, "Product found", response)
}

// AdminGetAllProducts handles fetching all products for admin.
func (ctrl *productController) AdminGetAllProducts(c *fiber.Ctx) error {
	products, err := ctrl.productService.AdminGetAllProducts()
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to retrieve products", nil)
	}

	response := make([]dtos.AdminProductResponseDTO, len(products))
	for i, p := range products {
		response[i] = dtos.ToAdminProductResponseDTO(p)
	}

	return utils.SendResponse(c, http.StatusOK, "Products retrieved successfully", response)
}

// GetProductPublicByID handles fetching a single product by ID for public users.
func (ctrl *productController) GetProductPublicByID(c *fiber.Ctx) error {
	id, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid product ID", nil)
	}

	product, err := ctrl.productService.GetProductPublicByID(id)
	if err != nil {
		return utils.SendResponse(c, http.StatusNotFound, "Product not found", nil)
	}

	response := dtos.ToProductResponseDTO(product)
	return utils.SendResponse(c, http.StatusOK, "Product found", response)
}

// GetAllProductsPublic handles fetching all products for public users.
func (ctrl *productController) GetAllProductsPublic(c *fiber.Ctx) error {
	products, err := ctrl.productService.GetAllProductsPublic()
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to retrieve products", nil)
	}

	response := dtos.ToProductResponseDTOs(products)
	return utils.SendResponse(c, http.StatusOK, "Products retrieved successfully", response)
}
