# Start from the official Golang image for building
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o client .

# Use a minimal image for running
FROM alpine:latest

WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/client .

# Expose port if needed (e.g., 8080)
EXPOSE 8080

# Run the binary
CMD ["./client"]