.PHONY: docs build run clean docker-build docker-run docker-compose-up docker-compose-down docker-compose-stage-up docker-compose-stage-down deploy

# Go related variables
BINARY_NAME=server
MAIN_PATH=./cmd/server

# Docker related variables
DOCKER_IMAGE=new-site-builder
DOCKER_TAG=latest
vet:
	go vet ./...
# Generate API documentation
init:
	@command -v swag >/dev/null 2>&1 || go install github.com/swaggo/swag/cmd/swag@latest
	@command -v air >/dev/null 2>&1 || go install github.com/air-verse/air@latest

docs:
	swag init -g cmd/server/main.go -o ./docs

# Build the application
build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Clean build files
clean:
	rm -f $(BINARY_NAME)
	rm -rf ./docs

# Docker build
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f ./docker/Dockerfile .

# Docker compose up for development
up:
	docker-compose -f ./docker/docker-compose.yml up redis mongodb mysql minio rabbitmq -d --force-recreate --build --remove-orphans
stage-up:
	docker-compose -f ./docker/docker-compose.yml up -d --force-recreate --build --remove-orphans
# Docker compose down for development
down:
	docker-compose -f ./docker/docker-compose.yml redis mongodb mysql minio rabbitmq down
stage-down:
	docker-compose -f ./docker/docker-compose.yml down

move-all-vendor:
	@echo "📦 Moving all vendor packages to GOPATH/src..."
	@sudo find vendor -type d -name '.git' -prune -o -type d -print | tail -n +2 | while read dir; do \
		relpath=$${dir#vendor/}; \
		dest="$(GOPATH)/$$relpath"; \
		echo "➡️  Moving $$dir to $$dest"; \
		mkdir -p "$$dest"; \
		cp -r $$dir/* "$$dest/" 2>/dev/null || true; \
	done
	@echo "✅ All vendor packages moved."

stage-deploy:
	git pull
	$(MAKE) move-all-vendor
	$(MAKE) clean
	$(MAKE) init
	$(MAKE) docs
	$(MAKE) stage-up

deploy:
	git pull
	$(MAKE) move-all-vendor
	$(MAKE) clean
	$(MAKE) init
	$(MAKE) docs
	$(MAKE) up
	air