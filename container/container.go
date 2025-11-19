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
	User         repositories.UserRepository
	Brand        repositories.BrandRepository
	Product      repositories.ProductRepository
	ShoppingCart repositories.ShoppingCartRepository
	Favorite     repositories.FavoriteRepository
	Review       repositories.ReviewRepository
	Wardrobe     repositories.WardrobeRepository
}

// Services holds all service instances
type Services struct {
	User         services.UserService
	Brand        services.BrandService
	Product      services.ProductService
	ShoppingCart services.ShoppingCartService
	Favorite     services.FavoriteService
	Review       services.ReviewService
	Wardrobe     services.WardrobeService
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
		User:         repositories.NewUserRepository(c.DB),
		Brand:        repositories.NewBrandRepository(c.DB),
		Product:      repositories.NewProductRepository(c.DB),
		ShoppingCart: repositories.NewShoppingCartRepository(c.DB),
		Favorite:     repositories.NewFavoriteRepository(c.DB),
		Review:       repositories.NewReviewRepository(c.DB),
		Wardrobe:     repositories.NewWardrobeRepository(c.DB),
	}
}

// initServices initializes all service instances
func (c *Container) initServices() {
	c.Services = &Services{
		User:         services.NewUserService(c.Repositories.User, c.Config),
		Brand:        services.NewBrandService(c.Repositories.Brand),
		Product:      services.NewProductService(c.Repositories.Product),
		ShoppingCart: services.NewShoppingCartService(c.Repositories.ShoppingCart, c.Repositories.Product),
		Favorite:     services.NewFavoriteService(c.Repositories.Favorite, c.Repositories.Product),
		Review:       services.NewReviewService(c.Repositories.Review, c.Repositories.Product),
		Wardrobe:     services.NewWardrobeService(c.Repositories.Wardrobe),
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
	}
}
