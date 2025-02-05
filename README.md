# Word of Wisdom Project

## Overview

The Word of Wisdom project implements a TCP server that provides a quote after the client successfully solves a Proof-of-Work (PoW) challenge. The PoW mechanism, which requires calculating a SHA-256 hash with a specified number of leading zeros, is used to protect the server from potential DDoS attacks.

**Project Requirements:**
- **Design and Implementation:**  
  Design and implement a "Word of Wisdom" TCP server.
- **DDoS Protection:**  
  Protect the TCP server from DDoS attacks using a Proof-of-Work (PoW) challenge-response protocol.
- **PoW Algorithm Choice:**  
  Explain the choice of the PoW algorithm.
- **Quote Delivery:**  
  After PoW verification, the server sends one of the quotes from the "Word of Wisdom" book or another collection.
- **Docker Support:**  
  Provide Dockerfiles for both the server and the client that solves the PoW challenge.

### Why SHA-256?
We chose the SHA-256 algorithm because:
- **Simplicity:** It is simple to implement.
- **Determinism:** The hash function is deterministic (same input always produces the same output), making verification straightforward.
- **Performance:** It offers fast verification on the server side.
- **Security:** SHA-256 is a well-known cryptographic hash function with strong security properties.
- **Flexibility with Challenge:** To ensure uniqueness, a dynamic challenge is generated for each client request.

## Prerequisites
- **Docker:** Ensure Docker is installed and running.
- **docker-compose:** Ensure docker-compose is installed.
- **Bash:** The provided scripts are written in Bash and require a Unix-like environment (Linux, macOS, or Windows with WSL).
- **Go:** Version 1.19+ is required for building and testing the application.

## Methods to Launch the Application
This project offers **two distinct methods** for launching the application:

## Method 1: Running via Individual Dockerfile Operations
This approach gives you granular control over each container using separate Dockerfiles.

### Interactive Script: `run_individual_docker.sh`
This script allows you to:
- Create the required Docker network.
- Build and run the server container.
- Build and run the client container.
- View logs for the server and client.
- Stop and remove the containers individually.

#### Steps to Use:
1. **Make the Script Executable**  
   In your project root, run:
   ```bash
   chmod +x run_individual_docker.sh
   ```
2. **Run the Script**  
   Execute the script:
   ```bash
   ./run_individual_docker.sh
   ```
### Follow the Interactive Menu
The script displays a menu with the following options:
- **Option 1**: Create the Docker network (`my_network`).
   - This network enables communication between the containers.
- **Option 2**: Build and run the server container (listening on port 8080).
- **Option 3**: View the server container logs.
- **Option 4**: Build and run the client container (listening on port 8081).
- **Option 5**: View the client container logs.
- **Option 6**: Stop and remove the server container.
- **Option 7**: Stop and remove the client container.
- **Option 0**: Exit the script.

### Sequence for Individual Dockerfile Approach
1. Create the Docker network (Option 1).
2. Build and run the server container (Option 2).
3. Build and run the client container (Option 4).
4. Monitor logs as needed (Options 3 and 5).
5. Stop and remove containers when done (Options 6 and 7).

## Method 2: Running via Docker Compose
This approach uses docker-compose to manage both containers together in detached mode. It’s ideal if you want to launch all services as a single unit and later monitor or restart them individually.

### Interactive Script: `run_docker_compose.sh`
This script allows you to:
- Start (build) all services via docker-compose in detached mode.
- View logs for the server and client containers.
- Restart the client container to request a new quote.
- View logs for both containers together.

#### Steps to Use:
1. **Make the Script Executable**  
   In your project root, run:
   ```bash
   chmod +x run_docker_compose.sh
   ```
2. **Run the Script**  
   Execute the script:
   ```bash
   ./run_docker_compose.sh
   ```
   
### Follow the Interactive Menu:
- **Option 1**: Start (build) services via docker-compose in detached mode.
   - The containers will run in the background, so you won’t be attached to their logs.
- **Option 2**: View the server container logs.
- **Option 3**: View the client container logs.
- **Option 4**: Restart the client container to get a new quote.
   - This allows you to request a new quote without stopping the entire service.
- **Option 5**: View logs for both the server and client containers.
- **Option 0**: Exit the script.

### Summary of the Sequence (Docker Compose Approach):
- Start services in detached mode (Option 1).
- Monitor logs separately (Options 2 and 3).
- Alternatively, view both logs together (Option 5).
- Restart the client container for a new quote (Option 4).

## Makefile Commands
In addition to the interactive scripts, a Makefile is provided with several commands for convenience. Some of the key commands include:
- `make lint` – Run code linting using golangci-lint.
- `make test` – Run unit tests with verbose output.
- `make docker-network` – Create the Docker network my_network.
- `make docker-run-server` – Build and run the server container.
- `make docker-logs-server` – View server container logs.
- `make docker-run-client` – Build and run the client container.
- `make docker-logs-client` – View client container logs.
- `make compose-up` – Start services using docker-compose in detached mode.
- `make docker-start-client` – Restart the client container for a new quote.
- `make docker-stop-server` & `make docker-stop-client` – Stop and remove the containers.

You can use these commands directly if you prefer not to use the interactive scripts.

# Additional Development Tasks

Before running the project, it is recommended to perform the following tasks:

## Code Linting

Run the following command to check the code quality using golangci-lint:
```bash
make lint
```

## Unit Tests

Run the following command to run unit tests:
```bash
make test
```

# Implementation Details

## PoW Challenge-Response Protocol

To protect the server from DDoS attacks, a Proof-of-Work (PoW) mechanism is implemented using a challenge-response protocol.

### Challenge Generation

Each time a client connects, the server generates a unique challenge (e.g., based on the current time in nanoseconds) and sends it to the client.

### Client Processing

The client receives the challenge and computes a PoW solution by finding a nonce such that the SHA-256 hash of the concatenated string (challenge + fixed prefix + nonce) starts with the required number of zeros.

### Server Verification

The server uses the SHA-256 hash function to quickly verify the client’s solution.

## Choice of SHA-256

The SHA-256 algorithm is used because:
- It is well-known and widely adopted.
- It is simple to implement.
- It allows fast verification on the server side.
- Its deterministic nature ensures that the same input always produces the same output, simplifying verification.

By combining SHA-256 with a unique challenge, we ensure that each PoW solution is unique and cannot be precomputed.

# Conclusion

The Word of Wisdom project offers two paths to launch the application:

## Individual Dockerfile Approach

Provides granular control through the interactive script `run_individual_docker.sh`. This approach requires you to manually create the network, build images, run containers, and view logs individually.

## Docker Compose Approach

Provides a unified solution via the interactive script `run_docker_compose.sh`, which launches both services in detached mode. You can then monitor logs and restart the client container as needed.