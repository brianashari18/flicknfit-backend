package routes

import (
	"flicknfit_backend/config"
	"flicknfit_backend/controllers"
	"flicknfit_backend/middlewares"
	"flicknfit_backend/repositories"
	"flicknfit_backend/services"
	"flicknfit_backend/utils"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes sets up all API routes for the application.
// It initializes the repositories, services, and controllers internally.
func SetupRoutes(app *fiber.App, db *gorm.DB) {
	validator := utils.NewValidator()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Config not loaded")
	}

	// Initialize Repositories
	brandRepository := repositories.NewBrandRepository(db)
	userRepository := repositories.NewUserRepository(db)
	productRepository := repositories.NewProductRepository(db)
	shoppingCartRepository := repositories.NewShoppingCartRepository(db)

	// Initialize Services
	brandService := services.NewBrandService(brandRepository)
	userService := services.NewUserService(userRepository, cfg)
	productService := services.NewProductService(productRepository)
	shoppingCartService := services.NewShoppingCartService(shoppingCartRepository, productRepository)

	// Initialize Controllers
	brandController := controllers.NewBrandController(brandService, validator)
	userController := controllers.NewUserController(userService, validator)
	productController := controllers.NewProductController(productService, validator)
	shoppingCartController := controllers.NewShoppingCartController(shoppingCartService, validator)

	api := app.Group("/api/v1")

	// ---
	// Setup public routes for the User resource.
	userRoutes := api.Group("/users")
	userRoutes.Post("/register", userController.RegisterUser)
	userRoutes.Post("/login", userController.LoginUser)
	userRoutes.Post("/forgot-password", userController.ForgotPassword)
	userRoutes.Post("/verify-otp", userController.VerifyOTP)
	userRoutes.Post("/reset-password", userController.ResetPassword)
	userRoutes.Post("/refresh-token", userController.RefreshToken)

	// Setup private routes for the User resource.
	privateRoutes := api.Group("/users")
	privateRoutes.Use(middlewares.AuthMiddleware())
	privateRoutes.Post("/logout", userController.LogoutUser)
	privateRoutes.Get("/me", userController.GetUserByAccessToken)
	privateRoutes.Patch("/edit-profile", userController.EditProfile)

	// ---
	// Setup admin routes for the User resource.
	userAdminRoutes := api.Group("/admin/users")
	userAdminRoutes.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	userAdminRoutes.Post("/", userController.AdminCreateUser)
	userAdminRoutes.Get("/", userController.AdminGetAllUsers)
	userAdminRoutes.Get("/:id", userController.AdminGetUserByID)
	userAdminRoutes.Put("/:id", userController.AdminUpdateUser)
	userAdminRoutes.Delete("/:id", userController.AdminDeleteUser)

	// ---
	// Setup public routes for the Product resource.
	productRoutes := api.Group("/products")
	productRoutes.Get("/", productController.GetAllProductsPublic)
	productRoutes.Get("/:id", productController.GetProductPublicByID)
	productRoutes.Get("/search", productController.SearchProductsPublic)
	productRoutes.Get("/filter", productController.GetAllProductsPublicWithFilter)

	// ---
	// Setup public routes for the Product reviews.
	reviewRoutes := api.Group("/products")
	reviewRoutes.Get("/:productID/reviews", productController.GetReviewsByProductIDPublic)
	reviewRoutes.Post("/:productID/reviews", middlewares.AuthMiddleware(), productController.CreateReview)

	// ---
	// Setup admin routes for the Product resource.
	productAdminRoutes := api.Group("/admin/products")
	productAdminRoutes.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	productAdminRoutes.Post("/", productController.AdminCreateProduct)
	productAdminRoutes.Get("/", productController.AdminGetAllProducts)
	productAdminRoutes.Get("/:id", productController.AdminGetProductByID)
	productAdminRoutes.Put("/:id", productController.AdminUpdateProduct)
	productAdminRoutes.Delete("/:id", productController.AdminDeleteProduct)

	// ---
	// Setup public routes for the Brand resource.
	brandRoutes := api.Group("/brands")
	brandRoutes.Get("/", brandController.GetAllBrands)
	brandRoutes.Get("/:id", brandController.GetBrandByID)

	// Setup admin routes for the Brand resource.
	brandAdminRoutes := api.Group("/admin/brands")
	brandAdminRoutes.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	brandAdminRoutes.Post("/", brandController.AdminCreateBrand)
	brandAdminRoutes.Put("/:id", brandController.AdminUpdateBrand)
	brandAdminRoutes.Delete("/:id", brandController.AdminDeleteBrand)
	brandAdminRoutes.Get("/:id", brandController.GetAllBrands)
	brandAdminRoutes.Get("/", brandController.GetBrandByID)

	// ---
	// Setup shopping cart routes.
	shoppingCartRoutes := api.Group("/cart")
	shoppingCartRoutes.Use(middlewares.AuthMiddleware())
	shoppingCartRoutes.Get("/", shoppingCartController.GetUserCart)
	shoppingCartRoutes.Post("/", shoppingCartController.AddProductItemToCart)
	shoppingCartRoutes.Put("/:itemId", shoppingCartController.UpdateProductItemInCart)
	shoppingCartRoutes.Delete("/:itemId", shoppingCartController.RemoveProductItemFromCart)
}
