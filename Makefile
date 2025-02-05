.PHONY: lint test docker-network docker-build-server docker-run-server docker-logs-server \
docker-build-client docker-run-client docker-logs-client compose-up docker-start-client \
docker-stop-server docker-stop-client help

# Code linting
lint:
	@golangci-lint run ./...

# Run tests with detailed output
test:
	@go test -v ./...

# Create a Docker network
docker-network:
	@echo "Creating Docker network 'my_network'..."
	@docker network create my_network || echo "Network 'my_network' already exists"

# Build the server Docker image
docker-build-server:
	@echo "Building server image..."
	@docker build -t server-image -f build/server/Dockerfile .

# Run the server container
docker-run-server: docker-network docker-build-server
	@echo "Running server container..."
	@docker run -d --name server-container --network my_network -p 8080:8080 server-image
	@echo "Container is running. Check the list of containers with: docker ps"
	@echo "Container logs: docker logs server-container"
	@echo "Use ServerAddressDocker = \"server-container:8080\" for connection"

# View server container logs
docker-logs-server:
	@docker logs server-container

# Build the client Docker image
docker-build-client:
	@echo "Building client image..."
	@docker build -t client-image -f build/client/Dockerfile .

# Run the client container
docker-run-client: docker-network docker-build-client
	@echo "Running client container..."
	@docker run -d --name client-container --network my_network -p 8081:8081 client-image
	@echo "Container is running. Check the list of containers with: docker ps"
	@echo "Container logs: docker logs client-container"

# View client container logs
docker-logs-client:
	@docker logs client-container

# Stop and remove the server container
docker-stop-server:
	@echo "Stopping and removing server container..."
	@docker stop server-container && docker rm server-container

# Stop and remove the client container
docker-stop-client:
	@echo "Stopping and removing client container..."
	@docker stop client-container && docker rm client-container

# Start both services using docker-compose (if docker-compose.yml is available)
compose-up:
	@echo "Starting server and client containers using docker-compose..."
	@docker-compose up --build

# Restart the client container to get a new quote
docker-start-client:
	@echo "Restarting client container..."
	@docker start client-container
	@echo "Check container logs with: docker logs client-container"

# Display help for Make commands
help:
	@echo "Available commands:"
	@echo "  make lint                - Run golangci-lint"
	@echo "  make test                - Run tests with detailed output"
	@echo "  make docker-network      - Create Docker network 'my_network'"
	@echo "  make docker-build-server - Build the server Docker image"
	@echo "  make docker-run-server   - Run the server container"
	@echo "  make docker-logs-server  - View server container logs"
	@echo "  make docker-build-client - Build the client Docker image"
	@echo "  make docker-run-client   - Run the client container"
	@echo "  make docker-logs-client  - View client container logs"
	@echo "  make docker-stop-server  - Stop and remove the server container"
	@echo "  make docker-stop-client  - Stop and remove the client container"
	@echo "  make compose-up          - Start both services using docker-compose"
	@echo "  make docker-start-client - Restart the client container to get a new quote"
	@echo "  make help                - Display this help"