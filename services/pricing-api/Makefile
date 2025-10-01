# Makefile für pricing-api gRPC Service
# Portfolio Finance - Go Microservice

.PHONY: proto clean build run test docker-build docker-run docker-stop help

# Variablen
SERVICE_NAME = pricing-api
PROTO_PATH = ../../protos
MODULE_NAME = pricingapi
DOCKER_COMPOSE_PATH = ../../deploy/compose

# Standard Target
help:
	@echo "=== Pricing API - Make Commands ==="
	@echo "proto        - Generiere Protobuf-Dateien"
	@echo "build        - Baue Go Binary lokal"
	@echo "run          - Starte Service lokal"
	@echo "test         - Führe Tests aus"
	@echo "clean        - Lösche generierte Dateien"
	@echo "docker-build - Baue Docker Image"
	@echo "docker-run   - Starte Docker Container"
	@echo "docker-stop  - Stoppe Docker Container"
	@echo "docker-logs  - Zeige Docker Logs"

# Protobuf Generierung
proto:
	@echo "🔧 Generiere Protobuf-Dateien..."
	@if not exist "gen\pricing\v1" mkdir gen\pricing\v1
	protoc --proto_path=$(PROTO_PATH) \
		--go_out=. --go_opt=module=$(MODULE_NAME) \
		--go-grpc_out=. --go-grpc_opt=module=$(MODULE_NAME) \
		pricing/v1/pricing.proto
	@echo "✅ Protobuf-Dateien erfolgreich generiert!"

# Aufräumen
clean:
	@echo "🧹 Lösche generierte Dateien..."
	@if exist "gen" rmdir /s /q gen
	@if exist "$(SERVICE_NAME).exe" del $(SERVICE_NAME).exe
	@if exist "$(SERVICE_NAME)" del $(SERVICE_NAME)
	@echo "✅ Aufgeräumt!"

# Dependencies laden
deps:
	@echo "📦 Lade Go Dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies geladen!"

# Build binary lokal
build: proto deps
	@echo "🔨 Baue $(SERVICE_NAME)..."
	go build -o $(SERVICE_NAME).exe
	@echo "✅ Build erfolgreich!"

# Tests ausführen
test:
	@echo "🧪 Führe Tests aus..."
	go test ./...
	@echo "✅ Tests abgeschlossen!"

# Service lokal laufen lassen
run: build
	@echo "🚀 Starte $(SERVICE_NAME) lokal..."
	./$(SERVICE_NAME).exe

# Docker Commands
docker-build:
	@echo "🐳 Baue Docker Image..."
	cd $(DOCKER_COMPOSE_PATH) && docker-compose build $(SERVICE_NAME)
	@echo "✅ Docker Image gebaut!"

docker-run:
	@echo "🐳 Starte Docker Container..."
	cd $(DOCKER_COMPOSE_PATH) && docker-compose up $(SERVICE_NAME) -d
	@echo "✅ Docker Container gestartet!"

docker-stop:
	@echo "🐳 Stoppe Docker Container..."
	cd $(DOCKER_COMPOSE_PATH) && docker-compose stop $(SERVICE_NAME)
	@echo "✅ Docker Container gestoppt!"

docker-logs:
	@echo "📋 Docker Logs:"
	cd $(DOCKER_COMPOSE_PATH) && docker-compose logs $(SERVICE_NAME)

# Vollständiger Docker Rebuild
docker-rebuild: clean
	@echo "🔄 Vollständiger Docker Rebuild..."
	cd $(DOCKER_COMPOSE_PATH) && docker-compose build --no-cache $(SERVICE_NAME)
	@echo "✅ Docker Rebuild abgeschlossen!"

# Entwicklungsworkflow
dev: clean proto build run

# Deployment workflow  
deploy: clean docker-build docker-run
	@echo "🎉 Deployment abgeschlossen!"