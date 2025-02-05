#!/bin/bash
# Interactive script to manage individual Docker containers (without docker-compose)
# Make sure this file has execution permissions: chmod +x run_individual_docker.sh

while true; do
  echo "======================================"
  echo "Select an action for individual Docker containers:"
  echo "1) Create Docker network (my_network)"
  echo "2) Build and run server container"
  echo "3) View server container logs"
  echo "4) Build and run client container"
  echo "5) View client container logs"
  echo "6) Stop and remove server container"
  echo "7) Stop and remove client container"
  echo "0) Exit"
  echo "======================================"

  read -p "Enter your choice: " choice

  case $choice in
    1)
      echo "Creating Docker network 'my_network'..."
      docker network create my_network || echo "Network 'my_network' already exists"
      ;;
    2)
      echo "Building server Docker image..."
      docker build -t server-image -f build/server/Dockerfile .
      echo "Running server container..."
      docker run -d --name server-container --network my_network -p 8080:8080 server-image
      ;;
    3)
      echo "Displaying server container logs..."
      docker logs server-container
      ;;
    4)
      echo "Building client Docker image..."
      docker build -t client-image -f build/client/Dockerfile .
      echo "Running client container..."
      docker run -d --name client-container --network my_network -p 8081:8081 client-image
      ;;
    5)
      echo "Displaying client container logs..."
      docker logs client-container
      ;;
    6)
      echo "Stopping and removing server container..."
      docker stop server-container && docker rm server-container
      ;;
    7)
      echo "Stopping and removing client container..."
      docker stop client-container && docker rm client-container
      ;;
    0)
      echo "Exiting..."
      exit 0
      ;;
    *)
      echo "Invalid choice, please try again."
      ;;
  esac

  echo ""
  echo "Press Enter to continue..."
  read -r
  clear
done
