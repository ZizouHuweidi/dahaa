# Variables
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= dahaa
MIGRATION_DIR ?= migrations

# Docker commands
.PHONY: docker-up
docker-up:
	docker-compose up -d

.PHONY: docker-down
docker-down:
	docker-compose down

# Database commands
.PHONY: db-create
db-create:
	PGPASSWORD=$(DB_PASSWORD) createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)

.PHONY: db-drop
db-drop:
	PGPASSWORD=$(DB_PASSWORD) dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)

# Migration commands
.PHONY: migrate-create
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $$name

.PHONY: migrate-up
migrate-up:
	migrate -path $(MIGRATION_DIR) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	migrate -path $(MIGRATION_DIR) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

.PHONY: migrate-force
migrate-force:
	@read -p "Enter version: " version; \
	migrate -path $(MIGRATION_DIR) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" force $$version

# Development commands
.PHONY: run
run:
	go run cmd/server/main.go

.PHONY: build
build:
	go build -o bin/server cmd/server/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run

# Setup commands
.PHONY: setup
setup: docker-up db-create migrate-up

.PHONY: clean
clean: docker-down db-drop

# Help command
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make setup              - Set up the development environment (network, database, migrations)"
	@echo "  make clean             - Clean up the development environment"
	@echo "  make run               - Run the development server"
	@echo "  make build             - Build the server"
	@echo "  make test              - Run tests"
	@echo "  make lint              - Run linting"
	@echo "  make migrate-create    - Create a new migration"
	@echo "  make migrate-up        - Run migrations up"
	@echo "  make migrate-down      - Run migrations down"
	@echo "  make migrate-force     - Force migration version" 