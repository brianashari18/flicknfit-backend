package admin

import (
	"flicknfit_backend/config"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"gorm.io/gorm"
)

// AdminConfig holds admin configuration
type AdminConfig struct {
	DB     *gorm.DB
	Config *config.Config
}

// DashboardData holds dashboard statistics for template rendering
type DashboardData struct {
	TotalUsers       int64   `json:"total_users"`
	TotalProducts    int64   `json:"total_products"`
	TotalBrands      int64   `json:"total_brands"`
	TotalReviews     int64   `json:"total_reviews"`
	ActiveUsers      int64   `json:"active_users"`
	NewUsersToday    int64   `json:"new_users_today"`
	RevenueThisMonth float64 `json:"revenue_this_month"`
	AIRequestsToday  int64   `json:"ai_requests_today"`
}

// SetupAdmin initializes a simple HTML-based admin dashboard for FlickNFit
func SetupAdmin(app *fiber.App, adminConfig *AdminConfig) error {
	log.Printf("Setting up FlickNFit Admin Dashboard...")

	// Create template engine
	engine := html.New("./admin/templates", ".html")
	engine.Reload(true) // for development

	// Setup admin routes
	setupAdminRoutes(app, adminConfig)

	log.Printf("Admin dashboard successfully initialized")
	return nil
}

