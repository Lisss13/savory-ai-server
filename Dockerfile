# Build stage
FROM golang:1.23-alpine AS builder

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

# Copy the configuration file
COPY config/config.docker.toml /app/config/config.toml

#COPY storage/selfsigned.crt /app/storage/selfsigned.crt
#COPY storage/selfsigned.key /app/storage/selfsigned.key

# Copy the storage directory for static files
COPY storage /app/storage

# Expose the port
EXPOSE 4000

# Run the application
# The CONFIG_FILE environment variable can be used to specify a custom config file
ENV CONFIG_FILE="/app/config/config.toml"
CMD ["/app/server", "-migrate"]
