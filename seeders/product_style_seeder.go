package seeders

import (
	"flicknfit_backend/models"

	"gorm.io/gorm"
)

func SeedProductStyles(db *gorm.DB) error {
	styles := []string{
		"Minimalist",
		"Casual",
		"Formal",
		"Sporty",
		"Streetwear",
		"Vintage",
		"Elegant",
		"Bohemian",
		"Preppy",
		"Athleisure",
	}

	var products []models.Product
	if err := db.Find(&products).Error; err != nil {
		return err
	}

	for _, product := range products {
		// Setiap produk dapat 1-2 style random
		numStyles := 1 + (int(product.ID) % 2) // 1 atau 2 style
		for i := 0; i < numStyles; i++ {
			style := styles[(int(product.ID)+i)%len(styles)]
			styleEntry := models.ProductStyle{
				ProductID: product.ID,
				Style:     style,
			}
			if err := db.Where(models.ProductStyle{ProductID: product.ID, Style: style}).FirstOrCreate(&styleEntry).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
