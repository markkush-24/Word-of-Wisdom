# Use a newer version of Go
FROM golang:1.23 as builder

# Create a working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o client ./cmd/client

# Use a lightweight image to run the client
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/client .
CMD ["./client"]
