help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Run all tests
	go test ./tests/acceptance -v

.PHONY: build
build: ## Build the service binary (no Docker)
	go build -o bin/delivery-server internal/server.go

.PHONY: run
run: build ## Build and run the service locally (no Docker)
	@test -f .env || cp .env.example .env
	ENV_FILE="$(CURDIR)/.env" ./bin/delivery-server

# Docker variables
DOCKER_REGISTRY ?= registry.example.com
IMAGE_NAME ?= delivery-service
IMAGE_TAG ?= latest
FULL_IMAGE_NAME = $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(FULL_IMAGE_NAME) .