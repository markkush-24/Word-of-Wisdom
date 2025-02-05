#!/bin/bash
# Interactive script to manage services via docker-compose in detached mode
# and to view logs from both server and client containers or restart the client container.
# Make sure this file has execution permissions: chmod +x run_docker_compose.sh

while true; do
  echo "======================================"
  echo "Select an action for docker-compose:"
  echo "1) Start (build) services via docker-compose in detached mode"
  echo "2) View server container logs"
  echo "3) View client container logs"
  echo "4) Restart client container to get a new quote"
  echo "5) View logs for both server and client containers"
  echo "0) Exit"
  echo "======================================"

  read -p "Enter your choice: " choice

  case $choice in
    1)
      echo "Starting services via docker-compose in detached mode..."
      docker-compose up --build -d
      ;;
    2)
      echo "Displaying server container logs..."
      docker logs server-container
      ;;
    3)
      echo "Displaying client container logs..."
      docker logs client-container
      ;;
    4)
      echo "Restarting client container..."
      docker restart client-container
      ;;
    5)
      echo "Displaying logs for both server and client containers..."
      echo "----- Server logs -----"
      docker logs server-container
      echo "----- Client logs -----"
      docker logs client-container
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
