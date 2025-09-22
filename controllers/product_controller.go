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
	GetReviewsByProductIDPublic(c *fiber.Ctx) error // Endpoint baru
	CreateReview(c *fiber.Ctx) error                // Endpoint baru
	SearchProductsPublic(c *fiber.Ctx) error
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
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}

	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
	}

	product, err := ctrl.productService.AdminCreateProduct(&dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to create product: "+err.Error(), nil)
	}

	response := dtos.ToAdminProductResponseDTO(product)
	return utils.SendResponse(c, http.StatusCreated, "Product created successfully", response)
}

// AdminUpdateProduct handles updating an existing product by an admin.
func (ctrl *productController) AdminUpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid product ID", nil)
	}

	var dto dtos.AdminProductUpdateRequestDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}

	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
	}

	product, err := ctrl.productService.AdminUpdateProduct(id, &dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to update product: "+err.Error(), nil)
	}

	response := dtos.ToAdminProductResponseDTO(product)
	return utils.SendResponse(c, http.StatusOK, "Product updated successfully", response)
}

// AdminDeleteProduct handles deleting a product by an admin.
func (ctrl *productController) AdminDeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid product ID", nil)
	}

	if err := ctrl.productService.AdminDeleteProduct(id); err != nil {
		return utils.SendResponse(c, http.StatusNotFound, "Failed to delete product: "+err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusOK, "Product deleted successfully", nil)
}

// AdminGetProductByID handles fetching a single product by ID for admin users.
func (ctrl *productController) AdminGetProductByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid product ID", nil)
	}

	product, err := ctrl.productService.AdminGetProductByID(id)
	if err != nil {
		return utils.SendResponse(c, http.StatusNotFound, "Product not found", nil)
	}

	response := dtos.ToAdminProductResponseDTO(product)
	return utils.SendResponse(c, http.StatusOK, "Product retrieved successfully", response)
}

// AdminGetAllProducts handles fetching all products for admin users.
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
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
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

// GetReviewsByProductIDPublic handles fetching all reviews for a specific product for public users.
func (ctrl *productController) GetReviewsByProductIDPublic(c *fiber.Ctx) error {
	productID, err := strconv.ParseUint(c.Params("productID"), 10, 64)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid product ID", nil)
	}

	reviews, err := ctrl.productService.AdminGetAllReviewsByProductID(productID)
	if err != nil {
		return utils.SendResponse(c, http.StatusNotFound, "Reviews not found for this product", nil)
	}

	response := dtos.ToAdminReviewResponseDTOs(reviews)
	return utils.SendResponse(c, http.StatusOK, "Reviews retrieved successfully", response)
}

// CreateReview handles the creation of a new review by an authenticated user.
func (ctrl *productController) CreateReview(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to get user ID from token", nil)
	}

	productID, err := strconv.ParseUint(c.Params("productID"), 10, 64)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid product ID", nil)
	}

	var dto dtos.ReviewCreateDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}

	// Override DTO productID and userID with values from URL and JWT
	dto.ProductID = productID
	dto.UserID = userID

	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
	}

	review, err := ctrl.productService.CreateReviewPublic(&dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to create review: "+err.Error(), nil)
	}

	response := dtos.ToReviewResponseDTO(review)
	return utils.SendResponse(c, http.StatusCreated, "Review created successfully", response)
}

// SearchProductsPublic handles product searches for public users.
func (ctrl *productController) SearchProductsPublic(c *fiber.Ctx) error {
	// Mengambil parameter kueri 'q' dari URL.
	query := c.Query("q")
	if query == "" {
		return utils.SendResponse(c, http.StatusBadRequest, "Search query 'q' is required", nil)
	}

	// Memanggil service untuk melakukan pencarian.
	products, err := ctrl.productService.SearchProductsPublic(query)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to search products", nil)
	}

	response := dtos.ToProductResponseDTOs(products)
	return utils.SendResponse(c, http.StatusOK, "Products found", response)
}
