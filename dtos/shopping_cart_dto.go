package dtos

import "flicknfit_backend/models"

// ShoppingCartDTO represents the shopping cart data returned to the user.
type ShoppingCartDTO struct {
	ID        uint64                `json:"id"`
	UserID    uint64                `json:"user_id"`
	CreatedAt string                `json:"created_at"`
	UpdatedAt string                `json:"updated_at"`
	Items     []ShoppingCartItemDTO `json:"items"`
}

// ShoppingCartItemDTO represents a single item in the shopping cart.
type ShoppingCartItemDTO struct {
	ID            uint64         `json:"id"`
	ProductItemID uint64         `json:"product_item_id"`
	Quantity      int            `json:"quantity"`
	ItemPrice     int            `json:"item_price"`
	ProductItem   ProductItemDTO `json:"product_item"`
}

// AddProductItemToCartRequestDTO is used for adding a new product item to the shopping cart.
type AddProductItemToCartRequestDTO struct {
	ProductItemID uint64 `json:"product_item_id" validate:"required"`
	Quantity      int    `json:"quantity" validate:"required,min=1"`
}

// UpdateProductItemInCartRequestDTO is used for updating an existing product item in the shopping cart.
type UpdateProductItemInCartRequestDTO struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}

func ToShoppingCartDTO(cart *models.ShoppingCart) ShoppingCartDTO {
	itemsDTOs := make([]ShoppingCartItemDTO, len(cart.ShoppingCartItems))
	for i, item := range cart.ShoppingCartItems {
		itemsDTOs[i] = ToShoppingCartItemDTO(&item)
	}

	return ShoppingCartDTO{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CreatedAt: cart.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: cart.UpdatedAt.Format("2006-01-02 15:04:05"),
		Items:     itemsDTOs,
	}
}

func ToShoppingCartItemDTO(item *models.ShoppingCartItem) ShoppingCartItemDTO {
	return ShoppingCartItemDTO{
		ID:            item.ID,
		ProductItemID: item.ProductItemID,
		Quantity:      item.Quantity,
		ItemPrice:     item.ProductItem.Price,
		ProductItem:   ToProductItemDTO(item.ProductItem),
	}
}
