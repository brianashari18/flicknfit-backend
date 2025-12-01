package seeders

import (
	"flicknfit_backend/models"

	"gorm.io/gorm"
)

func SeedProductCategories(db *gorm.DB) error {
	categories := []string{
		"T-Shirts",
		"Shirts",
		"Pants",
		"Jackets",
		"Cardigans",
		"Skirts",
		"Tank Tops",
		"Sweaters",
		"Jeans",
		"Crop Tops",
		"Dress",
	}

	var products []models.Product
	if err := db.Find(&products).Error; err != nil {
		return err
	}

	for _, product := range products {
		// Setiap produk hanya dapat 1 kategori (berdasarkan ID)
		cat := categories[int(product.ID)%len(categories)]
		catEntry := models.ProductCategory{
			ProductID: product.ID,
			Category:  cat,
		}
		if err := db.Where(models.ProductCategory{ProductID: product.ID, Category: cat}).FirstOrCreate(&catEntry).Error; err != nil {
			return err
		}
	}
	return nil
}

func SeedProductVariationsAndOptions(db *gorm.DB) error {
	// Variasi: Warna dan Ukuran
	colors := []string{"Merah", "Biru", "Hitam", "Putih", "Hijau", "Kuning"}
	sizes := []string{"S", "M", "L", "XL", "XXL"}

	// Buat variation Warna
	colorVar := models.ProductVariation{Name: "Warna"}
	if err := db.Where(models.ProductVariation{Name: colorVar.Name}).FirstOrCreate(&colorVar).Error; err != nil {
		return err
	}
	for _, c := range colors {
		opt := models.ProductVariationOption{
			ProductAttributeID: colorVar.ID,
			Value:              c,
		}
		if err := db.Where(models.ProductVariationOption{ProductAttributeID: colorVar.ID, Value: c}).FirstOrCreate(&opt).Error; err != nil {
			return err
		}
	}

	// Buat variation Ukuran
	sizeVar := models.ProductVariation{Name: "Ukuran"}
	if err := db.Where(models.ProductVariation{Name: sizeVar.Name}).FirstOrCreate(&sizeVar).Error; err != nil {
		return err
	}
	for _, s := range sizes {
		opt := models.ProductVariationOption{
			ProductAttributeID: sizeVar.ID,
			Value:              s,
		}
		if err := db.Where(models.ProductVariationOption{ProductAttributeID: sizeVar.ID, Value: s}).FirstOrCreate(&opt).Error; err != nil {
			return err
		}
	}

	return nil
}
