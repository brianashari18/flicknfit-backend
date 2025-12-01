package seeders

import (
	"flicknfit_backend/models"
	"math/rand"

	"gorm.io/gorm"
)

func SeedReviews(db *gorm.DB) error {
	// Ambil semua produk
	var products []models.Product
	if err := db.Find(&products).Error; err != nil {
		return err
	}

	// Ambil semua user (kalau ada, kalau belum ada user bisa skip atau buat dummy)
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	// Kalau belum ada user, buat dummy user dulu untuk review
	if len(users) == 0 {
		dummyUsers := []models.User{
			{Username: "reviewer1", Email: "reviewer1@test.com", Password: "dummy"},
			{Username: "reviewer2", Email: "reviewer2@test.com", Password: "dummy"},
			{Username: "reviewer3", Email: "reviewer3@test.com", Password: "dummy"},
		}
		for _, user := range dummyUsers {
			if err := db.Create(&user).Error; err != nil {
				return err
			}
			users = append(users, user)
		}
	}

	reviewTexts := []string{
		"Produk sangat bagus, sesuai ekspektasi!",
		"Kualitas oke, pengiriman cepat.",
		"Lumayan, tapi agak kekecilan.",
		"Recommended! Worth the price.",
		"Bahannya nyaman dipakai.",
		"Jelek, tidak sesuai deskripsi.",
		"Bagus banget, puas!",
		"Biasa saja, tidak istimewa.",
		"Sangat puas dengan pembelian ini.",
		"Kecewa, barang rusak.",
	}

	// Untuk setiap produk, buat 3-10 review random
	for _, product := range products {
		reviewCount := rand.Intn(8) + 3 // 3-10 review
		totalRating := 0

		for i := 0; i < reviewCount; i++ {
			user := users[rand.Intn(len(users))]
			rating := rand.Intn(5) + 1 // Rating 1-5
			totalRating += rating

			review := models.Review{
				UserID:     user.ID,
				ProductID:  product.ID,
				Rating:     rating,
				ReviewText: reviewTexts[rand.Intn(len(reviewTexts))],
			}

			if err := db.Create(&review).Error; err != nil {
				return err
			}
		}

		// Update rating dan reviewer di product
		avgRating := float64(totalRating) / float64(reviewCount)
		if err := db.Model(&models.Product{}).Where("id = ?", product.ID).Updates(map[string]interface{}{
			"rating":   avgRating,
			"reviewer": reviewCount,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}
