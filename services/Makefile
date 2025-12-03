.PHONY: help db-up db-down db-reset migrate-up migrate-down migrate-create migrate-status db-shell

# Load environment variables
include .env
export

help: ## Show this help message
	@echo "Database Management Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

db-up: ## Start PostgreSQL container
	cd ../deploy/compose && docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@timeout /t 10 /nobreak > nul
	@echo "PostgreSQL is running on localhost:5432"

db-down: ## Stop PostgreSQL container
	cd ../deploy/compose && docker-compose down postgres

db-reset: ## Reset database (delete volume and restart)
	cd ../deploy/compose && docker-compose down postgres
	docker volume rm compose_postgres_data || true
	cd ../deploy/compose && docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@timeout /t 10 /nobreak > nul
	$(MAKE) migrate-up

migrate-up: ## Run all pending migrations
	migrate -path migrations -database "$(DATABASE_URL)" up
	@echo "Migrations completed!"

migrate-down: ## Rollback last migration
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-create: ## Create a new migration (use: make migrate-create NAME=add_new_table)
	@if not defined NAME (echo "Usage: make migrate-create NAME=your_migration_name" && exit 1)
	migrate create -ext sql -dir migrations -seq $(NAME)

migrate-status: ## Show current migration version
	migrate -path migrations -database "$(DATABASE_URL)" version

migrate-force: ## Force migration version (use: make migrate-force VERSION=1)
	@if not defined VERSION (echo "Usage: make migrate-force VERSION=1" && exit 1)
	migrate -path migrations -database "$(DATABASE_URL)" force $(VERSION)

db-shell: ## Connect to PostgreSQL shell
	docker exec -it trading-postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

db-verify: ## Verify database tables
	docker exec -it trading-postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "\dt"

db-logs: ## Show PostgreSQL logs
	docker logs trading-postgres -f
