package controllers

import (
	"context"
	"flicknfit_backend/dtos"
	"flicknfit_backend/services"
	"flicknfit_backend/utils"

	"github.com/gofiber/fiber/v2"
)

// OAuthController handles OAuth authentication endpoints
type OAuthController struct {
	userService     services.UserService
	firebaseService *services.FirebaseService
}

// NewOAuthController creates a new OAuth controller
func NewOAuthController(userService services.UserService, firebaseService *services.FirebaseService) *OAuthController {
	return &OAuthController{
		userService:     userService,
		firebaseService: firebaseService,
	}
}

// GoogleLogin handles Google OAuth login
// @Summary Google OAuth login
// @Description Authenticate user with Google account via Firebase
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dtos.OAuthLoginRequestDTO true "Google login credentials"
// @Success 200 {object} utils.Response{data=dtos.UserLoginResponseDTO}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/google [post]
func (c *OAuthController) GoogleLogin(ctx *fiber.Ctx) error {
	logger := utils.GetLogger()
	validator := utils.NewValidator()

	var dto dtos.OAuthLoginRequestDTO
	if err := ctx.BodyParser(&dto); err != nil {
		logger.Error("Failed to parse request body: ", err)
		return utils.SendError(ctx, fiber.StatusBadRequest, "Invalid request payload")
	}

	// Validate request
	if err := validator.Struct(&dto); err != nil {
		logger.Error("Validation failed: ", err)
		return utils.SendError(ctx, fiber.StatusBadRequest, err.Error())
	}

	// Verify Firebase token
	token, err := c.firebaseService.VerifyIDToken(context.Background(), dto.FirebaseToken)
	if err != nil {
		logger.Error("Firebase token verification failed: ", err)
		return utils.SendError(ctx, fiber.StatusUnauthorized, "Invalid Firebase token")
	}

	// Ensure provider is Google
	dto.AuthProvider = "google"
	dto.AuthProviderID = token.UID

	// Get email from token if not provided
	if dto.Email == "" {
		if email, ok := token.Claims["email"].(string); ok {
			dto.Email = email
		} else {
			return utils.SendError(ctx, fiber.StatusBadRequest, "Email not found in token")
		}
	}

	// Perform OAuth login/registration
	tokens, err := c.userService.OAuthLogin(&dto, token.UID)
	if err != nil {
		logger.Error("OAuth login failed: ", err)
		return utils.SendError(ctx, fiber.StatusInternalServerError, "Failed to authenticate user")
	}

	response := dtos.UserLoginResponseDTO{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	logger.Info("Google OAuth login successful for email: ", dto.Email)
	return utils.SendSuccess(ctx, "Google login successful", response)
}

// FacebookLogin handles Facebook OAuth login
// @Summary Facebook OAuth login
// @Description Authenticate user with Facebook account via Firebase
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dtos.OAuthLoginRequestDTO true "Facebook login credentials"
// @Success 200 {object} utils.Response{data=dtos.UserLoginResponseDTO}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/facebook [post]
func (c *OAuthController) FacebookLogin(ctx *fiber.Ctx) error {
	logger := utils.GetLogger()
	validator := utils.NewValidator()

	var dto dtos.OAuthLoginRequestDTO
	if err := ctx.BodyParser(&dto); err != nil {
		logger.Error("Failed to parse request body: ", err)
		return utils.SendError(ctx, fiber.StatusBadRequest, "Invalid request payload")
	}

	// Validate request
	if err := validator.Struct(&dto); err != nil {
		logger.Error("Validation failed: ", err)
		return utils.SendError(ctx, fiber.StatusBadRequest, err.Error())
	}

	// Verify Firebase token
	token, err := c.firebaseService.VerifyIDToken(context.Background(), dto.FirebaseToken)
	if err != nil {
		logger.Error("Firebase token verification failed: ", err)
		return utils.SendError(ctx, fiber.StatusUnauthorized, "Invalid Firebase token")
	}

	// Ensure provider is Facebook
	dto.AuthProvider = "facebook"
	dto.AuthProviderID = token.UID

	// Get email from token if not provided
	if dto.Email == "" {
		if email, ok := token.Claims["email"].(string); ok {
			dto.Email = email
		} else {
			return utils.SendError(ctx, fiber.StatusBadRequest, "Email not found in token")
		}
	}

	// Perform OAuth login/registration
	tokens, err := c.userService.OAuthLogin(&dto, token.UID)
	if err != nil {
		logger.Error("OAuth login failed: ", err)
		return utils.SendError(ctx, fiber.StatusInternalServerError, "Failed to authenticate user")
	}

	response := dtos.UserLoginResponseDTO{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	logger.Info("Facebook OAuth login successful for email: ", dto.Email)
	return utils.SendSuccess(ctx, "Facebook login successful", response)
}
