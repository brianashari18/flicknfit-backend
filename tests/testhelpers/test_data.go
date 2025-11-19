package testhelpers

import (
	"flicknfit_backend/models"
	"time"
)

// CreateTestUser creates a test user for testing purposes
func CreateTestUser() models.User {
	birthday := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	return models.User{
		ID:          1,
		Email:       "test@example.com",
		Username:    "testuser",
		Password:    "hashedpassword",
		PhoneNumber: "+1234567890",
		Gender:      models.MaleGender,
		Birthday:    &birthday,
		Region:      "Jakarta",
		Role:        models.UserRole,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// CreateTestAdminUser creates a test admin user
func CreateTestAdminUser() models.User {
	user := CreateTestUser()
	user.ID = 2
	user.Email = "admin@example.com"
	user.Username = "admin"
	user.Role = models.AdminRole
	return user
}

// CreateTestProduct creates a test product for testing
func CreateTestProduct() models.Product {
	return models.Product{
		ID:          1,
		BrandID:     1,
		Name:        "Test Product",
		Description: "Test Description",
		Discount:    0.1,
		Rating:      4.5,
		Reviewer:    10,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// CreateTestProductItem creates a test product item
func CreateTestProductItem() models.ProductItem {
	return models.ProductItem{
		ID:        1,
		ProductID: 1,
		SKU:       "TEST-SKU-001",
		Price:     100000,
		Stock:     50,
		PhotoURL:  "https://example.com/test.jpg",
		Product:   CreateTestProduct(),
	}
}

// CreateTestFavorite creates a test favorite
func CreateTestFavorite() models.Favorite {
	return models.Favorite{
		ID:            1,
		UserID:        1,
		ProductItemID: 1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		User:          CreateTestUser(),
		ProductItem:   CreateTestProductItem(),
	}
}

// CreateTestReview creates a test review
func CreateTestReview() models.Review {
	return models.Review{
		ID:         1,
		UserID:     1,
		ProductID:  1,
		Rating:     5,
		ReviewText: "Great product!",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// CreateTestUserWardrobe creates a test wardrobe item
func CreateTestUserWardrobe() models.UserWardrobe {
	return models.UserWardrobe{
		ID:        1,
		UserID:    1,
		Category:  "Shirts",
		ImageURL:  "https://example.com/shirt.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		User:      CreateTestUser(),
	}
}
