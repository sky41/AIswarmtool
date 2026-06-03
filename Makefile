.PHONY: help dev build test up down clean

help:
	@echo "Usage:"
	@echo "  make dev      - Start development environment"
	@echo "  make build    - Build the application"
	@echo "  make test     - Run tests"
	@echo "  make up       - Start all services with Docker Compose"
	@echo "  make down     - Stop all services"
	@echo "  make clean    - Remove containers and volumes"

dev:
	@echo "Starting development environment..."
	cd backend && go run cmd/main.go

build:
	@echo "Building backend..."
	cd backend && go build -o bin/main cmd/main.go

test:
	@echo "Running tests..."
	cd backend && go test ./...

up:
	@echo "Starting all services with Docker Compose..."
	docker-compose up -d

down:
	@echo "Stopping all services..."
	docker-compose down

clean:
	@echo "Cleaning up..."
	docker-compose down -v

