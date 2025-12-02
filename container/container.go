package container

import (
	"flicknfit_backend/config"
	"flicknfit_backend/controllers"
	"flicknfit_backend/repositories"
	"flicknfit_backend/services"
	"flicknfit_backend/utils"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Container holds all application dependencies
type Container struct {
	DB           *gorm.DB
	Config       *config.Config
	Validator    *validator.Validate
	Repositories *Repositories
	Services     *Services
	Controllers  *Controllers
}

// Repositories holds all repository instances
type Repositories struct {
	User            repositories.UserRepository
	Brand           repositories.BrandRepository
	Product         repositories.ProductRepository
	ShoppingCart    repositories.ShoppingCartRepository
	Favorite        repositories.FavoriteRepository
	Review          repositories.ReviewRepository
	Wardrobe        repositories.WardrobeRepository
	FaceScanHistory repositories.FaceScanHistoryRepository
	BodyScanHistory repositories.BodyScanHistoryRepository
}

// Services holds all service instances
type Services struct {
	User            services.UserService
	Brand           services.BrandService
	Product         services.ProductService
	ShoppingCart    services.ShoppingCartService
	Favorite        services.FavoriteService
	Review          services.ReviewService
	Wardrobe        services.WardrobeService
	AI              services.AIService
	Firebase        *services.FirebaseService
	ScanHistory     services.ScanHistoryService
	SupabaseStorage services.SupabaseStorageService
}

// Controllers holds all controller instances
type Controllers struct {
	User         controllers.UserController
	Brand        controllers.BrandController
	Product      controllers.ProductController
	ShoppingCart controllers.ShoppingCartController
	Favorite     controllers.FavoriteController
	Review       controllers.ReviewController
	Wardrobe     controllers.WardrobeController
	AI           controllers.AIController
	Dashboard    controllers.DashboardController
	OAuth        *controllers.OAuthController
	ScanHistory  controllers.ScanHistoryController
}

// NewContainer creates and initializes a new container with all dependencies
func NewContainer(db *gorm.DB, cfg *config.Config) (*Container, error) {
	container := &Container{
		DB:     db,
		Config: cfg,
	}
	// Initialize validator
	container.Validator = utils.NewValidator()

	// Initialize repositories
	container.initRepositories()

	// Initialize services
	container.initServices()

	// Initialize controllers
	container.initControllers()

	return container, nil
}

// initRepositories initializes all repository instances
func (c *Container) initRepositories() {
	c.Repositories = &Repositories{
		User:            repositories.NewUserRepository(c.DB),
		Brand:           repositories.NewBrandRepository(c.DB),
		Product:         repositories.NewProductRepository(c.DB),
		ShoppingCart:    repositories.NewShoppingCartRepository(c.DB),
		Favorite:        repositories.NewFavoriteRepository(c.DB),
		Review:          repositories.NewReviewRepository(c.DB),
		Wardrobe:        repositories.NewWardrobeRepository(c.DB),
		FaceScanHistory: repositories.NewFaceScanHistoryRepository(c.DB),
		BodyScanHistory: repositories.NewBodyScanHistoryRepository(c.DB),
	}
}

// initServices initializes all service instances
func (c *Container) initServices() {
	// Initialize Firebase service
	firebaseService, err := services.NewFirebaseService(c.Config)
	if err != nil {
		utils.GetLogger().Warn("Firebase service initialization failed: ", err)
		// Continue without Firebase - it's optional
	}

	// Initialize Supabase Storage service
	var supabaseStorageService services.SupabaseStorageService
	if c.Config.SupabaseURL != "" && c.Config.SupabaseKey != "" {
		supabaseStorageService = services.NewSupabaseStorageService(c.Config)
		utils.GetLogger().Info("Supabase storage service initialized")
	} else {
		utils.GetLogger().Warn("Supabase storage service not initialized - credentials missing")
	}

	c.Services = &Services{
		User:            services.NewUserService(c.Repositories.User, c.Config),
		Brand:           services.NewBrandService(c.Repositories.Brand),
		Product:         services.NewProductService(c.Repositories.Product),
		ShoppingCart:    services.NewShoppingCartService(c.Repositories.ShoppingCart, c.Repositories.Product),
		Favorite:        services.NewFavoriteService(c.Repositories.Favorite, c.Repositories.Product),
		Review:          services.NewReviewService(c.Repositories.Review, c.Repositories.Product),
		Wardrobe:        services.NewWardrobeService(c.Repositories.Wardrobe),
		AI:              services.NewAIService(c.Config),
		Firebase:        firebaseService,
		SupabaseStorage: supabaseStorageService,
		ScanHistory:     services.NewScanHistoryService(c.Repositories.FaceScanHistory, c.Repositories.BodyScanHistory, supabaseStorageService),
	}
}

// initControllers initializes all controller instances
func (c *Container) initControllers() {
	c.Controllers = &Controllers{
		User:         controllers.NewUserController(c.Services.User, c.Validator),
		Brand:        controllers.NewBrandController(c.Services.Brand, c.Validator),
		Product:      controllers.NewProductController(c.Services.Product, c.Validator),
		ShoppingCart: controllers.NewShoppingCartController(c.Services.ShoppingCart, c.Validator),
		Favorite:     controllers.NewFavoriteController(c.Services.Favorite, c.Validator),
		Review:       controllers.NewReviewController(c.Services.Review, c.Validator),
		Wardrobe:     controllers.NewWardrobeController(c.Services.Wardrobe, c.Validator),
		AI:           controllers.NewAIController(c.Services.AI, c.Services.ScanHistory),
		Dashboard:    controllers.NewDashboardController(c.DB, c.Services.User, c.Services.Brand),
		OAuth:        controllers.NewOAuthController(c.Services.User, c.Services.Firebase),
		ScanHistory:  controllers.NewScanHistoryController(c.Services.ScanHistory, c.Services.SupabaseStorage),
	}
}
