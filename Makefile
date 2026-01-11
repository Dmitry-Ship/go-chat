.PHONY: help secret backend_build backend_test backend_lint backend_run
.PHONY: frontend_build frontend_test frontend_lint frontend_type frontend_dev frontend_prod
.PHONY: all_build all_test all_lint docker_up docker_up_detached docker_down docker_build clean

help:
	@echo "Available commands:"
	@echo "  secret           - Generate a cryptographically secure secret"
	@echo "  backend_build    - Build Go backend server"
	@echo "  backend_test     - Run Go backend tests"
	@echo "  backend_lint     - Run Go backend linter"
	@echo "  backend_run      - Run Go backend server"
	@echo ""
	@echo "  frontend_build   - Build Next.js frontend"
	@echo "  frontend_test    - Run frontend tests"
	@echo "  frontend_lint    - Run frontend linter"
	@echo "  frontend_type    - Run frontend type check"
	@echo "  frontend_dev     - Run frontend dev server"
	@echo "  frontend_prod    - Run frontend production server"
	@echo ""
	@echo "  all_build        - Build both backend and frontend"
	@echo "  all_test         - Run all tests (backend + frontend)"
	@echo "  all_lint         - Run all linters (backend + frontend)"
	@echo "  docker_up        - Start all services with Docker"
	@echo "  docker_up_detached - Start all services in detached mode"
	@echo "  docker_down      - Stop all Docker services"
	@echo "  docker_build     - Build all Docker images"
	@echo "  clean            - Remove build artifacts and dependencies"

backend_build:
	@echo "Building Go backend..."
	cd backend && go build -o go-bin ./cmd/server

backend_test:
	@echo "Running Go backend tests..."
	cd backend && go test -timeout 2m ./...

backend_lint:
	@echo "Running Go backend linter..."
	cd backend && golangci-lint run

backend_run:
	@echo "Running Go backend server..."
	cd backend && go run ./cmd/server

frontend_build:
	@echo "Building Next.js frontend..."
	cd frontend && npm run build

frontend_test:
	@echo "Running frontend tests..."
	cd frontend && npm test

frontend_lint:
	@echo "Running frontend linter..."
	cd frontend && npm run lint

frontend_type:
	@echo "Running frontend type check..."
	cd frontend && npx tsc --noEmit

frontend_dev:
	@echo "Running frontend dev server..."
	cd frontend && npm run dev

frontend_prod:
	@echo "Running frontend production server..."
	cd frontend && npm run start

all_build: backend_build frontend_build

all_test: backend_test frontend_test

all_lint: backend_lint frontend_lint frontend_type

docker_up:
	@echo "Starting all services with Docker..."
	docker-compose up --build

docker_up_detached:
	@echo "Starting all services in detached mode..."
	docker-compose up -d

docker_down:
	@echo "Stopping all Docker services..."
	docker-compose down

docker_build:
	@echo "Building all Docker images..."
	docker-compose build

clean:
	@echo "Cleaning build artifacts..."
	rm -rf backend/go-bin
	cd frontend && rm -rf node_modules .next

secret:
	@openssl rand -base64 32
