package seeders

import (
	"fmt"
	"math/rand"
	"time"

	"flicknfit_backend/models"

	"gorm.io/gorm"
)

func SeedProducts(db *gorm.DB) error {
	brands := []string{"Nike", "Adidas", "Zara", "H&M", "Uniqlo", "Gucci", "Louis Vuitton", "Chanel", "Hermès", "Prada"}
	photoURL := "https://static.nike.com/a/images/t_web_pw_592_v2/f_auto/7c4cce35-7b6d-44fe-8dc5-edea2b4796a5/AS+M+NL+SOLO+SWSH+BB+FZ+HOODIE.png"
	desc := "Produk fashion eksklusif dari brand ternama."
	rand.Seed(time.Now().UnixNano())

	for _, brandName := range brands {
		var brand models.Brand
		if err := db.Where("name = ?", brandName).First(&brand).Error; err != nil {
			return fmt.Errorf("brand %s not found: %w", brandName, err)
		}
		for i := 1; i <= 10; i++ {
			product := models.Product{
				BrandID:     brand.ID,
				Name:        fmt.Sprintf("%s Product %d", brand.Name, i),
				Description: desc,
				Discount:    float64(rand.Intn(50)),
				Rating:      0.0,
				Reviewer:    0,
				// CreatedAt dan UpdatedAt otomatis oleh GORM
			}
			if err := db.Where(models.Product{Name: product.Name, BrandID: brand.ID}).FirstOrCreate(&product).Error; err != nil {
				return err
			}

			// Buat 1 ProductItem per product
			item := models.ProductItem{
				ProductID: product.ID,
				SKU:       fmt.Sprintf("%s-%03d", brand.Name[:2], i),
				Price:     100000 + rand.Intn(900000),
				Stock:     10 + rand.Intn(90),
				PhotoURL:  photoURL,
			}
			if err := db.Where(models.ProductItem{SKU: item.SKU}).FirstOrCreate(&item).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
