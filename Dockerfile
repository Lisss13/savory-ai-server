# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Установить git/ssh/ca-certificates для go mod download (и обновить сертификаты)
RUN apk add --no-cache git openssh ca-certificates

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/main.go

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/server /app/server

# Copy configuration files
COPY config/config.docker.toml /app/config/config.docker.toml
COPY config/config.railway.toml /app/config/config.railway.toml

# Copy the storage directory for static files
COPY storage /app/storage

# Expose the port (Railway uses PORT env var)
EXPOSE 4000

# Default config (can be overridden by CONFIG_FILE env var)
ENV CONFIG_FILE="/app/config/config.docker.toml"

# Run the application with migrations
CMD ["/app/server", "-migrate"]
