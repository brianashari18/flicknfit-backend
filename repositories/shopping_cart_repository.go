package repositories

import (
	"flicknfit_backend/models"

	"gorm.io/gorm"
)

// ShoppingCartRepository defines the interface for data access operations on the ShoppingCart model.
//
//go:generate mockery --name ShoppingCartRepository
type ShoppingCartRepository interface {
	GetCartByUserID(userID uint64) (*models.ShoppingCart, error)
	CreateCart(cart *models.ShoppingCart) error
	GetCartItemByID(cartID uint64, productItemID uint64) (*models.ShoppingCartItem, error)
	AddCartItem(item *models.ShoppingCartItem) error
	UpdateCartItem(item *models.ShoppingCartItem) error
	DeleteCartItem(id uint64) error
	GetCartItemByItemID(id uint64) (*models.ShoppingCartItem, error)
	SaveCart(cart *models.ShoppingCart) error
}

// shoppingCartRepository is the implementation of ShoppingCartRepository.
type shoppingCartRepository struct {
	BaseRepository
}

// NewShoppingCartRepository creates and returns a new instance of ShoppingCartRepository.
func NewShoppingCartRepository(db *gorm.DB) ShoppingCartRepository {
	return &shoppingCartRepository{BaseRepository: BaseRepository{DB: db}}
}

// GetCartByUserID retrieves a shopping cart by its user ID.
func (r *shoppingCartRepository) GetCartByUserID(userID uint64) (*models.ShoppingCart, error) {
	var cart models.ShoppingCart
	err := r.DB.Where("user_id = ?", userID).Preload("ShoppingCartItems.ProductItem").First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// CreateCart creates a new shopping cart record in the database.
func (r *shoppingCartRepository) CreateCart(cart *models.ShoppingCart) error {
	return r.DB.Create(cart).Error
}

// GetCartItemByID retrieves a shopping cart item by its cart ID and product item ID.
func (r *shoppingCartRepository) GetCartItemByID(cartID uint64, productItemID uint64) (*models.ShoppingCartItem, error) {
	var item models.ShoppingCartItem
	err := r.DB.Where("shopping_cart_id = ? AND product_item_id = ?", cartID, productItemID).Preload("ProductItem").First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// AddCartItem adds a new item to the shopping cart.
func (r *shoppingCartRepository) AddCartItem(item *models.ShoppingCartItem) error {
	return r.DB.Create(item).Error
}

// UpdateCartItem updates an existing item in the shopping cart.
func (r *shoppingCartRepository) UpdateCartItem(item *models.ShoppingCartItem) error {
	return r.DB.Save(item).Error
}

// DeleteCartItem deletes an item from the shopping cart by its ID.
func (r *shoppingCartRepository) DeleteCartItem(id uint64) error {
	return r.DB.Delete(&models.ShoppingCartItem{}, id).Error
}

func (r *shoppingCartRepository) GetCartItemByItemID(id uint64) (*models.ShoppingCartItem, error) {
	var item models.ShoppingCartItem
	err := r.DB.Preload("ProductItem").First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *shoppingCartRepository) SaveCart(cart *models.ShoppingCart) error {
	return r.DB.Save(cart).Error
}
