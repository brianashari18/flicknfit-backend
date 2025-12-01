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

// ReviewController defines HTTP handlers for review operations
type ReviewController interface {
	GetProductReviews(c *fiber.Ctx) error
	CreateReview(c *fiber.Ctx) error
	UpdateReview(c *fiber.Ctx) error
	DeleteReview(c *fiber.Ctx) error
	GetUserReviews(c *fiber.Ctx) error
	GetProductReviewStats(c *fiber.Ctx) error
}

// reviewController implements ReviewController interface
type reviewController struct {
	service   services.ReviewService
	validator validator.Validate
}

// NewReviewController creates a new review controller
func NewReviewController(service services.ReviewService, validator *validator.Validate) ReviewController {
	return &reviewController{
		service:   service,
		validator: *validator,
	}
}

// GetProductReviews retrieves reviews for a specific product
// @Summary Get product reviews
// @Description Get paginated reviews for a specific product
// @Tags Reviews
// @Accept json
// @Produce json
// @Param productId path int true "Product ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response{data=dtos.ReviewListResponseDTO} "Reviews retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid product ID"
// @Failure 404 {object} utils.Response "Product not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /reviews/product/{productId} [get]
func (ctrl *reviewController) GetProductReviews(c *fiber.Ctx) error {
	productID, err := utils.GetUintParam(c, "productId")
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid product ID", nil)
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	result, err := ctrl.service.GetProductReviews(productID, page, limit)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusOK, "Reviews retrieved successfully", result)
}

// CreateReview creates a new review
// @Summary Create a product review
// @Description Create a new review for a product by authenticated user
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param review body dtos.CreateReviewDTO true "Review data"
// @Success 201 {object} utils.Response "Review created successfully"
// @Failure 400 {object} utils.Response "Invalid request body or validation failed"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Product not found"
// @Failure 409 {object} utils.Response "User already reviewed this product"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /reviews [post]
func (ctrl *reviewController) CreateReview(c *fiber.Ctx) error {
	var dto dtos.CreateReviewDTO
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

	err = ctrl.service.CreateReview(userID, &dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusCreated, "Review created successfully", nil)
}

// UpdateReview updates an existing review
func (ctrl *reviewController) UpdateReview(c *fiber.Ctx) error {
	reviewID, err := utils.GetUintParam(c, "reviewId")
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid review ID", nil)
	}

	var dto dtos.UpdateReviewDTO
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

	err = ctrl.service.UpdateReview(userID, reviewID, &dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusOK, "Review updated successfully", nil)
}

// DeleteReview deletes a review
func (ctrl *reviewController) DeleteReview(c *fiber.Ctx) error {
	reviewID, err := utils.GetUintParam(c, "reviewId")
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid review ID", nil)
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusUnauthorized, err.Error(), nil)
	}

	err = ctrl.service.DeleteReview(userID, reviewID)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusOK, "Review deleted successfully", nil)
}

// GetUserReviews retrieves all reviews by current user
func (ctrl *reviewController) GetUserReviews(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusUnauthorized, err.Error(), nil)
	}

	reviews, err := ctrl.service.GetUserReviews(userID)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	response := dtos.ToReviewResponseDTOs(reviews)
	return utils.SendResponse(c, http.StatusOK, "User reviews retrieved successfully", response)
}

// GetProductReviewStats retrieves review statistics for a product
func (ctrl *reviewController) GetProductReviewStats(c *fiber.Ctx) error {
	productID, err := utils.GetUintParam(c, "productId")
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid product ID", nil)
	}

	stats, err := ctrl.service.GetProductReviewStats(productID)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.SendResponse(c, http.StatusOK, "Review statistics retrieved successfully", stats)
}
