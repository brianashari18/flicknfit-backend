package repositories

import (
	"flicknfit_backend/models"
	"time"

	"gorm.io/gorm"
)

// ProductClickRepository handles database operations for product clicks
type ProductClickRepository interface {
	Create(click *models.ProductClick) error
	GetByID(id uint64) (*models.ProductClick, error)
	GetByUserID(userID uint64, limit int) ([]models.ProductClick, error)
	GetByProductID(productID uint64, limit int) ([]models.ProductClick, error)
	GetByBrandID(brandID uint64, limit int) ([]models.ProductClick, error)
	CountByProductID(productID uint64) (int64, error)
	CountByBrandID(brandID uint64) (int64, error)
	CountByBrandIDAndDateRange(brandID uint64, startDate, endDate time.Time) (int64, error)
	GetTopClickedProducts(limit int) ([]map[string]interface{}, error)
	GetClickStatsByBrand(brandID uint64, startDate, endDate time.Time) ([]map[string]interface{}, error)
}

type productClickRepository struct {
	db *gorm.DB
}

// NewProductClickRepository creates a new product click repository
func NewProductClickRepository(db *gorm.DB) ProductClickRepository {
	return &productClickRepository{db: db}
}

func (r *productClickRepository) Create(click *models.ProductClick) error {
	return r.db.Create(click).Error
}

func (r *productClickRepository) GetByID(id uint64) (*models.ProductClick, error) {
	var click models.ProductClick
	err := r.db.Preload("User").Preload("Product").Preload("Brand").
		First(&click, id).Error
	return &click, err
}

func (r *productClickRepository) GetByUserID(userID uint64, limit int) ([]models.ProductClick, error) {
	var clicks []models.ProductClick
	query := r.db.Where("user_id = ?", userID).
		Order("clicked_at DESC").
		Preload("Product").Preload("Brand")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&clicks).Error
	return clicks, err
}

func (r *productClickRepository) GetByProductID(productID uint64, limit int) ([]models.ProductClick, error) {
	var clicks []models.ProductClick
	query := r.db.Where("product_id = ?", productID).
		Order("clicked_at DESC").
		Preload("User")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&clicks).Error
	return clicks, err
}

func (r *productClickRepository) GetByBrandID(brandID uint64, limit int) ([]models.ProductClick, error) {
	var clicks []models.ProductClick
	query := r.db.Where("brand_id = ?", brandID).
		Order("clicked_at DESC").
		Preload("Product").Preload("User")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&clicks).Error
	return clicks, err
}

func (r *productClickRepository) CountByProductID(productID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&models.ProductClick{}).
		Where("product_id = ?", productID).
		Count(&count).Error
	return count, err
}

func (r *productClickRepository) CountByBrandID(brandID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&models.ProductClick{}).
		Where("brand_id = ?", brandID).
		Count(&count).Error
	return count, err
}

func (r *productClickRepository) CountByBrandIDAndDateRange(brandID uint64, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.ProductClick{}).
		Where("brand_id = ? AND clicked_at BETWEEN ? AND ?", brandID, startDate, endDate).
		Count(&count).Error
	return count, err
}

// GetTopClickedProducts returns the most clicked products with click counts
func (r *productClickRepository) GetTopClickedProducts(limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.db.Model(&models.ProductClick{}).
		Select("product_id, COUNT(*) as click_count").
		Group("product_id").
		Order("click_count DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

// GetClickStatsByBrand returns click statistics grouped by product for a specific brand
func (r *productClickRepository) GetClickStatsByBrand(brandID uint64, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.db.Model(&models.ProductClick{}).
		Select("product_id, COUNT(*) as click_count, COUNT(DISTINCT user_id) as unique_users").
		Where("brand_id = ? AND clicked_at BETWEEN ? AND ?", brandID, startDate, endDate).
		Group("product_id").
		Order("click_count DESC").
		Scan(&results).Error

	return results, err
}
