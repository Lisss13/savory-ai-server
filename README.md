# Savory AI Server

This is the server component of the Savory AI application. It provides APIs for restaurant management, table management, chat functionality, and more.

## Quick Start with Docker

The easiest way to run the Savory AI server is using Docker and Docker Compose.

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Running the Server

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd savory_ai/server
   ```

2. Start the services:
   ```bash
   docker-compose up -d
   ```

3. The server will be available at http://localhost:4000

4. To stop the services:
   ```bash
   docker-compose down
   ```

## Running Without Docker

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.18 or higher)
- [PostgreSQL](https://www.postgresql.org/download/) (version 13 or higher)

### Database Setup

1. Create a PostgreSQL database:
   ```sql
   CREATE DATABASE attic_db;
   CREATE USER attic_admin WITH ENCRYPTED PASSWORD 'my_pass';
   GRANT ALL PRIVILEGES ON DATABASE attic_db TO attic_admin;
   ```

2. Update the database connection string in `config/config.toml` if needed:
   ```toml
   [db.postgres]
   dsn = "postgresql://attic_admin:my_pass@127.0.0.1:5432/attic_db"
   ```

### Running the Server

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd savory_ai/server
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build and run the server:
   ```bash
   go run cmd/main.go -migrate
   ```

4. The server will be available at http://localhost:4000

## Configuration

The configuration file is located at `config/config.toml`. Here are the main configuration sections:

- **app**: Application settings like name, port, and timeouts
- **db.postgres**: Database connection settings
- **logger**: Logging configuration
- **middleware**: Settings for various middleware components

When running with Docker, a special configuration file `config/config.docker.toml` is used by default, which is configured to work with the Docker environment.

### Custom Configuration with Docker

You can provide a custom configuration file when running with Docker by:

1. Creating your custom configuration file (e.g., `config/custom-config.toml`)
2. Mounting the configuration directory as a volume (already done in docker-compose.yaml)
3. Setting the `CONFIG_FILE` environment variable to point to your custom configuration file

Example in docker-compose.yaml:
```yaml
services:
  app:
    environment:
      - CONFIG_FILE=/app/config/custom-config.toml
```
delete all dockers 
```bash
docker stop $(docker ps -aq) && docker rm $(docker ps -aq) && docker rmi -f $(docker images -q)
```

Or when running with docker run:
```bash
docker run -v ./config:/app/config -e CONFIG_FILE=/app/config/custom-config.toml -p 4000:4000 savory-ai-server
```

## API Documentation

The API documentation is available in the `http` directory, which contains HTTP request examples for testing the API endpoints.

## Project Structure

```
server/
├── cmd/main.go              # Entry point (Uber FX for DI)
├── app/
│   ├── middleware/          # HTTP middleware (CORS, JWT, Rate Limit)
│   ├── module/              # 11 business modules
│   ├── router/              # API route registration
│   └── storage/             # GORM database models
├── internal/bootstrap/      # Application initialization
│   └── database/            # PostgreSQL connection
├── utils/                   # Utilities (config, jwt, response, helpers)
├── config/                  # TOML configurations
├── http/                    # API request examples
└── storage/public/          # Static files (images, QR codes)
```

### Business Modules (11 total)

Each module follows the MVC pattern:

| Module | Description |
|--------|-------------|
| **auth** | Authentication & registration |
| **user** | User management |
| **organization** | Organization management |
| **restaurant** | Restaurant information |
| **menu_category** | Menu categories |
| **dish** | Menu dishes |
| **table** | Restaurant tables |
| **question** | Questions (with multilingual support) |
| **qr_code** | QR code generation |
| **file_upload** | File upload service |
| **chat** | Chat functionality |

### Module Structure

Each module follows this structure:

```
module/
├── {module}_module.go       # Module initialization & DI
├── controller/
│   └── {module}_controller.go   # HTTP request handlers
├── service/
│   └── {module}_service.go      # Business logic
├── repository/
│   └── {module}_repository.go   # Data access layer
└── payload/
    ├── request.go               # Request DTOs
    └── response.go              # Response DTOs
```

### Technology Stack

| Technology | Purpose |
|------------|---------|
| **Fiber v2** | Web framework |
| **GORM** | ORM for PostgreSQL |
| **Uber FX** | Dependency injection |
| **JWT** | Authentication |
| **Zerolog** | Structured logging |
| **Docker** | Containerization |
| **go-validator** | Request validation |
| **go-qrcode** | QR code generation |

### Architecture Patterns

- **Clean Architecture** - Separation of concerns
- **Repository Pattern** - Data access abstraction
- **Service Layer** - Business logic isolation
- **DTO Pattern** - Request/Response objects
- **Dependency Injection** - Using Uber FX
- **Middleware Pattern** - Fiber middleware composition

