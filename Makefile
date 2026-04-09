.PHONY: run build test migrate docker-up docker-down tidy

run:
	go run ./cmd/api/main.go

build:
	go build -o bin/api ./cmd/api/main.go

test:
	go test ./...

migrate:
	@echo "Run: migrate -path migrations -database $${DATABASE_URL} up"

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

tidy:
	go mod tidy
