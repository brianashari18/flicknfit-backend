# FlickNFit Backend API

A modern, scalable REST API backend for FlickNFit application built with Go, Fiber framework, and GORM.

## 🚀 Features

- **Clean Architecture**: Organized with proper separation of concerns
- **Authentication & Authorization**: JWT-based auth with role-based access control
- **Security**: Rate limiting, CORS, Helmet, input validation
- **Error Handling**: Comprehensive error handling with structured responses
- **Validation**: Custom validators with user-friendly error messages
- **Database**: MySQL with GORM ORM and auto-migration
- **Email**: SMTP integration for OTP and notifications
- **Logging**: Structured logging with request tracking
- **Documentation**: Well-documented code and API endpoints

## 📁 Project Structure

```
flicknfit_backend/
├── config/           # Configuration management
├── constants/        # Application constants
├── controllers/      # HTTP handlers
├── database/         # Database connection and migration
├── dtos/            # Data Transfer Objects
├── errors/          # Custom error types
├── middlewares/     # HTTP middlewares
├── models/          # Database models (GORM)
├── repositories/    # Data access layer
├── routes/          # Route definitions
├── services/        # Business logic layer
├── utils/           # Utility functions
├── validators/      # Custom validation rules
├── main.go          # Application entry point
├── .env.example     # Environment variables template
└── README.md        # This file
```

## 🛠️ Technology Stack

- **Language**: Go 1.21+
- **Framework**: Fiber v2 (Express-inspired web framework)
- **Database**: MySQL with GORM ORM
- **Authentication**: JWT (JSON Web Tokens)
- **Validation**: go-playground/validator
- **Email**: SMTP with HTML templates
- **Security**: Rate limiting, CORS, Helmet middleware
- **Configuration**: Environment variables with godotenv

## ⚡ Quick Start

### Prerequisites

- Go 1.21 or higher
- MySQL 5.7 or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd flicknfit_backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Configure your `.env` file**
   ```env
   # Application
   APP_PORT=8000

   # Database
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_HOST=localhost
   DB_PORT=3306
   DB_NAME=flicknfit_backend

   # JWT
   JWT_SECRET_KEY=your_very_secure_jwt_secret_key_here

   # SMTP
   SMTP_HOST=smtp.gmail.com
   SMTP_USER=your_email@gmail.com
   SMTP_PASSWORD=your_app_password
   SMTP_PORT=587
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8000`

## 📚 API Documentation

### Base URL
```
http://localhost:8000/api/v1
```

### Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/users/register` | Register new user |
| POST | `/users/login` | User login |
| POST | `/users/forgot-password` | Request password reset |
| POST | `/users/verify-otp` | Verify OTP |
| POST | `/users/reset-password` | Reset password |
| POST | `/users/refresh-token` | Refresh JWT token |

### User Endpoints (Protected)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/users/logout` | User logout | ✅ |
| GET | `/users/me` | Get current user | ✅ |
| PATCH | `/users/edit-profile` | Update profile | ✅ |

### Admin User Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/admin/users` | Create user | ✅ Admin |
| GET | `/admin/users` | List all users | ✅ Admin |
| GET | `/admin/users/:id` | Get user by ID | ✅ Admin |
| PUT | `/admin/users/:id` | Update user | ✅ Admin |
| DELETE | `/admin/users/:id` | Delete user | ✅ Admin |

### Product Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/products` | List products | ❌ |
| GET | `/products/:id` | Get product | ❌ |
| GET | `/products/search` | Search products | ❌ |
| GET | `/products/filter` | Filter products | ❌ |
| GET | `/products/:id/reviews` | Get reviews | ❌ |
| POST | `/products/:id/reviews` | Create review | ✅ |

### Admin Product Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/admin/products` | Create product | ✅ Admin |
| GET | `/admin/products` | List all products | ✅ Admin |
| GET | `/admin/products/:id` | Get product | ✅ Admin |
| PUT | `/admin/products/:id` | Update product | ✅ Admin |
| DELETE | `/admin/products/:id` | Delete product | ✅ Admin |

### Brand Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/brands` | List brands | ❌ |
| GET | `/brands/:id` | Get brand | ❌ |

### Admin Brand Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/admin/brands` | Create brand | ✅ Admin |
| PUT | `/admin/brands/:id` | Update brand | ✅ Admin |
| DELETE | `/admin/brands/:id` | Delete brand | ✅ Admin |

### Shopping Cart Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/cart` | Get user cart | ✅ |
| POST | `/cart` | Add to cart | ✅ |
| PUT | `/cart/:itemId` | Update cart item | ✅ |
| DELETE | `/cart/:itemId` | Remove from cart | ✅ |

## 🔒 Security Features

### Rate Limiting
- **Global**: 100 requests per minute per IP
- **API Endpoints**: 50 requests per minute per IP  
- **Auth Endpoints**: 10 requests per minute per IP

### Security Headers
- Helmet middleware for security headers
- CORS protection with configurable origins
- Request ID tracking

### Input Validation
- Custom validators for passwords, usernames, phone numbers
- Strict JSON parsing with unknown field detection
- Comprehensive validation error messages

### Authentication
- JWT tokens with configurable expiry
- Role-based access control (user/admin)
- Secure password hashing

## 🗃️ Database Schema

The application uses GORM for database operations with auto-migration enabled. Key models include:

- **User**: User accounts with roles
- **Product**: Product catalog
- **Brand**: Product brands
- **ShoppingCart**: User shopping carts
- **Review**: Product reviews
- **ProductCategory**: Product categorization

## 📝 Development

### Code Style
- Follow Go conventions and best practices
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions small and focused

### Testing
```bash
go test ./...
```

### Building for Production
```bash
go build -o flicknfit-api .
```

## 🚀 Deployment

### Environment Variables for Production
Make sure to set secure values for:
- `JWT_SECRET_KEY`: Use a strong, random secret
- `DB_PASSWORD`: Use a strong database password
- `SMTP_PASSWORD`: Use app-specific password for email

### Docker Support
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you have any questions or need help with setup, please create an issue in the repository.

---

**Built with ❤️ using Go and Fiber**
