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

- **app**: Application code
  - **middleware**: Middleware components
  - **module**: Business logic modules
  - **router**: API routes
  - **storage**: Database models
- **cmd**: Application entry point
- **config**: Configuration files
- **http**: HTTP request examples
- **internal**: Internal packages
- **storage**: Static files and uploads
- **utils**: Utility functions

## License

[License information]
