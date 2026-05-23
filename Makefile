.PHONY: run build migrate-up migrate-down migrate-create test
include .env
export

# Database connection string for goose
DB_STRING := root:$(DB_PASSWORD)@tcp(localhost:3306)/$(DB_NAME)?parseTime=true

run:
	air

build:
	go build -o bin/api ./cmd/api/main.go

migrate-up:
	goose -dir migrations mysql "$(DB_STRING)" up

migrate-down:
	goose -dir migrations mysql "$(DB_STRING)" down

migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir migrations create $$name sql

test:
	go test -v ./...

deps:
	go mod tidy
	go mod download