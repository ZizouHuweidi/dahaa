set dotenv-load := true

db_url := env_var_or_default('DATABASE_URL', 'postgres://' + env_var_or_default('POSTGRES_USER', 'postgres') + ':' + env_var_or_default('POSTGRES_PASSWORD', 'postgres') + '@' + env_var_or_default('POSTGRES_HOST', 'localhost') + ':' + env_var_or_default('POSTGRES_PORT', '5432') + '/' + env_var_or_default('POSTGRES_DB', 'dahaa') + '?sslmode=disable')

default:
    just --list

tools:
    go install github.com/pressly/goose/v3/cmd/goose@latest
    go install github.com/air-verse/air@latest
    go install mvdan.cc/gofumpt@latest

dev:
    air

build:
    go build -o bin/server ./cmd/api

test:
    go test -race ./...

fmt:
    gofumpt -w .

compose-up:
    podman compose -f compose.yml up --build

compose-down:
    podman compose -f compose.yml down

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
