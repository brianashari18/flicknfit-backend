package routes

import (
	"flicknfit_backend/container"
	"flicknfit_backend/middlewares"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"gorm.io/gorm"
)

// SetupRoutes sets up all API routes for the application with clean structure.
func SetupRoutes(app *fiber.App, db *gorm.DB, container *container.Container) {
	// Setup common security and performance middlewares
	middlewares.SetupCommonMiddlewares(app)

	// Add global error handling middlewares
	app.Use(middlewares.ErrorHandler())
	app.Use(middlewares.RecoverHandler())

	// Setup Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Setup API routes
	setupAPIRoutes(app, container)
}

// setupAPIRoutes configures all API routes
func setupAPIRoutes(app *fiber.App, container *container.Container) {
	api := app.Group("/api/v1")

	// Setup API-specific middlewares
	middlewares.SetupAPIMiddlewares(api)
	// Setup user routes
	setupUserRoutes(api, container)

	// Setup product routes
	setupProductRoutes(api, container)

	// Setup brand routes
	setupBrandRoutes(api, container)

	// Setup shopping cart routes
	setupShoppingCartRoutes(api, container)
	// Setup new feature routes
	setupFavoriteRoutes(api, container)
	setupReviewRoutes(api, container)
	setupWardrobeRoutes(api, container)

	// Setup AI prediction routes
	setupAIRoutes(api, container)
}

// setupUserRoutes configures all user-related routes
func setupUserRoutes(api fiber.Router, c *container.Container) {
	// Public user routes with auth rate limiting
	userRoutes := api.Group("/users")
	middlewares.SetupAuthMiddlewares(userRoutes)

	userRoutes.Post("/register", c.Controllers.User.RegisterUser)
	userRoutes.Post("/login", c.Controllers.User.LoginUser)
	userRoutes.Post("/forgot-password", c.Controllers.User.ForgotPassword)
	userRoutes.Post("/verify-otp", c.Controllers.User.VerifyOTP)
	userRoutes.Post("/reset-password", c.Controllers.User.ResetPassword)
	userRoutes.Post("/refresh-token", c.Controllers.User.RefreshToken)

	// Private user routes (requires authentication)
	privateRoutes := api.Group("/users")
	privateRoutes.Use(middlewares.AuthMiddleware())
	privateRoutes.Post("/logout", c.Controllers.User.LogoutUser)
	privateRoutes.Get("/me", c.Controllers.User.GetUserByAccessToken)
	privateRoutes.Patch("/edit-profile", c.Controllers.User.EditProfile)

	// Admin user routes (requires authentication + admin role)
	userAdminRoutes := api.Group("/admin/users")
	userAdminRoutes.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	userAdminRoutes.Post("/", c.Controllers.User.AdminCreateUser)
	userAdminRoutes.Get("/", c.Controllers.User.AdminGetAllUsers)
	userAdminRoutes.Get("/:id", c.Controllers.User.AdminGetUserByID)
	userAdminRoutes.Put("/:id", c.Controllers.User.AdminUpdateUser)
	userAdminRoutes.Delete("/:id", c.Controllers.User.AdminDeleteUser)
}

// setupProductRoutes configures all product-related routes
func setupProductRoutes(api fiber.Router, c *container.Container) {
	// Public product routes
	productRoutes := api.Group("/products")
	productRoutes.Get("/", c.Controllers.Product.GetAllProductsPublic)
	productRoutes.Get("/:id", c.Controllers.Product.GetProductPublicByID)
	productRoutes.Get("/search", c.Controllers.Product.SearchProductsPublic)
	productRoutes.Get("/filter", c.Controllers.Product.GetAllProductsPublicWithFilter)

	// Admin product routes
	productAdminRoutes := api.Group("/admin/products")
	productAdminRoutes.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	productAdminRoutes.Post("/", c.Controllers.Product.AdminCreateProduct)
	productAdminRoutes.Get("/", c.Controllers.Product.AdminGetAllProducts)
	productAdminRoutes.Get("/:id", c.Controllers.Product.AdminGetProductByID)
	productAdminRoutes.Put("/:id", c.Controllers.Product.AdminUpdateProduct)
	productAdminRoutes.Delete("/:id", c.Controllers.Product.AdminDeleteProduct)
}

