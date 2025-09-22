package services

import (
	"errors"
	"flicknfit_backend/dtos"
	"flicknfit_backend/models"
	"flicknfit_backend/repositories"

	"gorm.io/gorm"
)

// ShoppingCartService defines the interface for business logic related to shopping carts.
type ShoppingCartService interface {
	GetOrCreateCart(userID uint64) (*models.ShoppingCart, error)
	AddProductItemToCart(userID uint64, dto *dtos.AddProductItemToCartRequestDTO) (*models.ShoppingCart, error)
	UpdateProductItemInCart(userID uint64, cartItemID uint64, dto *dtos.UpdateProductItemInCartRequestDTO) (*models.ShoppingCart, error)
	RemoveProductItemFromCart(userID uint64, cartItemID uint64) error
	GetCartItems(userID uint64) (*models.ShoppingCart, error)
}

// shoppingCartService is the implementation of ShoppingCartService.
type shoppingCartService struct {
	cartRepository    repositories.ShoppingCartRepository
	productRepository repositories.ProductRepository
}

// NewShoppingCartService creates and returns a new instance of ShoppingCartService.
func NewShoppingCartService(cartRepo repositories.ShoppingCartRepository, productRepo repositories.ProductRepository) ShoppingCartService {
	return &shoppingCartService{
		cartRepository:    cartRepo,
		productRepository: productRepo,
	}
}

// GetOrCreateCart retrieves a user's shopping cart or creates a new one if it doesn't exist.
func (s *shoppingCartService) GetOrCreateCart(userID uint64) (*models.ShoppingCart, error) {
	cart, err := s.cartRepository.GetCartByUserID(userID)
	if err != nil {
		// If the cart doesn't exist, create a new one.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newCart := &models.ShoppingCart{
				UserID: userID,
			}
			if err := s.cartRepository.CreateCart(newCart); err != nil {
				return nil, errors.New("failed to create shopping cart")
			}
			return newCart, nil
		}
		return nil, errors.New("failed to retrieve shopping cart")
	}
	return cart, nil
}

// GetCartItems retrieves all items from a user's shopping cart.
func (s *shoppingCartService) GetCartItems(userID uint64) (*models.ShoppingCart, error) {
	cart, err := s.cartRepository.GetCartByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to retrieve shopping cart")
	}
	return cart, nil
}

// AddProductItemToCart handles adding a product item to the user's cart.
func (s *shoppingCartService) AddProductItemToCart(userID uint64, dto *dtos.AddProductItemToCartRequestDTO) (*models.ShoppingCart, error) {
	cart, err := s.GetOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	productItem, err := s.productRepository.GetProductItemByID(dto.ProductItemID)
	if productItem == nil && err != nil {
		return nil, errors.New("product item not found")
	}

	// Check if item already exists in the cart.
	var cartItem *models.ShoppingCartItem
	for _, item := range cart.ShoppingCartItems {
		if item.ProductItemID == dto.ProductItemID {
			cartItem = &item
			break
		}
	}

	if cartItem != nil {
		// Update quantity if item exists.
		cartItem.Quantity += dto.Quantity
		if err := s.cartRepository.UpdateCartItem(cartItem); err != nil {
			return nil, errors.New("failed to update cart item")
		}
	} else {
		// Add new item to cart.
		cartItem = &models.ShoppingCartItem{
			ShoppingCartID: cart.ID,
			ProductItemID:  dto.ProductItemID,
			Quantity:       dto.Quantity,
		}
		if err := s.cartRepository.AddCartItem(cartItem); err != nil {
			return nil, errors.New("failed to add item to cart")
		}
	}

	// Retrieve the updated cart with the new or updated item.
	return s.cartRepository.GetCartByUserID(userID)
}

// UpdateProductItemInCart handles updating the quantity of a product item in the user's cart.
func (s *shoppingCartService) UpdateProductItemInCart(userID uint64, cartItemID uint64, dto *dtos.UpdateProductItemInCartRequestDTO) (*models.ShoppingCart, error) {
	cart, err := s.cartRepository.GetCartByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to retrieve shopping cart")
	}

	cartItem, err := s.cartRepository.GetCartItemByItemID(cartItemID)
	if err != nil {
		return nil, errors.New("cart item not found")
	}

	// Verify the cart item belongs to the user's cart.
	if cartItem.ShoppingCartID != cart.ID {
		return nil, errors.New("cart item does not belong to this user")
	}

	cartItem.Quantity = dto.Quantity
	if err := s.cartRepository.UpdateCartItem(cartItem); err != nil {
		return nil, errors.New("failed to update cart item quantity")
	}

	return s.cartRepository.GetCartByUserID(userID)
}

// RemoveProductItemFromCart handles removing an item from the user's cart.
func (s *shoppingCartService) RemoveProductItemFromCart(userID uint64, cartItemID uint64) error {
	cart, err := s.cartRepository.GetCartByUserID(userID)
	if err != nil {
		return errors.New("failed to retrieve shopping cart")
	}

	cartItem, err := s.cartRepository.GetCartItemByItemID(cartItemID)
	if err != nil {
		return errors.New("cart item not found")
	}

	// Verify the cart item belongs to the user's cart.
	if cartItem.ShoppingCartID != cart.ID {
		return errors.New("cart item does not belong to this user")
	}

	return s.cartRepository.DeleteCartItem(cartItemID)
}
