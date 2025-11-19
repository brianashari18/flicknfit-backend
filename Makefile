# FlickNFit Backend Makefile

# Variables
BINARY_NAME=flicknfit-api
DOCKER_IMAGE=flicknfit/backend
VERSION=1.0.0

# Go commands
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: run
run: ## Run the application
	go run main.go

.PHONY: build
build: ## Build the application
	go build -o $(BINARY_NAME) .

.PHONY: build-windows
build-windows: ## Build for Windows
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe .

.PHONY: build-linux
build-linux: ## Build for Linux
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) .

.PHONY: build-mac
build-mac: ## Build for macOS
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME) .

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: test-unit
test-unit: ## Run unit tests only
	go test -v ./tests/unit/...

.PHONY: test-integration
test-integration: ## Run integration tests only
	go test -v ./tests/integration/...

.PHONY: test-watch
test-watch: ## Run tests in watch mode
	@echo "Running tests in watch mode..."
	@while true; do \
		go test -v ./...; \
		sleep 2; \
	done

.PHONY: test-benchmarks
test-benchmarks: ## Run benchmark tests
	go test -bench=. -benchmem ./...

.PHONY: clean
clean: ## Clean build artifacts
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe
	rm -f coverage.out
	rm -f coverage.html

.PHONY: docs-generate
docs-generate: ## Generate Swagger/OpenAPI documentation
	go run github.com/swaggo/swag/cmd/swag init

.PHONY: docs-serve
docs-serve: ## Start server and open Swagger UI
	@echo "Starting server with Swagger documentation..."
	@echo "Swagger UI will be available at: http://localhost:8080/swagger/index.html"
	go run main.go

.PHONY: docs-clean
docs-clean: ## Clean generated documentation files
	rm -f docs/docs.go docs/swagger.json docs/swagger.yaml

.PHONY: deps
deps: ## Download dependencies
	go mod download
	go mod tidy

.PHONY: deps-upgrade
deps-upgrade: ## Upgrade dependencies
	go get -u ./...
	go mod tidy

.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: format
format: ## Format code
	go fmt ./...
	goimports -w .

.PHONY: security
security: ## Run security scan
	gosec ./...

# Docker commands
.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(VERSION) .
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest

.PHONY: docker-run
docker-run: ## Run Docker container
	docker run -p 8000:8000 --env-file .env $(DOCKER_IMAGE):latest

.PHONY: docker-compose-up
docker-compose-up: ## Start all services with docker-compose
	docker-compose up -d

.PHONY: docker-compose-down
docker-compose-down: ## Stop all services
	docker-compose down

.PHONY: docker-compose-logs
docker-compose-logs: ## View docker-compose logs
	docker-compose logs -f

.PHONY: docker-clean
docker-clean: ## Clean Docker images and containers
	docker-compose down -v
	docker rmi $(DOCKER_IMAGE):latest $(DOCKER_IMAGE):$(VERSION) 2>/dev/null || true

# Database commands
.PHONY: db-migrate
db-migrate: ## Run database migrations
	go run main.go migrate

.PHONY: db-seed
db-seed: ## Seed database with sample data
	go run scripts/seed.go

.PHONY: db-reset
db-reset: ## Reset database (WARNING: This will delete all data)
	@echo "Are you sure you want to reset the database? [y/N]" && read ans && [ $${ans:-N} = y ]
	go run scripts/reset.go

# Development commands
.PHONY: dev
dev: ## Run in development mode with hot reload
	air

.PHONY: dev-setup
dev-setup: ## Setup development environment
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env file from template"; fi
	@echo "Please edit .env file with your configuration"
	go mod download
	@echo "Development environment setup complete!"

.PHONY: generate-key
generate-key: ## Generate JWT secret key
	@echo "Generated JWT secret key:"
	@openssl rand -hex 32

# API commands
.PHONY: api-docs
api-docs: ## Generate API documentation
	swag init

.PHONY: api-test
api-test: ## Run API tests
	@echo "Running API integration tests..."
	go test -tags=integration ./tests/...

# Production commands
.PHONY: deploy-staging
deploy-staging: ## Deploy to staging
	@echo "Deploying to staging..."
	# Add your staging deployment commands here

.PHONY: deploy-prod
deploy-prod: ## Deploy to production
	@echo "Deploying to production..."
	# Add your production deployment commands here

.PHONY: backup-db
backup-db: ## Backup production database
	@echo "Creating database backup..."
	# Add your database backup commands here

# Git hooks
.PHONY: install-hooks
install-hooks: ## Install git hooks
	@echo "Installing git hooks..."
	cp scripts/pre-commit .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
	@echo "Git hooks installed successfully"

# Health check
.PHONY: health
health: ## Check application health
	@curl -f http://localhost:8000/health || echo "Application is not running"

.PHONY: all
all: deps format lint test build ## Run all checks and build