// setupBrandRoutes configures all brand-related routes
func setupBrandRoutes(api fiber.Router, c *container.Container) {
	// Public brand routes
	brandRoutes := api.Group("/brands")
	brandRoutes.Get("/", c.Controllers.Brand.GetAllBrands)
	brandRoutes.Get("/:id", c.Controllers.Brand.GetBrandByID)

	// Admin brand routes
	brandAdminRoutes := api.Group("/admin/brands")
	brandAdminRoutes.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	brandAdminRoutes.Post("/", c.Controllers.Brand.AdminCreateBrand)
	brandAdminRoutes.Put("/:id", c.Controllers.Brand.AdminUpdateBrand)
	brandAdminRoutes.Delete("/:id", c.Controllers.Brand.AdminDeleteBrand)
}

// setupShoppingCartRoutes configures all shopping cart routes
func setupShoppingCartRoutes(api fiber.Router, c *container.Container) {
	// All cart routes require authentication
	shoppingCartRoutes := api.Group("/cart")
	shoppingCartRoutes.Use(middlewares.AuthMiddleware())
	shoppingCartRoutes.Get("/", c.Controllers.ShoppingCart.GetUserCart)
	shoppingCartRoutes.Post("/", c.Controllers.ShoppingCart.AddProductItemToCart)
	shoppingCartRoutes.Put("/:itemId", c.Controllers.ShoppingCart.UpdateProductItemInCart)
	shoppingCartRoutes.Delete("/:itemId", c.Controllers.ShoppingCart.RemoveProductItemFromCart)
}

// setupFavoriteRoutes configures all favorite-related routes
func setupFavoriteRoutes(api fiber.Router, c *container.Container) {
	// All favorite routes require authentication
	favoriteRoutes := api.Group("/favorites")
	favoriteRoutes.Use(middlewares.AuthMiddleware())
	favoriteRoutes.Get("/", c.Controllers.Favorite.GetUserFavorites)
	favoriteRoutes.Post("/:productId", c.Controllers.Favorite.ToggleFavorite)
	favoriteRoutes.Delete("/:productId", c.Controllers.Favorite.RemoveFavorite)
}

// setupReviewRoutes configures all review-related routes
func setupReviewRoutes(api fiber.Router, c *container.Container) {
	// Public review routes
	reviewRoutes := api.Group("/reviews")
	reviewRoutes.Get("/product/:productId", c.Controllers.Review.GetProductReviews)
	reviewRoutes.Get("/product/:productId/stats", c.Controllers.Review.GetProductReviewStats)

	// Authenticated review routes
	authReviewRoutes := api.Group("/reviews")
	authReviewRoutes.Use(middlewares.AuthMiddleware())
	authReviewRoutes.Post("/", c.Controllers.Review.CreateReview)
	authReviewRoutes.Put("/:reviewId", c.Controllers.Review.UpdateReview)
	authReviewRoutes.Delete("/:reviewId", c.Controllers.Review.DeleteReview)
	authReviewRoutes.Get("/user", c.Controllers.Review.GetUserReviews)
}

// setupWardrobeRoutes configures all wardrobe-related routes
func setupWardrobeRoutes(api fiber.Router, c *container.Container) {
	// All wardrobe routes require authentication
	wardrobeRoutes := api.Group("/wardrobe")
	wardrobeRoutes.Use(middlewares.AuthMiddleware())
	wardrobeRoutes.Get("/", c.Controllers.Wardrobe.GetUserWardrobe)
	wardrobeRoutes.Post("/", c.Controllers.Wardrobe.CreateWardrobeItem)
	wardrobeRoutes.Put("/:itemId", c.Controllers.Wardrobe.UpdateWardrobeItem)
	wardrobeRoutes.Delete("/:itemId", c.Controllers.Wardrobe.DeleteWardrobeItem)
	wardrobeRoutes.Get("/categories", c.Controllers.Wardrobe.GetWardrobeCategories)
}

// setupAIRoutes configures all AI prediction-related routes
func setupAIRoutes(api fiber.Router, c *container.Container) {
	// All AI prediction routes require authentication
	aiRoutes := api.Group("/ai")
	aiRoutes.Use(middlewares.AuthMiddleware())

	// AI prediction endpoints
	predictionRoutes := aiRoutes.Group("/predict")
	predictionRoutes.Post("/skin-color-tone", c.Controllers.AI.PredictSkinColorTone)
	predictionRoutes.Post("/woman-body-scan", c.Controllers.AI.PredictWomanBodyScan)
	predictionRoutes.Post("/men-body-scan", c.Controllers.AI.PredictMenBodyScan)
}
