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
	GetAllProductsPublicWithFilter(c *fiber.Ctx) error
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
// @Summary Create a new product (Admin only)
// @Description Create a new product with full details
// @Tags Admin - Product Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body dtos.AdminProductCreateRequestDTO true "Product creation data"
// @Success 201 {object} utils.Response{data=dtos.ProductResponseDTO} "Product created successfully"
// @Failure 400 {object} utils.Response "Invalid request body or validation failed"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/products [post]
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
// @Summary Update product (Admin only)
// @Description Update an existing product's information
// @Tags Admin - Product Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param product body dtos.AdminProductUpdateRequestDTO true "Product update data"
// @Success 200 {object} utils.Response{data=dtos.AdminProductResponseDTO} "Product updated successfully"
// @Failure 400 {object} utils.Response "Invalid request body or product ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 404 {object} utils.Response "Product not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/products/{id} [put]
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
// @Summary Delete product (Admin only)
// @Description Permanently delete a product from the system
// @Tags Admin - Product Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} utils.Response "Product deleted successfully"
// @Failure 400 {object} utils.Response "Invalid product ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 404 {object} utils.Response "Product not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/products/{id} [delete]
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
// @Summary Get product by ID (Admin only)
// @Description Retrieve detailed product information for admin users
// @Tags Admin - Product Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} utils.Response{data=dtos.AdminProductResponseDTO} "Product retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid product ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 404 {object} utils.Response "Product not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/products/{id} [get]
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
// @Summary Get all products (Admin only)
// @Description Retrieve a list of all products for admin management
// @Tags Admin - Product Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]dtos.AdminProductResponseDTO} "Products retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/products [get]
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
// @Summary Get product by ID
// @Description Retrieve detailed product information for public viewing
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} utils.Response{data=dtos.ProductResponseDTO} "Product found"
// @Failure 400 {object} utils.Response "Invalid product ID"
// @Failure 404 {object} utils.Response "Product not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /products/{id} [get]
func (ctrl *productController) GetProductPublicByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid product ID", nil)
	}

	product, err := ctrl.productService.GetProductPublicByID(id)
	if err != nil {
		return utils.SendResponse(c, http.StatusNotFound, "Product not found", nil)
	}

	response := dtos.ToProductPublicResponseDTO(product)
	return utils.SendResponse(c, http.StatusOK, "Product found", response)
}

// GetAllProductsPublic handles fetching all products for public users.
// @Summary Get all products
// @Description Retrieve all available products for public viewing
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]dtos.ProductResponseDTO} "Products retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /products [get]
func (ctrl *productController) GetAllProductsPublic(c *fiber.Ctx) error {
	products, err := ctrl.productService.GetAllProductsPublic()
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to retrieve products", nil)
	}
	response := dtos.ToProductResponseDTOs(products)
	return utils.SendResponse(c, http.StatusOK, "Products retrieved successfully", response)
}

// GetReviewsByProductIDPublic handles fetching all reviews for a specific product for public users.
// @Summary Get product reviews
// @Description Retrieve all reviews for a specific product
// @Tags Products
// @Accept json
// @Produce json
// @Param productID path string true "Product ID"
// @Success 200 {object} utils.Response{data=[]dtos.AdminReviewResponseDTO} "Reviews retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid product ID"
// @Failure 404 {object} utils.Response "Reviews not found for this product"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /products/{productId}/reviews-list [get]
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
// @Summary Create product review
// @Description Create a new review for a specific product
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param productID path string true "Product ID"
// @Param review body dtos.ReviewCreateDTO true "Review data"
// @Success 201 {object} utils.Response "Review created successfully"
// @Failure 400 {object} utils.Response "Invalid request body or product ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /products/{productId}/review [post]
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
// @Summary Search products
// @Description Search for products by name, description, or other criteria
// @Tags Products
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {object} utils.Response{data=[]dtos.ProductResponseDTO} "Products found successfully"
// @Failure 400 {object} utils.Response "Search query 'q' is required"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /products/search [get]
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

// GetAllProductsPublicWithFilter handles fetching products with filters for public users.
// @Summary Get products with filters
// @Description Retrieve products with advanced filtering options (category, brand, price range, etc.)
// @Tags Products
// @Accept json
// @Produce json
// @Param category query string false "Product category"
// @Param brand query string false "Brand name"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param sort query string false "Sort by (price, name, created_at)"
// @Param order query string false "Order direction (asc, desc)"
// @Success 200 {object} utils.Response{data=[]dtos.ProductResponseDTO} "Products retrieved successfully with filters"
// @Failure 400 {object} utils.Response "Invalid filter parameters"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /products/filter [get]
func (ctrl *productController) GetAllProductsPublicWithFilter(c *fiber.Ctx) error {
	var filterParams dtos.ProductFilterRequestDTO
	if err := c.QueryParser(&filterParams); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid filter parameters: "+err.Error(), nil)
	}

	products, err := ctrl.productService.GetAllProductsPublicWithFilter(&filterParams)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to retrieve products with filters", nil)
	}

	response := dtos.ToProductResponseDTOs(products)
	return utils.SendResponse(c, http.StatusOK, "Products retrieved successfully with filters", response)
}
