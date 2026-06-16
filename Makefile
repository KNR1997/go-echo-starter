.PHONY: run build migrate-up migrate-down migrate-create migrate-status migrate-reset test deps
include .env
export

# Database connection string for PostgreSQL
# Format: postgres://user:password@host:port/dbname?sslmode=disable
DB_STRING := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

run:
	air

build:
	go build -o bin/api ./cmd/api/main.go

# PostgreSQL migration commands
migrate-up:
	goose -dir migrations postgres "$(DB_STRING)" up

migrate-down:
	goose -dir migrations postgres "$(DB_STRING)" down

migrate-down-to:
	@read -p "Enter target version: " version; \
	goose -dir migrations postgres "$(DB_STRING)" down-to $$version

migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir migrations create $$name sql

migrate-status:
	goose -dir migrations postgres "$(DB_STRING)" status

migrate-reset:
	goose -dir migrations postgres "$(DB_STRING)" reset

migrate-version:
	goose -dir migrations postgres "$(DB_STRING)" version

test:
	go test -v ./...

deps:
	go mod tidy
	go mod download

# Help command to show available commands
help:
	@echo "Available commands:"
	@echo "  make run                 - Run the application with air (hot reload)"
	@echo "  make build              - Build the binary"
	@echo "  make migrate-up         - Apply all available migrations"
	@echo "  make migrate-down       - Roll back the last migration"
	@echo "  make migrate-down-to    - Roll back to a specific version"
	@echo "  make migrate-create     - Create a new migration file"
	@echo "  make migrate-status     - Check migration status"
	@echo "  make migrate-reset      - Roll back all migrations"
	@echo "  make migrate-version    - Show current migration version"
	@echo "  make test               - Run tests"
	@echo "  make deps               - Download dependencies"