# FlickNFit Backend - OpenAPI/Swagger Documentation

## 📚 **API Documentation Implementation Complete**

### **✅ What Has Been Implemented**

#### **1. OpenAPI/Swagger Setup**
- ✅ Added Swagger dependencies (swaggo/swag, fiber-swagger, swaggo/files)
- ✅ Configured main.go with comprehensive API metadata
- ✅ Added Swagger UI route at `/swagger/*`
- ✅ Generated complete API documentation

#### **2. API Metadata Configuration**
```go
// @title FlickNFit API
// @version 1.0
// @description FlickNFit backend API - Fashion recommendation and wardrobe management system
// @host localhost:8000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

#### **3. Controller Annotations Added**

**Authentication Endpoints:**
- ✅ `POST /auth/register` - User registration
- ✅ `POST /auth/login` - User authentication with JWT tokens

**Admin - User Management:**
- ✅ `POST /admin/users` - Create user (Admin only)

**Product Endpoints:**
- ✅ `GET /products` - Get all products (public)
- ✅ `POST /admin/products` - Create product (Admin only)

**Favorites Management:**
- ✅ `GET /favorites` - Get user's favorite products
- ✅ `POST /favorites` - Add product to favorites

**Reviews System:**
- ✅ `GET /products/{productId}/reviews` - Get product reviews with pagination
- ✅ `POST /reviews` - Create product review

**Wardrobe Management:**
- ✅ `GET /wardrobe` - Get user's wardrobe organized by category
- ✅ `POST /wardrobe` - Add item to wardrobe

**Shopping Cart:**
- ✅ `GET /cart` - Get user's shopping cart

#### **4. Documentation Features**

**Request/Response Models:**
- ✅ All DTOs properly documented with JSON schema
- ✅ Request body validation specifications
- ✅ Response structure definitions

**Security Documentation:**
- ✅ Bearer Token authentication documented
- ✅ Protected endpoints marked with `@Security BearerAuth`

**Error Responses:**
- ✅ Comprehensive error response documentation
- ✅ HTTP status codes with descriptions
- ✅ Error message structures

**Pagination Support:**
- ✅ Query parameters documented (page, limit)
- ✅ Pagination response structure

### **🚀 How to Access Documentation**

#### **1. Start the Server**
```bash
go run main.go
```

#### **2. Access Swagger UI**
Open your browser and navigate to:
```
http://localhost:8000/swagger/index.html
```

#### **3. Generated Documentation Files**
- `docs/swagger.json` - OpenAPI 3.0 JSON specification
- `docs/swagger.yaml` - OpenAPI 3.0 YAML specification  
- `docs/docs.go` - Generated Go documentation

### **📋 Available API Categories**

1. **Authentication** - User registration, login, JWT token management
2. **Admin - User Management** - Administrative user operations
3. **Admin - Product Management** - Product CRUD operations
4. **Products** - Public product browsing and information
5. **Favorites** - User favorite product management
6. **Reviews** - Product review system with ratings
7. **Wardrobe** - Personal wardrobe management
8. **Shopping Cart** - Shopping cart operations

### **🔐 Security Implementation**

**Bearer Token Authentication:**
```json
{
  "Authorization": "Bearer <your-jwt-token>"
}
```

**Protected Endpoints:**
- All user-specific operations require authentication
- Admin operations require admin role verification
- Public endpoints (products, reviews viewing) don't require authentication

### **📊 API Response Format**

**Standard Response Structure:**
```json
{
  "status": "OK",
  "message": "Operation successful",
  "data": {
    // Actual response data
  }
}
```

**Error Response Structure:**
```json
{
  "status": "Bad Request",
  "message": "Validation failed: field is required",
  "data": null
}
```

### **🛠️ Development Commands**

**Regenerate Documentation:**
```bash
# After adding new API annotations
go run github.com/swaggo/swag/cmd/swag init
```

**Install Swagger CLI (optional):**
```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init  # If installed globally
```

### **📝 Adding New API Documentation**

**Controller Method Example:**
```go
// @Summary Brief description
// @Description Detailed description
// @Tags Category Name
// @Accept json
// @Produce json
// @Security BearerAuth (if authentication required)
// @Param paramName body/path/query dtos.RequestDTO true "Description"
// @Success 200 {object} utils.Response{data=dtos.ResponseDTO} "Success description"
// @Failure 400 {object} utils.Response "Error description"
// @Router /endpoint/path [post]
func (ctrl *controller) Method(c *fiber.Ctx) error {
    // Implementation
}
```

### **✨ Key Features Implemented**

1. **Interactive Documentation** - Test API endpoints directly from browser
2. **Model Schemas** - Auto-generated from Go structs
3. **Authentication Integration** - Bearer token support with "Authorize" button
4. **Request Examples** - Sample JSON payloads for all endpoints
5. **Response Examples** - Expected response formats
6. **Error Documentation** - Comprehensive error code coverage
7. **Search & Filter** - Find endpoints quickly
8. **Download Options** - Export as JSON/YAML

### **🎯 Benefits**

- **Developer Experience** - Easy API exploration and testing
- **Frontend Integration** - Clear contract definition for frontend teams
- **API Testing** - Built-in testing interface
- **Documentation Maintenance** - Auto-generated from code annotations
- **Standardization** - Consistent API documentation format
- **Onboarding** - New developers can quickly understand API structure

The FlickNFit API now has comprehensive, professional-grade OpenAPI documentation that makes it easy for developers to understand, test, and integrate with the backend services.
