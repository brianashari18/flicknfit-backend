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
// @Summary Create a new user (Admin only)
// @Description Create a new user account with admin privileges
// @Tags Admin - User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body dtos.UserAdminCreateRequestDTO true "User creation data"
// @Success 201 {object} utils.Response{data=dtos.UserResponseDTO} "User created successfully"
// @Failure 400 {object} utils.Response "Invalid request body"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 409 {object} utils.Response "User already exists"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/users [post]
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
// @Summary Get all users (Admin only)
// @Description Retrieve a list of all registered users
// @Tags Admin - User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]dtos.UserAdminResponseDTO} "Users retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/users [get]
func (ctrl *userController) AdminGetAllUsers(c *fiber.Ctx) error {
	users, err := ctrl.service.GetAllUsers()
	if err != nil {
		return utils.SendResponse(c, http.StatusInternalServerError, "Failed to retrieve users", nil)
	}
	response := dtos.ToUserAdminResponseDTOs(users)
	return utils.SendResponse(c, http.StatusOK, "Users retrieved successfully", response)
}

// AdminGetUserByID retrieves a specific user for an admin.
// @Summary Get user by ID (Admin only)
// @Description Retrieve detailed information about a specific user
// @Tags Admin - User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} utils.Response{data=dtos.UserAdminResponseDTO} "User retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/users/{id} [get]
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
// @Summary Update user (Admin only)
// @Description Update user information by admin
// @Tags Admin - User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param user body dtos.UserAdminUpdateRequestDTO true "User update data"
// @Success 200 {object} utils.Response "User updated successfully"
// @Failure 400 {object} utils.Response "Invalid request body or validation failed"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/users/{id} [put]
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
// @Summary Delete user (Admin only)
// @Description Permanently delete a user account
// @Tags Admin - User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} utils.Response "User deleted successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden - Admin access required"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/users/{id} [delete]
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
// @Summary Register a new user
// @Description Create a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body dtos.UserRegisterRequestDTO true "User registration data"
// @Success 201 {object} utils.Response "User registered successfully"
// @Failure 400 {object} utils.Response "Invalid request body or validation failed"
// @Failure 409 {object} utils.Response "User already exists"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /auth/register [post]
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
// @Summary User login
// @Description Authenticate user and return JWT tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body dtos.UserLoginRequestDTO true "Login credentials"
// @Success 200 {object} utils.Response{data=dtos.UserLoginResponseDTO} "Login successful"
// @Failure 400 {object} utils.Response "Invalid request body or validation failed"
// @Failure 401 {object} utils.Response "Invalid credentials"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /auth/login [post]
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
// @Summary Forgot password
// @Description Initiate password reset process by sending OTP to email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param email body dtos.ForgotPasswordDTO true "Email for password reset"
// @Success 200 {object} utils.Response "OTP sent to email"
// @Failure 400 {object} utils.Response "Invalid request body or validation failed"
// @Failure 404 {object} utils.Response "Email not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /auth/forgot-password [post]
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
// @Summary Verify OTP
// @Description Verify OTP code for password reset
// @Tags Authentication
// @Accept json
// @Produce json
// @Param otp body dtos.VerifyOTPDTO true "OTP verification data"
// @Success 200 {object} utils.Response{data=dtos.UserResetTokenResponseDTO} "OTP verified successfully"
// @Failure 400 {object} utils.Response "Invalid OTP or expired"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /auth/verify-otp [post]
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
// @Summary Reset password
// @Description Set a new password using reset token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param password body dtos.UserResetPasswordRequestDTO true "New password data"
// @Success 200 {object} utils.Response "Password reset successfully"
// @Failure 400 {object} utils.Response "Invalid request body or token expired"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /auth/reset-password [post]
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
// @Summary User logout
// @Description Logout user by invalidating refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response "Logout successful"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /auth/logout [post]
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
// @Summary Get current user profile
// @Description Get current authenticated user's profile information
// @Tags User Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=dtos.UserResponseDTO} "User retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /user/profile [get]
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
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token body object{refresh_token=string} true "Refresh token"
// @Success 200 {object} utils.Response{data=dtos.UserLoginResponseDTO} "Token refreshed successfully"
// @Failure 400 {object} utils.Response "Invalid request body"
// @Failure 401 {object} utils.Response "Invalid or expired refresh token"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /auth/refresh [post]
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

// EditProfile allows users to update their profile information.
// @Summary Edit user profile
// @Description Update current user's profile information
// @Tags User Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body dtos.UserEditProfileRequestDTO true "Profile update data"
// @Success 200 {object} utils.Response{data=dtos.UserEditProfileResponseDTO} "Profile updated successfully"
// @Failure 400 {object} utils.Response "Invalid request body or validation failed"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /user/profile [put]
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
