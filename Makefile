help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: codegen-types
codegen-types: ## Compile types
	oapi-codegen -generate="types" -package delivery openapi/delivery.yml > internal/generated/Types.gen.go
test:
	go test ./tests/integration -v

# Docker variables
DOCKER_REGISTRY ?= registry.example.com
IMAGE_NAME ?= delivery-service
IMAGE_TAG ?= latest
FULL_IMAGE_NAME = $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(FULL_IMAGE_NAME) .