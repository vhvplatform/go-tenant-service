.PHONY: help build test lint clean run docker-build docker-push proto

# Variables
SERVICE_NAME := tenant-service
DOCKER_REGISTRY ?= ghcr.io/vhvplatform
VERSION ?= $(shell git describe --tags --always --dirty)
GO_VERSION := 1.25.5

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the service
	@echo "Building $(SERVICE_NAME)..."
	@go build -o bin/$(SERVICE_NAME) ./cmd/main.go
	@echo "Build complete!"

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linters
	@echo "Running linters..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -s -w .

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/ dist/ coverage.* *.out
	@go clean -testcache
	@echo "Clean complete!"

run: ## Run the service locally
	@echo "Running $(SERVICE_NAME)..."
	@go run ./cmd/main.go

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

proto: ## Generate protobuf files (if applicable)
	@if [ -d "proto" ]; then \
		echo "Generating protobuf files..."; \
		protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			proto/*.proto; \
	fi

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_REGISTRY)/$(SERVICE_NAME):$(VERSION) .
	@docker tag $(DOCKER_REGISTRY)/$(SERVICE_NAME):$(VERSION) $(DOCKER_REGISTRY)/$(SERVICE_NAME):latest
	@echo "Docker image built: $(DOCKER_REGISTRY)/$(SERVICE_NAME):$(VERSION)"

docker-push: docker-build ## Push Docker image
	@echo "Pushing Docker image..."
	@docker push $(DOCKER_REGISTRY)/$(SERVICE_NAME):$(VERSION)
	@docker push $(DOCKER_REGISTRY)/$(SERVICE_NAME):latest
	@echo "Docker image pushed!"

docker-run: ## Run Docker container locally
	@echo "Running Docker container..."
	@docker run --rm -p 8080:8080 -p 50051:50051 \
		--name $(SERVICE_NAME) \
		$(DOCKER_REGISTRY)/$(SERVICE_NAME):latest

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@if [ -d "proto" ]; then \
		go install google.golang.org/protobuf/cmd/protoc-gen-go@latest; \
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest; \
	fi
	@echo "Tools installed!"

# Tenant-specific targets
deploy-tenant: ## Deploy tenant-specific instance (usage: make deploy-tenant TENANT_ID=<id>)
	@if [ -z "$(TENANT_ID)" ]; then \
		echo "Error: TENANT_ID is required. Usage: make deploy-tenant TENANT_ID=<id>"; \
		exit 1; \
	fi
	@echo "Deploying tenant-specific instance for tenant: $(TENANT_ID)"
	@docker build -t $(DOCKER_REGISTRY)/$(SERVICE_NAME):$(TENANT_ID)-$(VERSION) \
		--build-arg TENANT_ID=$(TENANT_ID) .
	@docker push $(DOCKER_REGISTRY)/$(SERVICE_NAME):$(TENANT_ID)-$(VERSION)
	@echo "Tenant deployment complete!"

migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	@echo "Note: Add your migration tool command here"

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	@echo "Note: Add your migration tool command here"

perf-test: ## Run performance tests
	@echo "Running performance tests..."
	@if command -v hey > /dev/null; then \
		hey -n 1000 -c 10 http://localhost:8083/health; \
	else \
		echo "hey is not installed. Install with: go install github.com/rakyll/hey@latest"; \
	fi

security-scan: ## Run security scanning with gosec
	@echo "Running security scan..."
	@if command -v gosec > /dev/null; then \
		gosec -fmt=json -out=security-report.json ./...; \
		echo "Security report generated: security-report.json"; \
	else \
		echo "gosec is not installed. Install with: go install github.com/securego/gosec/v2/cmd/gosec@latest"; \
	fi

bench: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

install-dev-tools: install-tools ## Install additional development tools
	@echo "Installing additional development tools..."
	@go install github.com/rakyll/hey@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "Additional tools installed!"

docker-build-optimized: ## Build optimized Docker image with build cache
	@echo "Building optimized Docker image..."
	@docker build --target builder --tag $(SERVICE_NAME)-builder:$(VERSION) .
	@docker build --cache-from=$(SERVICE_NAME)-builder:$(VERSION) \
		-t $(DOCKER_REGISTRY)/$(SERVICE_NAME):$(VERSION) .
	@docker tag $(DOCKER_REGISTRY)/$(SERVICE_NAME):$(VERSION) $(DOCKER_REGISTRY)/$(SERVICE_NAME):latest
	@echo "Optimized Docker image built!"

docker-run-dev: ## Run Docker container in development mode with volume mounts
	@echo "Running Docker container in development mode..."
	@docker run --rm -it \
		-p 8080:8080 -p 50051:50051 \
		-v $(PWD):/app \
		-e LOG_LEVEL=debug \
		--name $(SERVICE_NAME)-dev \
		$(DOCKER_REGISTRY)/$(SERVICE_NAME):latest

.DEFAULT_GOAL := help
