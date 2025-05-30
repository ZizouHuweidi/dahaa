.PHONY: up down logs backend frontend restart clean test-backend test-frontend install-deps format-frontend lint-frontend

# Docker commands
up:
	docker-compose up --build

down:
	docker-compose down

logs:
	docker-compose logs -f

restart:
	docker-compose down && docker-compose up --build

clean:
	docker-compose down -v
	docker system prune -f

# Development commands
backend:
	cd backend && go run cmd/api/main.go

frontend:
	cd frontend && npm start

# Testing commands
test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && npm test

# Dependency installation
install-deps:
	cd backend && go mod download
	cd frontend && npm install

# Database commands
db-migrate:
	cd backend && go run cmd/migrate/main.go up

db-rollback:
	cd backend && go run cmd/migrate/main.go down

# Linting and formatting
lint-backend:
	cd backend && golangci-lint run

lint-frontend:
	cd frontend && npm run lint

format-backend:
	cd backend && go fmt ./...

format-frontend:
	cd frontend && npx prettier --write .

# Code generation
generate:
	cd backend && go generate ./...

# Help command
help:
	@echo "Available commands:"
	@echo "  make up              - Start all services (backend, db, redis)"
	@echo "  make down            - Stop all services"
	@echo "  make logs            - View logs from all services"
	@echo "  make restart         - Restart all services"
	@echo "  make clean           - Remove all containers and volumes"
	@echo ""
	@echo "Development:"
	@echo "  make backend         - Start backend service (without Docker)"
	@echo "  make frontend        - Start frontend development server"
	@echo ""
	@echo "Testing:"
	@echo "  make test-backend    - Run backend tests"
	@echo "  make test-frontend   - Run frontend tests"
	@echo ""
	@echo "Dependencies:"
	@echo "  make install-deps    - Install all dependencies"
	@echo ""
	@echo "Database:"
	@echo "  make db-migrate      - Run database migrations"
	@echo "  make db-rollback     - Rollback database migrations"
	@echo ""
	@echo "Code Quality:"
	@echo "  make lint-backend    - Run backend linter"
	@echo "  make lint-frontend   - Run frontend linter"
	@echo "  make format-backend  - Format backend code"
	@echo "  make format-frontend - Format frontend code"
	@echo "  make generate        - Generate backend code"
