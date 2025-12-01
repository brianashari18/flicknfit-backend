package seeders

import (
	"flicknfit_backend/models"
	"fmt"
	"math/rand"

	"gorm.io/gorm"
)

// SeedProductConfigurationsAllCombinations
// Setiap kombinasi warna + ukuran akan menghasilkan 1 ProductItem unik dengan stok sendiri
func SeedProductConfigurationsAllCombinations(db *gorm.DB) error {
	// Hapus semua ProductItem dan ProductConfiguration yang ada
	if err := db.Exec("DELETE FROM product_configurations").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM product_items").Error; err != nil {
		return err
	}

	// Ambil semua produk
	var products []models.Product
	if err := db.Find(&products).Error; err != nil {
		return err
	}

	// Ambil semua opsi warna dan ukuran
	var colorVar models.ProductVariation
	if err := db.Where("name = ?", "Warna").First(&colorVar).Error; err != nil {
		return fmt.Errorf("variation warna tidak ditemukan: %w", err)
	}
	var sizeVar models.ProductVariation
	if err := db.Where("name = ?", "Ukuran").First(&sizeVar).Error; err != nil {
		return fmt.Errorf("variation ukuran tidak ditemukan: %w", err)
	}

	var colorOpts []models.ProductVariationOption
	if err := db.Where("product_attribute_id = ?", colorVar.ID).Find(&colorOpts).Error; err != nil {
		return err
	}
	var sizeOpts []models.ProductVariationOption
	if err := db.Where("product_attribute_id = ?", sizeVar.ID).Find(&sizeOpts).Error; err != nil {
		return err
	}

	// Untuk setiap produk, generate ProductItem untuk setiap kombinasi warna + ukuran
	for _, product := range products {
		totalSold := 0
		for _, color := range colorOpts {
			for _, size := range sizeOpts {
				// Buat SKU unik untuk kombinasi ini
				sku := fmt.Sprintf("%s-%d-%s-%s", product.Name[:2], product.ID, color.Value[:3], size.Value)

				sold := rand.Intn(50) // Random sold 0-49
				totalSold += sold

				// Buat ProductItem baru untuk kombinasi ini
				productItem := models.ProductItem{
					ProductID: product.ID,
					SKU:       sku,
					Price:     450000 + rand.Intn(550000), // Random price 450k-1jt
					Stock:     rand.Intn(20) + 1,          // Random stock 1-20
					Sold:      sold,
					PhotoURL:  "https://static.nike.com/a/images/t_web_pw_592_v2/f_auto/7c4cce35-7b6d-44fe-8dc5-edea2b4796a5/AS+M+NL+SOLO+SWSH+BB+FZ+HOODIE.png",
				}
				if err := db.Create(&productItem).Error; err != nil {
					return fmt.Errorf("gagal create product item: %w", err)
				} // Konfigurasi warna untuk ProductItem ini
				confColor := models.ProductConfiguration{
					ProductItemID:           productItem.ID,
					ProductAttributeValueID: color.ID,
				}
				if err := db.Create(&confColor).Error; err != nil {
					return fmt.Errorf("gagal create config warna: %w", err)
				}

				// Konfigurasi ukuran untuk ProductItem ini
				confSize := models.ProductConfiguration{
					ProductItemID:           productItem.ID,
					ProductAttributeValueID: size.ID,
				}
				if err := db.Create(&confSize).Error; err != nil {
					return fmt.Errorf("gagal create config ukuran: %w", err)
				}
			}
		}
		// Update sold produk (sum dari semua varian)
		if err := db.Model(&models.Product{}).Where("id = ?", product.ID).Update("sold", totalSold).Error; err != nil {
			return fmt.Errorf("gagal update sold: %w", err)
		}
	}
	return nil
}
