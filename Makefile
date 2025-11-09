.PHONY: help build run test clean install dev

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install dependencies
	go mod download
	go mod tidy

build: ## Build the application
	go build -o golang-starter-kit main.go

run: ## Run the application
	go run main.go

dev: ## Run with hot reload (requires air)
	air

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean: ## Clean build artifacts
	rm -f golang-starter-kit
	rm -f coverage.out
	go clean

format: ## Format code
	go fmt ./...

lint: ## Run linter (requires golangci-lint)
	golangci-lint run

deps-update: ## Update dependencies
	go get -u ./...
	go mod tidy

docker-build: ## Build Docker image
	docker build -t golang-starter-kit .

docker-run: ## Run Docker container
	docker run -p 3000:3000 --env-file .env golang-starter-kit



