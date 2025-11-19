package integration

import (
	"bytes"
	"encoding/json"
	"flicknfit_backend/config"
	"flicknfit_backend/container"
	"flicknfit_backend/database"
	"flicknfit_backend/dtos"
	"flicknfit_backend/routes"
	"flicknfit_backend/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type FavoriteIntegrationTestSuite struct {
	suite.Suite
	app       *fiber.App
	db        *gorm.DB
	container *container.Container
	userToken string
}

func (suite *FavoriteIntegrationTestSuite) SetupSuite() {
	// Setup test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	suite.db = db
	// Auto-migrate all models
	database.Migrate(db, utils.NewLogger())

	// Setup test config
	cfg := &config.Config{
		JwtSecretKey: "test-secret",
		AppPort:      "8080",
		DBHost:       "localhost",
		DBPort:       "3306",
		DBUser:       "test",
		DBPassword:   "test",
		DBName:       "flicknfit_test",
	}

	// Create container
	container, err := container.NewContainer(db, cfg)
	suite.Require().NoError(err)
	suite.container = container

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	// Setup routes
	routes.SetupRoutes(app, db, container)
	suite.app = app

	// Create test user and get token
	suite.createTestUserAndGetToken()
}

func (suite *FavoriteIntegrationTestSuite) createTestUserAndGetToken() {
	// Create test user
	registerDTO := dtos.UserRegisterRequestDTO{
		Email:       "test@example.com",
		Username:    "testuser",
		Password:    "TestPassword123!",
		PhoneNumber: "+1234567890",
	}

	registerBody, _ := json.Marshal(registerDTO)
	req := httptest.NewRequest("POST", "/api/v1/users/register", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)
	suite.Require().Equal(http.StatusCreated, resp.StatusCode)

	// Login to get token
	loginDTO := dtos.UserLoginRequestDTO{
		Email:    "test@example.com",
		Password: "TestPassword123!",
	}

	loginBody, _ := json.Marshal(loginDTO)
	req = httptest.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.app.Test(req)
	suite.Require().NoError(err)
	suite.Require().Equal(http.StatusOK, resp.StatusCode)

	var loginResponse struct {
		Success bool `json:"success"`
		Data    struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	suite.Require().NoError(err)
	suite.userToken = loginResponse.Data.AccessToken
}

func (suite *FavoriteIntegrationTestSuite) TearDownSuite() {
	// Clean up database
	suite.db.Exec("DELETE FROM favorites")
	suite.db.Exec("DELETE FROM users")
	suite.db.Exec("DELETE FROM products")
	suite.db.Exec("DELETE FROM brands")
}

func (suite *FavoriteIntegrationTestSuite) TestGetUserFavorites() {
	req := httptest.NewRequest("GET", "/api/v1/favorites", nil)
	req.Header.Set("Authorization", "Bearer "+suite.userToken)

	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response utils.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err)

	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "Favorites retrieved successfully", response.Message)
}

func (suite *FavoriteIntegrationTestSuite) TestGetUserFavoritesUnauthorized() {
	req := httptest.NewRequest("GET", "/api/v1/favorites", nil)
	// No authorization header

	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *FavoriteIntegrationTestSuite) TestToggleFavoriteInvalidJSON() {
	invalidJSON := `{"invalid": json}`
	req := httptest.NewRequest("POST", "/api/v1/favorites/1", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.userToken)

	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *FavoriteIntegrationTestSuite) TestToggleFavoriteValidationError() {
	// Missing required fields
	dto := map[string]interface{}{
		"product_item_id": "", // Invalid - should be uint64
	}

	body, _ := json.Marshal(dto)
	req := httptest.NewRequest("POST", "/api/v1/favorites/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.userToken)

	resp, err := suite.app.Test(req)
	suite.Require().NoError(err)

	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	var response utils.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.Require().NoError(err)

	assert.False(suite.T(), response.Success)
	assert.Contains(suite.T(), response.Message, "Validation failed")
}

func TestFavoriteIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(FavoriteIntegrationTestSuite))
}
