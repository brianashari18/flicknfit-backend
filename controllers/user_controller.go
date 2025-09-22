package controllers

import (
	"flicknfit_backend/dtos"
	"flicknfit_backend/services"
	"flicknfit_backend/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// UserController defines the HTTP handlers for user-related operations using Fiber.
type UserController interface {
	// Admin Routes
	AdminCreateUser(c *fiber.Ctx) error
	AdminGetAllUsers(c *fiber.Ctx) error
	AdminGetUserByID(c *fiber.Ctx) error
	AdminUpdateUser(c *fiber.Ctx) error
	AdminDeleteUser(c *fiber.Ctx) error

	// Public Routes
	RegisterUser(c *fiber.Ctx) error
	LoginUser(c *fiber.Ctx) error
	ForgotPassword(c *fiber.Ctx) error
	VerifyOTP(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
	LogoutUser(c *fiber.Ctx) error
	GetUserByAccessToken(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	EditProfile(c *fiber.Ctx) error
}

// userController is the implementation of UserController.
type userController struct {
	service   services.UserService
	validator validator.Validate
}

// NewUserController creates and returns a new instance of UserController.
func NewUserController(service services.UserService, validator *validator.Validate) UserController {
	return &userController{service: service, validator: *validator}
}

// AdminCreateUser handles creating a new user by an admin.
func (ctrl *userController) AdminCreateUser(c *fiber.Ctx) error {
	var dto dtos.UserAdminCreateRequestDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}

	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
	}

	if err := ctrl.service.AdminCreateUser(&dto); err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to create user", nil)
	}
	return utils.SendResponse(c, http.StatusCreated, "User created successfully", nil)
}

// AdminGetAllUsers retrieves all users for an admin.
func (ctrl *userController) AdminGetAllUsers(c *fiber.Ctx) error {
	users, err := ctrl.service.GetAllUsers()
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to retrieve users", nil)
	}
	response := dtos.ToUserAdminResponseDTOs(users)
	return utils.SendResponse(c, http.StatusOK, "Users retrieved successfully", response)
}

// AdminGetUserByID retrieves a specific user for an admin.
func (ctrl *userController) AdminGetUserByID(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}
	user, err := ctrl.service.GetUserByID(userID)
	if err != nil {
		return utils.SendResponse(c, http.StatusNotFound, "User not found", nil)
	}
	response := dtos.ToUserAdminResponseDTO(user)
	return utils.SendResponse(c, http.StatusOK, "User retrieved successfully", response)
}

// AdminUpdateUser handles updating an existing user by an admin.
func (ctrl *userController) AdminUpdateUser(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}
	var dto dtos.UserAdminUpdateRequestDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}
	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
	}
	if err := ctrl.service.AdminUpdateUser(userID, &dto); err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to update user", nil)
	}
	return utils.SendResponse(c, http.StatusOK, "User updated successfully", nil)
}

// AdminDeleteUser handles deleting a user by an admin.
func (ctrl *userController) AdminDeleteUser(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}
	if err := ctrl.service.AdminDeleteUser(userID); err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to delete user", nil)
	}
	return utils.SendResponse(c, http.StatusOK, "User deleted successfully", nil)
}

// RegisterUser handles a new user registration.
func (ctrl *userController) RegisterUser(c *fiber.Ctx) error {
	var dto dtos.UserRegisterRequestDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}
	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
	}
	if err := ctrl.service.RegisterUser(&dto); err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to register user", nil)
	}
	return utils.SendResponse(c, http.StatusCreated, "User registered successfully", nil)
}

// LoginUser handles user authentication.
func (ctrl *userController) LoginUser(c *fiber.Ctx) error {
	var dto dtos.UserLoginRequestDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}
	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
	}
	tokens, err := ctrl.service.LoginUser(&dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusUnauthorized, err.Error(), nil)
	}
	response := dtos.ToUserLoginResponseDTO(tokens)
	return utils.SendResponse(c, http.StatusOK, "Login successful", response)
}

// ForgotPassword initiates the password reset process.
func (ctrl *userController) ForgotPassword(c *fiber.Ctx) error {
	var dto dtos.ForgotPasswordDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}
	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
	}
	if err := ctrl.service.ForgotPassword(&dto); err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}
	return utils.SendResponse(c, http.StatusOK, "OTP sent to email", nil)
}

// VerifyOTP handles OTP verification.
func (ctrl *userController) VerifyOTP(c *fiber.Ctx) error {
	var dto dtos.VerifyOTPDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}
	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
	}
	tokens, err := ctrl.service.VerifyOTP(&dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, err.Error(), nil)
	}
	response := dtos.ToUserResetTokenResponseDTO(tokens)
	return utils.SendResponse(c, http.StatusOK, "OTP verified successfully", response)
}

// ResetPassword allows the user to set a new password.
func (ctrl *userController) ResetPassword(c *fiber.Ctx) error {
	var dto dtos.UserResetPasswordRequestDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}
	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
	}
	if err := ctrl.service.ResetPassword(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, err.Error(), nil)
	}
	return utils.SendResponse(c, http.StatusOK, "Password reset successfully", nil)
}

// LogoutUser logs out the current user by invalidating their refresh token.
func (ctrl *userController) LogoutUser(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}
	if err := ctrl.service.LogoutUser(userID); err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to logout", nil)
	}
	return utils.SendResponse(c, http.StatusOK, "Logout successful", nil)
}

// GetUserByAccessToken retrieves a user based on a valid access token.
func (ctrl *userController) GetUserByAccessToken(c *fiber.Ctx) error {
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}
	user, err := ctrl.service.GetUserByID(userID)
	if err != nil {
		return utils.SendResponse(c, http.StatusNotFound, "User not found", nil)
	}
	response := dtos.ToUserResponseDTO(user)

	return utils.SendResponse(c, http.StatusOK, "User retrieved successfully", response)
}

// RefreshToken handles generating a new access token using a refresh token.
func (ctrl *userController) RefreshToken(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := utils.StrictBodyParser(c, &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	tokens, err := ctrl.service.RefreshToken(body.RefreshToken)
	if err != nil {
		return utils.SendResponse(c, http.StatusUnauthorized, err.Error(), nil)
	}

	response := dtos.ToUserLoginResponseDTO(tokens)
	return utils.SendResponse(c, http.StatusOK, "Token refreshed successfully", response)
}

func (ctrl *userController) EditProfile(c *fiber.Ctx) error {
	var dto dtos.UserEditProfileRequestDTO
	if err := utils.StrictBodyParser(c, &dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error(), nil)
	}
	if err := ctrl.validator.Struct(&dto); err != nil {
		return utils.SendResponse(c, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	updatedUsers, err := ctrl.service.EditProfile(userID, &dto)
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	response := dtos.ToUserEditProfileResponseDTO(*updatedUsers)
	return utils.SendResponse(c, http.StatusOK, "profile updated", response)

}
