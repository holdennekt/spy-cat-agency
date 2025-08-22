# Spy Cat Agency Management System

CRUD REST API for managing spy cats, missions, and targets built with Go, Gin, GORM, and PostgreSQL.

## Features

- **Spy Cats Management**: Create, read, update, and delete spy cats with breed validation
- **Mission Management**: Create missions with targets, assign cats, track completion
- **Target Management**: Add, update, and remove targets from missions
- **Business Rules Enforcement**: Comprehensive validation of business logic
- **External API Integration**: Breed validation using TheCatAPI
- **Database Migrations**: Automated PostgreSQL schema management
- **Docker Support**: Complete containerization for development and deployment
- **API Documentation**: Comprehensive OpenAPI/Swagger documentation
- **Unit Tests**: Test coverage for business logic

## Business Rules

- Cats can only have one active mission at a time
- Missions must have 1-3 targets
- Cannot update notes if target/mission is completed
- Cannot delete completed targets
- Cannot add targets to completed missions
- Cannot delete assigned missions
- Cat breeds are validated against TheCatAPI

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)
- PostgreSQL (for local development without Docker)

### Using Docker (Recommended)

1. Clone the repository:
   ```bash
   git clone https://github.com/holdennekt/spy-cat-agency
   cd spy-cat-agency
   ```

2. Copy environment file:
   ```bash
   cp .env.example .env
   ```

3. Start the application:
   ```bash
   docker-compose up --build
   ```

The API will be available at `http://localhost:8080`

## API Documentation

The API documentation is available in OpenAPI/Swagger format at `docs/swagger.yaml`.

Or at `http://localhost:8080/swagger/index.html`

## Testing

Run unit tests:
```bash
go test ./internal/service/tests/...
```

### Database Migrations

Migrations are automatically applied on application startup using GORM's AutoMigrate feature. Manual SQL migrations are also available in the `migrations/` directory.
