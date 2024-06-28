# Movie Project

A Go-based microservice for managing movie information using gRPC and RESTful APIs.

## Features

- CRUD operations for movies
- gRPC and REST API support
- Swagger UI for API documentation
- Prometheus metrics
- Structured logging (log/slog)
- PostgreSQL database

## Prerequisites

- Go 1.16+
- PostgreSQL
- Docker (optional)

## Quick Start

1. Clone the repository:
   ```
   git clone https://github.com/dkdkv/movie-project.git
   cd movie-project
   ```

2. Set up the environment variables (copy `.env.example` to `.env` and edit as needed).

3. Run database migrations:
   ```
   go run cmd/migrate/main.go
   ```

4. Start the server:
   ```
   go run cmd/server/main.go
   ```

5. Access the Swagger UI at `http://localhost:8080/docs`

## API Endpoints

- gRPC: `localhost:50051`
- REST: `http://localhost:8080`
- Metrics: `http://localhost:8080/metrics`

## Development

Generate Protocol Buffers:
```
./generate_proto.bat  # On Windows
./generate_proto.sh   # On Unix-based systems
```
