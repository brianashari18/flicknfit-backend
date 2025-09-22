package repositories

import (
	"flicknfit_backend/dtos"
	"flicknfit_backend/models"

	"gorm.io/gorm"
)

// ProductRepository mendefinisikan antarmuka untuk operasi akses data pada model Product.
type ProductRepository interface {
	// Metode untuk CRUD produk oleh admin
	CreateProduct(product *models.Product) error
	UpdateProduct(product *models.Product) error
	DeleteProduct(id uint64) error
	GetProductByID(id uint64) (*models.Product, error)
	GetAllProducts() ([]*models.Product, error)
	CreateReview(review *models.Review) error
	GetReviewsByProductID(productID uint64) ([]*models.Review, error)
	UpdateReview(review *models.Review) error
	DeleteReview(reviewID uint64) error
	GetReviewByID(reviewID uint64) (*models.Review, error)
	GetProductItemByID(id uint64) (*models.ProductItem, error)

	// Metode untuk user biasa
	GetProductPublicByID(id uint64) (*models.Product, error)
	GetAllProductsPublic() ([]*models.Product, error)

	// Metode baru untuk pencarian produk.
	// SearchProducts mencari produk berdasarkan nama atau deskripsi.
	SearchProducts(query string) ([]*models.Product, error)

	GetAllProductsPublicWithFilter(filter *dtos.ProductFilterRequestDTO) ([]*models.Product, error)
}

// productRepository adalah implementasi dari ProductRepository.
type productRepository struct {
	BaseRepository
}

// NewProductRepository membuat dan mengembalikan instance baru dari ProductRepository.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{BaseRepository: BaseRepository{DB: db}}
}

// CreateProduct membuat produk baru di database.
func (r *productRepository) CreateProduct(product *models.Product) error {
	return r.DB.Create(product).Error
}

// UpdateProduct memperbarui produk yang sudah ada.
func (r *productRepository) UpdateProduct(product *models.Product) error {
	return r.DB.Save(product).Error
}

// DeleteProduct menghapus produk dari database.
func (r *productRepository) DeleteProduct(id uint64) error {
	return r.DB.Delete(&models.Product{}, id).Error
}

// GetProductByID mengambil produk berdasarkan ID-nya, memuat juga variasi dan konfigurasinya.
func (r *productRepository) GetProductByID(id uint64) (*models.Product, error) {
	var product models.Product
	if err := r.DB.
		Preload("ProductItems.Configurations.ProductVariationOption").
		First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// GetAllProducts mengambil semua produk, memuat juga variasi dan konfigurasinya.
func (r *productRepository) GetAllProducts() ([]*models.Product, error) {
	var products []*models.Product
	if err := r.DB.
		Preload("ProductItems.Configurations.ProductVariationOption").
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetProductPublicByID mengambil produk berdasarkan ID-nya untuk pengguna publik, memuat variasi, gambar, dan review.
func (r *productRepository) GetProductPublicByID(id uint64) (*models.Product, error) {
	var product models.Product
	if err := r.DB.
		Preload("ProductItems.Configurations.ProductVariationOption").
		Preload("Reviews").
		First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// GetAllProductsPublic mengambil semua produk untuk pengguna publik, memuat variasi, gambar, dan review.
func (r *productRepository) GetAllProductsPublic() ([]*models.Product, error) {
	var products []*models.Product
	if err := r.DB.
		Preload("ProductItems.Configurations.ProductVariationOption").
		Preload("Reviews").
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// CreateReview membuat review baru di database.
func (r *productRepository) CreateReview(review *models.Review) error {
	return r.DB.Create(review).Error
}

// GetReviewsByProductID mengambil semua review untuk suatu produk.
func (r *productRepository) GetReviewsByProductID(productID uint64) ([]*models.Review, error) {
	var reviews []*models.Review
	if err := r.DB.Where("product_id = ?", productID).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

// UpdateReview memperbarui review yang sudah ada.
func (r *productRepository) UpdateReview(review *models.Review) error {
	return r.DB.Save(review).Error
}

// DeleteReview menghapus review dari database.
func (r *productRepository) DeleteReview(reviewID uint64) error {
	return r.DB.Delete(&models.Review{}, reviewID).Error
}

// GetReviewByID mengambil review berdasarkan ID-nya.
func (r *productRepository) GetReviewByID(reviewID uint64) (*models.Review, error) {
	var review models.Review
	if err := r.DB.First(&review, reviewID).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

// GetProductItemByID mengambil item produk berdasarkan ID-nya.
func (r *productRepository) GetProductItemByID(id uint64) (*models.ProductItem, error) {
	var item models.ProductItem
	if err := r.DB.Preload("Product").First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// SearchProducts mengimplementasikan pencarian produk berdasarkan nama atau deskripsi.
func (r *productRepository) SearchProducts(query string) ([]*models.Product, error) {
	var products []*models.Product
	searchTerm := "%" + query + "%"
	// Menggunakan `Where` dengan `OR` untuk mencari kecocokan pada nama atau deskripsi.
	// `Preload` digunakan untuk memuat data terkait (product items).
	if err := r.DB.
		Preload("ProductItems.Configurations.ProductVariationOption").
		Preload("Reviews").
		Where("name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm).Preload("ProductItems").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetAllProductsPublicWithFilter mengambil semua produk publik dengan filter.
func (r *productRepository) GetAllProductsPublicWithFilter(filter *dtos.ProductFilterRequestDTO) ([]*models.Product, error) {
	var products []*models.Product
	tx := r.DB.Model(&models.Product{})

	tx.
		Preload("ProductItems.Configurations.ProductVariationOption").
		Preload("Reviews")

	// Filter berdasarkan nama produk
	if filter.Name != "" {
		tx = tx.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	// Filter berdasarkan rentang harga
	if filter.MinPrice > 0 {
		// Menggunakan Joins untuk mengakses tabel product_items
		tx = tx.Joins("JOIN product_items ON products.id = product_items.product_id").
			Where("product_items.price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		tx = tx.Joins("JOIN product_items ON products.id = product_items.product_id").
			Where("product_items.price <= ?", filter.MaxPrice)
	}

	// Filter berdasarkan brand
	if filter.BrandName != "" {
		tx = tx.Joins("JOIN brands ON products.brand_id = brands.id").
			Where("brands.name ILIKE ?", "%"+filter.BrandName+"%")
	}

	// Filter berdasarkan kategori
	if filter.Category != "" {
		tx = tx.Joins("JOIN product_categories ON products.id = product_categories.product_id").
			Where("product_categories.category ILIKE ?", "%"+filter.Category+"%")
	}

	// Filter berdasarkan rating minimum
	if filter.MinRating > 0 {
		tx = tx.Where("rating >= ?", filter.MinRating)
	}

	// Menghindari duplikat jika menggunakan Joins
	tx = tx.Group("products.id")

	if err := tx.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
