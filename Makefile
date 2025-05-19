.PHONY: docs build run clean docker-build docker-run docker-compose-up docker-compose-down docker-compose-stage-up docker-compose-stage-down

# Go related variables
BINARY_NAME=server
MAIN_PATH=./src/cmd/server

# Docker related variables
DOCKER_IMAGE=new-site-builder
DOCKER_TAG=latest

# Generate API documentation
docs:
	swag init -g */**/*.go -o ./docs

# Build the application
build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	go run $(MAIN_PATH)

# Clean build files
clean:
	rm -f $(BINARY_NAME)
	rm -rf ./docs

# Docker build
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f ./docker/Dockerfile .

# Docker run
docker-run:
	docker run -p 9595:8585 $(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker compose up for development
docker-compose-up:
	docker-compose -f ./docker/docker-compose.yml up -d

# Docker compose down for development
docker-compose-down:
	docker-compose -f ./docker/docker-compose.yml down

# Docker compose up for staging environment
docker-compose-stage-up:
	docker-compose -f ./docker/docker-compose.stage.yml up -d

# Docker compose down for staging environment
docker-compose-stage-down:
	docker-compose -f ./docker/docker-compose.stage.yml down

