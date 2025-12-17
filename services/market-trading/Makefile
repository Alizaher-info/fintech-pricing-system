.PHONY: help proto build run test clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

proto: ## Generate protobuf code
	@echo "Generating protobuf code..."
	@mkdir -p gen/market/v1
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		-I=../../protos \
		../../protos/market/v1/market.proto
	@echo "✅ Protobuf code generated"

build: ## Build the application
	@echo "Building market-trading service..."
	go build -o bin/market-trading ./cmd/server
	@echo "✅ Build complete"

run: ## Run the application
	@echo "Running market-trading service..."
	go run ./cmd/server/main.go

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf gen/
	@echo "✅ Clean complete"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t market-trading:latest .
	@echo "✅ Docker image built"

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 50052:50052 --env-file .env market-trading:latest

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies downloaded"
