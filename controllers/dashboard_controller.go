package controllers

import (
	"flicknfit_backend/services"
	"flicknfit_backend/utils"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// DashboardController handles admin dashboard operations
type DashboardController interface {
	GetDashboardStats(c *fiber.Ctx) error
	GetUserAnalytics(c *fiber.Ctx) error
	GetProductAnalytics(c *fiber.Ctx) error
	GetRevenueAnalytics(c *fiber.Ctx) error
}

type dashboardController struct {
	db           *gorm.DB
	userService  services.UserService
	brandService services.BrandService
}

// NewDashboardController creates a new dashboard controller
func NewDashboardController(db *gorm.DB, userService services.UserService, brandService services.BrandService) DashboardController {
	return &dashboardController{
		db:           db,
		userService:  userService,
		brandService: brandService,
	}
}

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalUsers       int64   `json:"total_users"`
	TotalProducts    int64   `json:"total_products"`
	TotalBrands      int64   `json:"total_brands"`
	TotalReviews     int64   `json:"total_reviews"`
	ActiveUsers      int64   `json:"active_users"`
	RevenueThisMonth float64 `json:"revenue_this_month"`
	NewUsersToday    int64   `json:"new_users_today"`
	AIRequestsToday  int64   `json:"ai_requests_today"`
}

// UserAnalytics represents user analytics data
type UserAnalytics struct {
	NewUsersLastWeek    []int               `json:"new_users_last_week"`
	UsersByRole         []RoleCount         `json:"users_by_role"`
	ActiveUsersLastHour int64               `json:"active_users_last_hour"`
	TopUsersByFavorites []UserFavoriteCount `json:"top_users_by_favorites"`
}

type RoleCount struct {
	Role  string `json:"role"`
	Count int64  `json:"count"`
}

type UserFavoriteCount struct {
	Username      string `json:"username"`
	FavoriteCount int64  `json:"favorite_count"`
}

// ProductAnalytics represents product analytics
type ProductAnalytics struct {
	TopProductsByReviews []ProductReviewCount `json:"top_products_by_reviews"`
	ProductsByCategory   []CategoryCount      `json:"products_by_category"`
	TopBrands            []BrandProductCount  `json:"top_brands"`
	RecentProducts       []RecentProduct      `json:"recent_products"`
}

type ProductReviewCount struct {
	ProductName   string  `json:"product_name"`
	ReviewCount   int64   `json:"review_count"`
	AverageRating float64 `json:"average_rating"`
}

type CategoryCount struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}

type BrandProductCount struct {
	BrandName    string `json:"brand_name"`
	ProductCount int64  `json:"product_count"`
}

type RecentProduct struct {
	Name      string    `json:"name"`
	Brand     string    `json:"brand"`
	CreatedAt time.Time `json:"created_at"`
}

// GetDashboardStats returns overall dashboard statistics
// @Summary Get dashboard statistics
// @Description Get overview statistics for admin dashboard
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response "Dashboard statistics retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/dashboard/stats [get]
func (ctrl *dashboardController) GetDashboardStats(c *fiber.Ctx) error {
	var stats DashboardStats

	// Get total users
	ctrl.db.Table("users").Where("deleted_at IS NULL").Count(&stats.TotalUsers)

	// Get total products
	ctrl.db.Table("products").Count(&stats.TotalProducts)

	// Get total brands
	ctrl.db.Table("brands").Count(&stats.TotalBrands)

	// Get total reviews
	ctrl.db.Table("reviews").Count(&stats.TotalReviews)

	// Get active users (users who logged in in last 24 hours)
	yesterday := time.Now().AddDate(0, 0, -1)
	ctrl.db.Table("users").Where("last_login > ? AND deleted_at IS NULL", yesterday).Count(&stats.ActiveUsers)

	// Get new users today
	today := time.Now().Truncate(24 * time.Hour)
	ctrl.db.Table("users").Where("created_at >= ? AND deleted_at IS NULL", today).Count(&stats.NewUsersToday)

	// Mock revenue data (you can implement actual revenue calculation)
	stats.RevenueThisMonth = 15420.50

	// Mock AI requests (you can implement actual AI request tracking)
	stats.AIRequestsToday = 127

	return utils.SendResponse(c, http.StatusOK, "Dashboard statistics retrieved successfully", stats)
}

