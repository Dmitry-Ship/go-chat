.DEFAULT_GOAL := help

BACKEND_DIR := backend
FRONTEND_DIR := frontend
DOCKER_COMPOSE := docker-compose

.PHONY: \
	help secret clean \
	backend_build backend_test backend_lint backend_run \
	frontend_build frontend_test frontend_lint frontend_type frontend_dev frontend_prod \
	docker_deps_up docker_deps_up_detached docker_deps_down

help:
	@echo "Available commands:"
	@echo "  secret                  - Generate a cryptographically secure secret"
	@echo "  clean                   - Remove build artifacts and frontend dependencies"
	@echo ""
	@echo "  backend_build           - Build Go backend server locally (no Docker)"
	@echo "  backend_test            - Run Go backend tests"
	@echo "  backend_lint            - Run Go backend linter"
	@echo "  backend_run             - Run Go backend server"
	@echo ""
	@echo "  frontend_build          - Build Next.js frontend locally (no Docker)"
	@echo "  frontend_test           - Run frontend tests"
	@echo "  frontend_lint           - Run frontend linter"
	@echo "  frontend_type           - Run frontend type check"
	@echo "  frontend_dev            - Run frontend dev server"
	@echo "  frontend_prod           - Run frontend production server"
	@echo ""
	@echo "  docker_deps_up          - Start only Postgres and Redis containers"
	@echo "  docker_deps_up_detached - Start Postgres and Redis in detached mode"
	@echo "  docker_deps_down        - Stop Postgres and Redis containers"

backend_build:
	@echo "Building Go backend..."
	cd $(BACKEND_DIR) && go build -o go-bin ./cmd/server

backend_test:
	@echo "Running Go backend tests..."
	cd $(BACKEND_DIR) && go test -timeout 2m ./...

backend_lint:
	@echo "Running Go backend linter..."
	cd $(BACKEND_DIR) && golangci-lint run

backend_run:
	@echo "Running Go backend server..."
	cd $(BACKEND_DIR) && go run ./cmd/server

frontend_build:
	@echo "Building Next.js frontend..."
	cd $(FRONTEND_DIR) && npm run build

frontend_test:
	@echo "Running frontend tests..."
	cd $(FRONTEND_DIR) && npm run test

frontend_lint:
	@echo "Running frontend linter..."
	cd $(FRONTEND_DIR) && npm run lint

frontend_type:
	@echo "Running frontend type check..."
	cd $(FRONTEND_DIR) && npx tsc --noEmit

frontend_dev:
	@echo "Running frontend dev server..."
	cd $(FRONTEND_DIR) && npm run dev

frontend_prod:
	@echo "Running frontend production server..."
	cd $(FRONTEND_DIR) && npm run start

docker_deps_up:
	@echo "Starting Postgres and Redis containers..."
	$(DOCKER_COMPOSE) up postgres redis

docker_deps_up_detached:
	@echo "Starting Postgres and Redis containers in detached mode..."
	$(DOCKER_COMPOSE) up -d postgres redis

docker_deps_down:
	@echo "Stopping Postgres and Redis containers..."
	$(DOCKER_COMPOSE) stop postgres redis

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BACKEND_DIR)/go-bin
	cd $(FRONTEND_DIR) && rm -rf node_modules .next

secret:
	@openssl rand -base64 32
