package dtos

import "flicknfit_backend/models"

// SavedItemsDTO represents the shopping cart data returned to the user.
type SavedItemsDTO struct {
	ID        uint64         `json:"id"`
	UserID    uint64         `json:"user_id"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
	Items     []SavedItemDTO `json:"items"`
}

// SavedItemDTO represents a single item in the shopping cart.
type SavedItemDTO struct {
	ID            uint64         `json:"id"`
	ProductItemID uint64         `json:"product_item_id"`
	Quantity      int            `json:"quantity"`
	ItemPrice     int            `json:"item_price"`
	ProductItem   ProductItemDTO `json:"product_item"`
}

// AddProductItemToSavedItemsRequestDTO is used for adding a new product item to the shopping cart.
type AddProductItemToSavedItemsRequestDTO struct {
	ProductItemID uint64 `json:"product_item_id" validate:"required"`
	Quantity      int    `json:"quantity" validate:"required,min=1"`
}

// UpdateProductItemInSavedItemsRequestDTO is used for updating an existing product item in the shopping cart.
type UpdateProductItemInSavedItemsRequestDTO struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}

func ToSavedItemsDTO(cart *models.SavedItems) SavedItemsDTO {
	itemsDTOs := make([]SavedItemDTO, len(cart.SavedItemsList))
	for i, item := range cart.SavedItemsList {
		itemsDTOs[i] = ToSavedItemDTO(&item)
	}

	return SavedItemsDTO{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CreatedAt: cart.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: cart.UpdatedAt.Format("2006-01-02 15:04:05"),
		Items:     itemsDTOs,
	}
}

func ToSavedItemDTO(item *models.SavedItemsList) SavedItemDTO {
	return SavedItemDTO{
		ID:            item.ID,
		ProductItemID: item.ProductItemID,
		Quantity:      item.Quantity,
		ItemPrice:     item.ProductItem.Price,
		ProductItem:   ToProductItemDTO(item.ProductItem),
	}
}
