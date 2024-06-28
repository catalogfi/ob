.PHONY: build run stop logs status restart pull

DOCKER_FILE=Dockerfile
DOCKER=docker
CONTAINER_NAME=ghcr.io/catalogfi/orderbook

all: build

define print_message
	@echo "\033[1;34m$(1)\033[0m"
endef

# Build and start orderbook containers
build:
	$(call print_message, "Building orderbook")
	@$(DOCKER) build ./ -f $(DOCKER_FILE) -t $(CONTAINER_NAME)

# Run the container
run:
	$(call print_message, "Running orderbook container")
	@$(DOCKER) run -d --name $(CONTAINER_NAME) $(CONTAINER_NAME)

pull:
	$(call print_message, "Pulling orderbook container")
	@$(DOCKER) pull $(CONTAINER_NAME)

# Stop the container
stop:
	$(call print_message, "Stopping orderbook container")
	@$(DOCKER) stop $(CONTAINER_NAME) 
	@$(DOCKER) rm $(CONTAINER_NAME) 

# View logs of the container
logs:
	$(call print_message, "Viewing orderbook container logs")
	@$(DOCKER) logs -f $(CONTAINER_NAME)

# Check the status of the container
status:
	$(call print_message, "Checking orderbook container status")
	@$(DOCKER) ps -a | grep $(CONTAINER_NAME) || echo "$(CONTAINER_NAME) is not running"

# Restart the container
restart: stop run
	$(call print_message, "Restarting orderbook container")