// setupAdminRoutes creates the admin dashboard routes
func setupAdminRoutes(app *fiber.App, adminConfig *AdminConfig) {
	// Admin dashboard main route
	app.Get("/admin", func(c *fiber.Ctx) error {
		// Get dashboard data
		data := getDashboardData(adminConfig.DB)

		// Render with basic HTML response for now
		html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>FlickNFit Admin Dashboard</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
</head>
<body class="bg-gray-100">
    <div class="min-h-screen">
        <!-- Header -->
        <header class="bg-white shadow">
            <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
                <h1 class="text-3xl font-bold text-gray-900">FlickNFit Admin Dashboard</h1>
            </div>
        </header>

        <!-- Main Content -->
        <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
            <!-- Stats Overview -->
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                <div class="bg-white overflow-hidden shadow rounded-lg">
                    <div class="p-5">
                        <div class="flex items-center">
                            <div class="flex-shrink-0">
                                <div class="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center">
                                    <span class="text-white text-sm font-medium">👥</span>
                                </div>
                            </div>
                            <div class="ml-5 w-0 flex-1">
                                <dl>
                                    <dt class="text-sm font-medium text-gray-500 truncate">Total Users</dt>
                                    <dd class="text-lg font-medium text-gray-900">%d</dd>
                                </dl>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="bg-white overflow-hidden shadow rounded-lg">
                    <div class="p-5">
                        <div class="flex items-center">
                            <div class="flex-shrink-0">
                                <div class="w-8 h-8 bg-green-500 rounded-full flex items-center justify-center">
                                    <span class="text-white text-sm font-medium">📦</span>
                                </div>
                            </div>
                            <div class="ml-5 w-0 flex-1">
                                <dl>
                                    <dt class="text-sm font-medium text-gray-500 truncate">Total Products</dt>
                                    <dd class="text-lg font-medium text-gray-900">%d</dd>
                                </dl>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="bg-white overflow-hidden shadow rounded-lg">
                    <div class="p-5">
                        <div class="flex items-center">
                            <div class="flex-shrink-0">
                                <div class="w-8 h-8 bg-purple-500 rounded-full flex items-center justify-center">
                                    <span class="text-white text-sm font-medium">🏷️</span>
                                </div>
                            </div>
                            <div class="ml-5 w-0 flex-1">
                                <dl>
                                    <dt class="text-sm font-medium text-gray-500 truncate">Total Brands</dt>
                                    <dd class="text-lg font-medium text-gray-900">%d</dd>
                                </dl>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="bg-white overflow-hidden shadow rounded-lg">
                    <div class="p-5">
                        <div class="flex items-center">
                            <div class="flex-shrink-0">
                                <div class="w-8 h-8 bg-red-500 rounded-full flex items-center justify-center">
                                    <span class="text-white text-sm font-medium">⭐</span>
                                </div>
                            </div>
                            <div class="ml-5 w-0 flex-1">
                                <dl>
                                    <dt class="text-sm font-medium text-gray-500 truncate">Total Reviews</dt>
                                    <dd class="text-lg font-medium text-gray-900">%d</dd>
                                </dl>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Management Links -->
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
                <div class="bg-white overflow-hidden shadow rounded-lg">
                    <div class="p-6">
                        <h3 class="text-lg font-medium text-gray-900 mb-4">User Management</h3>
                        <div class="space-y-2">
                            <a href="/admin/users" class="block text-blue-600 hover:text-blue-800">View All Users</a>
                            <a href="/admin/users/create" class="block text-blue-600 hover:text-blue-800">Create New User</a>
                        </div>
                    </div>
                </div>

                <div class="bg-white overflow-hidden shadow rounded-lg">
                    <div class="p-6">
                        <h3 class="text-lg font-medium text-gray-900 mb-4">Product Management</h3>
                        <div class="space-y-2">
                            <a href="/admin/products" class="block text-blue-600 hover:text-blue-800">View All Products</a>
                            <a href="/admin/products/create" class="block text-blue-600 hover:text-blue-800">Create New Product</a>
                        </div>
                    </div>
                </div>

                <div class="bg-white overflow-hidden shadow rounded-lg">
                    <div class="p-6">
                        <h3 class="text-lg font-medium text-gray-900 mb-4">Brand Management</h3>
                        <div class="space-y-2">
                            <a href="/admin/brands" class="block text-blue-600 hover:text-blue-800">View All Brands</a>
                            <a href="/admin/brands/create" class="block text-blue-600 hover:text-blue-800">Create New Brand</a>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Analytics Charts -->
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <div class="bg-white shadow rounded-lg p-6">
                    <h3 class="text-lg font-medium text-gray-900 mb-4">User Growth</h3>
                    <canvas id="userChart" width="400" height="200"></canvas>
                </div>

                <div class="bg-white shadow rounded-lg p-6">
                    <h3 class="text-lg font-medium text-gray-900 mb-4">Product Categories</h3>
                    <canvas id="categoryChart" width="400" height="200"></canvas>
                </div>
            </div>

            <!-- API Information -->
            <div class="mt-8 bg-white shadow rounded-lg p-6">
                <h3 class="text-lg font-medium text-gray-900 mb-4">API Information</h3>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <p class="text-sm text-gray-600 mb-2">API Documentation:</p>
                        <a href="/swagger/index.html" class="text-blue-600 hover:text-blue-800" target="_blank">View Swagger Documentation</a>
                    </div>
                    <div>
                        <p class="text-sm text-gray-600 mb-2">Server Status:</p>
                        <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                            🟢 Online
                        </span>
                    </div>
                </div>
            </div>
        </main>
    </div>

    <script>
        // User Growth Chart
        const userCtx = document.getElementById('userChart').getContext('2d');
        new Chart(userCtx, {
            type: 'line',
            data: {
                labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
                datasets: [{
                    label: 'New Users',
                    data: [12, 19, 3, 5, 2, 3],
                    borderColor: 'rgb(59, 130, 246)',
                    backgroundColor: 'rgba(59, 130, 246, 0.1)',
                    tension: 0.1
                }]
            },
            options: {
                responsive: true,
                scales: {
                    y: {
                        beginAtZero: true
                    }
                }
            }
        });

        // Product Categories Chart
        const categoryCtx = document.getElementById('categoryChart').getContext('2d');
        new Chart(categoryCtx, {
            type: 'doughnut',
            data: {
                labels: ['Clothing', 'Shoes', 'Accessories', 'Bags'],
                datasets: [{
                    data: [45, 25, 20, 10],
                    backgroundColor: [
                        'rgb(59, 130, 246)',
                        'rgb(16, 185, 129)',
                        'rgb(245, 158, 11)',
                        'rgb(239, 68, 68)'
                    ]
                }]
            },
            options: {
                responsive: true,
                plugins: {
                    legend: {
                        position: 'bottom'
                    }
                }
            }
        });
    </script>
</body>
</html>
		`, data.TotalUsers, data.TotalProducts, data.TotalBrands, data.TotalReviews)

		c.Set("Content-Type", "text/html")
		return c.SendString(html)
	})

	// Admin API endpoints for AJAX calls
	app.Get("/admin/api/stats", func(c *fiber.Ctx) error {
		data := getDashboardData(adminConfig.DB)
		return c.JSON(data)
	})
}

// getDashboardData retrieves dashboard statistics from database
func getDashboardData(db *gorm.DB) DashboardData {
	var data DashboardData

	// Get total users
	db.Table("users").Where("deleted_at IS NULL").Count(&data.TotalUsers)

	// Get total products
	db.Table("products").Count(&data.TotalProducts)

	// Get total brands
	db.Table("brands").Count(&data.TotalBrands)

	// Get total reviews
	db.Table("reviews").Count(&data.TotalReviews)

	// Get active users (users who logged in in last 24 hours)
	yesterday := time.Now().AddDate(0, 0, -1)
	db.Table("users").Where("last_login > ? AND deleted_at IS NULL", yesterday).Count(&data.ActiveUsers)

	// Get new users today
	today := time.Now().Truncate(24 * time.Hour)
	db.Table("users").Where("created_at >= ? AND deleted_at IS NULL", today).Count(&data.NewUsersToday)

	// Mock data for now
	data.RevenueThisMonth = 15420.50
	data.AIRequestsToday = 127

	return data
}
