package seeders

import (
	"flicknfit_backend/models"

	"gorm.io/gorm"
)

func SeedBrands(db *gorm.DB) error {
	brands := []models.Brand{
		{
			Name:        "Nike",
			Description: "The Hero (Performa & Motivasi)",
			LogoURL:     "https://upload.wikimedia.org/wikipedia/commons/a/a6/Logo_NIKE.svg",
			WebsiteURL:  "https://www.nike.com/",
		},
		{
			Name:        "Adidas",
			Description: "The Creator (Kolaborasi & Budaya)",
			LogoURL:     "https://upload.wikimedia.org/wikipedia/commons/2/20/Adidas_Logo.svg",
			WebsiteURL:  "https://www.adidas.com/",
		},
		{
			Name:        "Zara",
			Description: "The Ruler (Otoritas Tren Cepat)",
			LogoURL:     "https://upload.wikimedia.org/wikipedia/commons/f/fd/Zara_Logo.svg",
			WebsiteURL:  "https://www.zara.com/",
		},
		{
			Name:        "H&M",
			Description: "The Everyman (Aksesibilitas Ceria)",
			LogoURL:     "https://upload.wikimedia.org/wikipedia/commons/5/53/H%26M-Logo.svg",
			WebsiteURL:  "https://www2.hm.com/",
		},
		{
			Name:        "Uniqlo",
			Description: "The Sage (Rasionalitas & Kualitas)",
			LogoURL:     "https://upload.wikimedia.org/wikipedia/commons/9/92/UNIQLO_logo.svg",
			WebsiteURL:  "https://www.uniqlo.com/",
		},
		{
			Name:        "Gucci",
			Description: "The Jester/Rebel (Ekspresi Diri Eklektik)",
			LogoURL:     "https://upload.wikimedia.org/wikipedia/commons/2/2e/Gucci_Logo.svg",
			WebsiteURL:  "https://www.gucci.com/",
		},
		{
			Name:        "Louis Vuitton",
			Description: "The Explorer (Perjalanan & Warisan)",
			LogoURL:     "LV Logo (PNG)",
			WebsiteURL:  "https://us.louisvuitton.com/",
		},
		{
			Name:        "Chanel",
			Description: "The Lover (Intimacy & Misteri)",
			LogoURL:     "https://upload.wikimedia.org/wikipedia/commons/3/35/Chanel_logo.svg",
			WebsiteURL:  "https://www.chanel.com/",
		},
		{
			Name:        "Hermès",
			Description: "The Creator (Pengrajin Artistik)",
			LogoURL:     "https://upload.wikimedia.org/wikipedia/commons/c/c9/Hermes_wordmark.svg",
			WebsiteURL:  "https://www.hermes.com/",
		},
		{
			Name:        "Prada",
			Description: "The Magician (Intelek & Subversi)",
			LogoURL:     "https://upload.wikimedia.org/wikipedia/commons/b/b8/Prada-Logo.svg",
			WebsiteURL:  "https://www.prada.com/",
		},
	}
	for _, brand := range brands {
		if err := db.Where(models.Brand{Name: brand.Name}).FirstOrCreate(&brand).Error; err != nil {
			return err
		}
	}
	return nil
}
