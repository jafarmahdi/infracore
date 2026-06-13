## InfraCore Makefile
## Usage: make <target>

APP_NAME    := infracore
MODULE      := github.com/infracore/infracore
CMD_API     := ./cmd/api
CMD_WORKER  := ./cmd/worker
CMD_AGENT   := ./cmd/agent
BUILD_DIR   := ./build
MIGRATE_DIR := ./migrations/postgres
DB_URL      ?= postgres://infracore:infracore_secret@localhost:5432/infracore?sslmode=disable

.PHONY: all build build-api build-worker build-agent \
        run-api run-worker \
        migrate-up migrate-down migrate-create \
        test test-unit test-integration \
        lint vet fmt \
        docker-build docker-up docker-down \
        swag clean

## ── Build ────────────────────────────────────────────────────
all: build

build: build-api build-worker build-agent

build-api:
	@echo "Building API server..."
	go build -ldflags "-s -w" -o $(BUILD_DIR)/$(APP_NAME)-api $(CMD_API)

build-worker:
	@echo "Building background worker..."
	go build -ldflags "-s -w" -o $(BUILD_DIR)/$(APP_NAME)-worker $(CMD_WORKER)

build-agent:
	@echo "Building monitoring agent..."
	go build -ldflags "-s -w" -o $(BUILD_DIR)/$(APP_NAME)-agent $(CMD_AGENT)

## ── Run ──────────────────────────────────────────────────────
run-api:
	go run $(CMD_API)/main.go

run-worker:
	go run $(CMD_WORKER)/main.go

## ── Database Migrations ──────────────────────────────────────
migrate-up:
	migrate -path $(MIGRATE_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATE_DIR) -database "$(DB_URL)" down 1

migrate-down-all:
	migrate -path $(MIGRATE_DIR) -database "$(DB_URL)" down

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir $(MIGRATE_DIR) -seq $$name

migrate-force:
	@read -p "Version: " ver; \
	migrate -path $(MIGRATE_DIR) -database "$(DB_URL)" force $$ver

## ── Testing ──────────────────────────────────────────────────
test:
	go test ./... -v -race -timeout 120s

test-unit:
	go test ./internal/... -v -race -short

test-integration:
	go test ./tests/integration/... -v -race -timeout 300s

test-coverage:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html

## ── Code Quality ─────────────────────────────────────────────
lint:
	golangci-lint run ./...

vet:
	go vet ./...

fmt:
	gofmt -w -s .
	goimports -w .

## ── Swagger ──────────────────────────────────────────────────
swag:
	swag init -g $(CMD_API)/main.go -o ./docs/swagger

## ── Docker ───────────────────────────────────────────────────
docker-build:
	docker build -f deployments/docker/Dockerfile.api -t $(APP_NAME)-api:latest .
	docker build -f deployments/docker/Dockerfile.worker -t $(APP_NAME)-worker:latest .
	docker build -f deployments/docker/Dockerfile.agent -t $(APP_NAME)-agent:latest .

docker-up:
	docker compose -f deployments/compose/docker-compose.yml up -d

docker-down:
	docker compose -f deployments/compose/docker-compose.yml down

docker-logs:
	docker compose -f deployments/compose/docker-compose.yml logs -f

## ── Deps ─────────────────────────────────────────────────────
deps:
	go mod download
	go mod tidy

## ── Cleanup ──────────────────────────────────────────────────
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