// GetUserAnalytics returns user analytics data
// @Summary Get user analytics
// @Description Get detailed user analytics for admin dashboard
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response "User analytics retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/dashboard/user-analytics [get]
func (ctrl *dashboardController) GetUserAnalytics(c *fiber.Ctx) error {
	var analytics UserAnalytics

	// Get new users for last 7 days
	analytics.NewUsersLastWeek = make([]int, 7)
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Truncate(24 * time.Hour)
		nextDate := date.Add(24 * time.Hour)

		var count int64
		ctrl.db.Table("users").
			Where("created_at >= ? AND created_at < ? AND deleted_at IS NULL", date, nextDate).
			Count(&count)
		analytics.NewUsersLastWeek[6-i] = int(count)
	}

	// Get users by role
	var roleCounts []RoleCount
	ctrl.db.Table("users").
		Select("COALESCE(role, 'user') as role, COUNT(*) as count").
		Where("deleted_at IS NULL").
		Group("role").
		Scan(&roleCounts)
	analytics.UsersByRole = roleCounts

	// Get active users in last hour
	lastHour := time.Now().Add(-1 * time.Hour)
	ctrl.db.Table("users").
		Where("last_login > ? AND deleted_at IS NULL", lastHour).
		Count(&analytics.ActiveUsersLastHour)

	// Get top users by favorites count
	var topUsers []UserFavoriteCount
	ctrl.db.Table("users").
		Select("users.username, COUNT(favorites.id) as favorite_count").
		Joins("LEFT JOIN favorites ON users.id = favorites.user_id").
		Where("users.deleted_at IS NULL").
		Group("users.id, users.username").
		Order("favorite_count DESC").
		Limit(10).
		Scan(&topUsers)
	analytics.TopUsersByFavorites = topUsers

	return utils.SendResponse(c, http.StatusOK, "User analytics retrieved successfully", analytics)
}

// GetProductAnalytics returns product analytics data
// @Summary Get product analytics
// @Description Get detailed product analytics for admin dashboard
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response "Product analytics retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/dashboard/product-analytics [get]
func (ctrl *dashboardController) GetProductAnalytics(c *fiber.Ctx) error {
	var analytics ProductAnalytics

	// Get top products by review count and rating
	var topProducts []ProductReviewCount
	ctrl.db.Table("products").
		Select("products.name as product_name, COUNT(reviews.id) as review_count, AVG(reviews.rating) as average_rating").
		Joins("LEFT JOIN reviews ON products.id = reviews.product_id").
		Group("products.id, products.name").
		Order("review_count DESC, average_rating DESC").
		Limit(10).
		Scan(&topProducts)
	analytics.TopProductsByReviews = topProducts

	// Get products by category
	var categoryCounts []CategoryCount
	ctrl.db.Table("products").
		Select("category, COUNT(*) as count").
		Group("category").
		Order("count DESC").
		Scan(&categoryCounts)
	analytics.ProductsByCategory = categoryCounts

	// Get top brands by product count
	var topBrands []BrandProductCount
	ctrl.db.Table("brands").
		Select("brands.name as brand_name, COUNT(products.id) as product_count").
		Joins("LEFT JOIN products ON brands.id = products.brand_id").
		Group("brands.id, brands.name").
		Order("product_count DESC").
		Limit(10).
		Scan(&topBrands)
	analytics.TopBrands = topBrands

	// Get recent products (last 10)
	var recentProducts []RecentProduct
	ctrl.db.Table("products").
		Select("products.name, brands.name as brand, products.created_at").
		Joins("LEFT JOIN brands ON products.brand_id = brands.id").
		Order("products.created_at DESC").
		Limit(10).
		Scan(&recentProducts)
	analytics.RecentProducts = recentProducts

	return utils.SendResponse(c, http.StatusOK, "Product analytics retrieved successfully", analytics)
}

// GetRevenueAnalytics returns revenue analytics (mock data for now)
// @Summary Get revenue analytics
// @Description Get revenue analytics for admin dashboard
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response "Revenue analytics retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /admin/dashboard/revenue-analytics [get]
func (ctrl *dashboardController) GetRevenueAnalytics(c *fiber.Ctx) error {
	// Mock revenue data - in real app, you'd calculate from orders/payments table
	revenueData := fiber.Map{
		"monthly_revenue": []float64{12340.50, 15670.25, 18920.75, 22150.00, 19875.50, 25340.25},
		"revenue_by_category": []fiber.Map{
			{"category": "Clothing", "revenue": 15420.75},
			{"category": "Shoes", "revenue": 8750.25},
			{"category": "Accessories", "revenue": 3290.50},
			{"category": "Bags", "revenue": 2180.00},
		},
		"total_revenue":  154250.75,
		"revenue_growth": 12.5, // percentage
	}

	return utils.SendResponse(c, http.StatusOK, "Revenue analytics retrieved successfully", revenueData)
}
