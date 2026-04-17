.PHONY: help build run dev test clean db-up db-down docker-up docker-down

help:
	@echo "Available commands:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make dev         - Run in development mode with hot reload"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make db-up       - Start PostgreSQL database"
	@echo "  make db-down     - Stop PostgreSQL database"
	@echo "  make docker-up   - Start all services with Docker Compose"
	@echo "  make docker-down - Stop all services"
	@echo "  make migrate     - Run database migrations"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Run linter"

build:
	@echo "Building application..."
	go build -o bin/task-api ./cmd/server

run: build
	@echo "Running application..."
	./bin/task-api

dev:
	@echo "Running in development mode..."
	go run ./cmd/server/main.go

test:
	@echo "Running tests..."
	go test -v -cover ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

db-up:
	@echo "Starting PostgreSQL database..."
	docker run -d --name task_db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=task_management -p 5432:5432 postgres:16-alpine

db-down:
	@echo "Stopping PostgreSQL database..."
	docker stop task_db
	docker rm task_db

docker-up:
	@echo "Starting all services with Docker Compose..."
	docker-compose up -d

docker-down:
	@echo "Stopping all services..."
	docker-compose down

migrate:
	@echo "Running database migrations..."
	psql -h localhost -U postgres -d task_management -f migrations/001_initial_schema.sql

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Running linter..."
	golangci-lint run ./...

.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
