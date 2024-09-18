# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	
	@go build -o main cmd/app/main.go

# Run the application
run:
	@go run cmd/app/main.go

# Create backend containers
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown backend containers
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Build backend containers
docker-build:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up -d --build; \
	fi

# Rebuild the application container
docker-rebuild-app:
	@docker-compose up --build --force-recreate --no-deps -d app

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

# Cleanup all docker data
docker-clean:
	@echo "Listing containers..."
	@containers=$$(docker ps -qa); \
	echo "containers: $$containers"; \
	if [ ! -z "$$containers" ]; then \
	    echo "Stopping containers..."; \
	    docker stop $$containers; \
	    echo "Removing containers..."; \
	    docker rm $$containers; \
	else \
	    echo "No containers found"; \
	fi
	@echo "Listing images..."
	@images=$$(docker images -qa); \
	echo "images: $$images"; \
	if [ ! -z "$$images" ]; then \
	    echo "Removing images..."; \
	    docker rmi -f $$images; \
	else \
	    echo "No images found"; \
	fi
	@echo "Listing volumes..."
	@volumes=$$(docker volume ls -q); \
	echo "volumes: $$volumes"; \
	if [ ! -z "$$volumes" ]; then \
	    echo "Removing volumes..."; \
	    docker volume rm $$volumes; \
	else \
	    echo "No volumes found"; \
	fi
	@echo "Listing networks..."
	@networks=$$(docker network ls -q); \
	echo "networks: $$networks"; \
	if [ ! -z "$$networks" ]; then \
	    echo "Removing networks..."; \
	    docker network rm $$networks; \
	else \
	    echo "No networks found"; \
	fi
	@echo "These should not output any items:"
	@docker ps -a
	@docker images -a
	@docker volume ls
	@echo "This should only show the default networks:"
	@docker network ls

.PHONY: all build run test clean docker-clean
