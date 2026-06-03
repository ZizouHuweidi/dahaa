set dotenv-load := true

db_url := env_var_or_default('DATABASE_URL', 'postgres://' + env_var_or_default('POSTGRES_USER', 'postgres') + ':' + env_var_or_default('POSTGRES_PASSWORD', 'postgres') + '@' + env_var_or_default('POSTGRES_HOST', 'localhost') + ':' + env_var_or_default('POSTGRES_PORT', '5432') + '/' + env_var_or_default('POSTGRES_DB', 'dahaa') + '?sslmode=disable')
podman := env_var_or_default('PODMAN', 'podman')
compose := env_var_or_default('PODMAN_COMPOSE', podman + ' compose')
compose_file := env_var_or_default('COMPOSE_FILE', 'compose.yml')
api_image := env_var_or_default('API_IMAGE', 'dahaa-api:dev')
web_image := env_var_or_default('WEB_IMAGE', 'dahaa-web:dev')
sqlc_image := env_var_or_default('SQLC_IMAGE', 'docker.io/sqlc/sqlc:1.31.1')

default:
    just --list

# Install local development tools.
tools:
    go install github.com/pressly/goose/v3/cmd/goose@latest
    go install github.com/air-verse/air@latest
    go install mvdan.cc/gofumpt@latest

# Generate sqlc code using the sqlc container image.
sqlc-generate:
    {{podman}} run --rm -v "$PWD:/src:Z" -w /src {{sqlc_image}} generate

# Vet sqlc queries using the sqlc container image.
sqlc-vet:
    {{podman}} run --rm -v "$PWD:/src:Z" -w /src {{sqlc_image}} vet

# Generate and vet database query code.
db-check: sqlc-generate sqlc-vet

# Run the API locally with air. Use `just up-deps` first for Postgres and Redis.
dev:
    air

# Build the API binary locally.
build:
    go build -o bin/server ./cmd/api

# Run Go tests with the race detector.
test:
    go test -race ./...

# Format Go code.
fmt:
    gofumpt -w .

# Check frontend formatting/linting with Biome.
frontend-check:
    npm --prefix frontend run check

# Format frontend files with Biome.
frontend-format:
    npm --prefix frontend run format:write

# Start the full stack with Podman Compose.
up:
    {{compose}} -f {{compose_file}} up --build

# Start only infrastructure dependencies.
up-deps:
    {{compose}} -f {{compose_file}} up -d postgres redis

# Start the full stack in the background.
up-detached:
    {{compose}} -f {{compose_file}} up -d --build

# Stop and remove the Podman Compose stack.
down:
    {{compose}} -f {{compose_file}} down

# Stop and remove the stack, including named volumes.
down-volumes:
    {{compose}} -f {{compose_file}} down --volumes

# Show Podman Compose service status.
ps:
    {{compose}} -f {{compose_file}} ps

# Follow logs for all services, or pass a service name.
logs service="":
    {{compose}} -f {{compose_file}} logs -f {{service}}

# Restart one service, for example `just restart api`.
restart service:
    {{compose}} -f {{compose_file}} restart {{service}}

# Open a shell in a running service container.
shell service="api":
    {{compose}} -f {{compose_file}} exec {{service}} sh

# Build the API image with Podman.
image-api:
    {{podman}} build -f Containerfile -t {{api_image}} .

# Build the web image with Podman.
image-web:
    {{podman}} build -f frontend/Containerfile -t {{web_image}} frontend

# Build both application images with Podman.
images: image-api image-web

# Run the API image directly. Use `just up-deps` first.
run-api: image-api
    env_file=""; [ ! -f .env ] || env_file="--env-file .env"; {{podman}} run --rm --name dahaa-api --network host $env_file {{api_image}}

# Remove unused Podman data.
prune:
    {{podman}} system prune

# Backwards-compatible aliases.
compose-up: up
compose-down: down

migrate-up:
    goose -dir migrations postgres "{{db_url}}" up

migrate-down:
    goose -dir migrations postgres "{{db_url}}" down

migrate-status:
    goose -dir migrations postgres "{{db_url}}" status

migration name:
    goose -dir migrations create {{name}} sql

frontend-dev:
    npm --prefix frontend run dev

frontend-build:
    npm --prefix frontend run build

frontend-typecheck:
    npm --prefix frontend run typecheck

# Verify frontend typecheck, Biome, and production build.
frontend-verify: frontend-typecheck frontend-check frontend-build

# Check API health endpoint.
health:
    curl -fsS http://localhost:8080/health

# Check API readiness endpoint.
ready:
    curl -fsS http://localhost:8080/ready

# Verify backend and frontend locally.
verify: test frontend-verify
