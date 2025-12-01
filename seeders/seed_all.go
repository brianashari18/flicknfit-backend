package seeders

import "gorm.io/gorm"

func SeedAll(db *gorm.DB) error {
	if err := SeedBrands(db); err != nil {
		return err
	}
	if err := SeedProducts(db); err != nil {
		return err
	}
	if err := SeedProductCategories(db); err != nil {
		return err
	}
	if err := SeedProductStyles(db); err != nil {
		return err
	}
	if err := SeedProductVariationsAndOptions(db); err != nil {
		return err
	}
	if err := SeedProductConfigurationsAllCombinations(db); err != nil {
		return err
	}
	if err := SeedReviews(db); err != nil {
		return err
	}
	return nil
}
