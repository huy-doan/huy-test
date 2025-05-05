# Makeshop Payment API Project Guide

## I. Introduction

Makeshop Payment API is a backend service built with Golang that manages payment processing for the Makeshop platform. The project follows Domain-Driven Design (DDD) principles.

### Technology Stack

- **Go (v1.24)** - Primary programming language
- **GORM** - ORM for MySQL database interaction
- **MySQL** - Database system
- **JWT** - Authentication mechanism
- **Cobra** - CLI framework for command execution
- **Goose** - Database migration management
- **Swagger** - API documentation
- **Docker** - Containerization
- **Mockery** - Testing utilities

## II. Project Structure
### Description of Key Components:

#### 1. Domain Layer (src/domain)
- **models**: Core business entities like merchants, payment providers, transactions, users, etc.
- **repositories**: Interfaces defining data access methods

#### 2. Usecase Layer (src/usecase)
- Business logic implementation including audit logs, two-factor authentication, and user management

#### 3. Infrastructure Layer (src/infrastructure)
- **auth**: JWT authentication services
- **config**: Application configuration
- **email**: Email service implementations
- **logger**: Logging utilities
- **persistence**: Database connection and repository implementations

#### 4. API Layer (src/api)
- **server.go**: API server setup
- **http**: HTTP server implementation
  - **errors**: Custom error handling
  - **handlers**: Request handlers
  - **middleware**: Authentication and other cross-cutting concerns
  - **response**: Response formatters
  - **router**: API routing
  - **serializers**: Data serialization
  - **validator**: Request validation

#### 5. Database (database)
- **migrations**: Database schema versioning (using Goose)
- **seeds**: Initial data for the application
  - **master**: Master data including roles, permissions, payment providers, etc.

#### 6. Operations (ops)
- **go**: Docker configurations for different environments
  - **developer**: Development environment
  - **sandbox**: Sandbox environment

#### 7. Command Line Interface (cmd)
- Command-line tools for running the service and utilities

## III. Using the Makefile

The project provides multiple make commands to simplify development tasks:

```bash
# Display a list of available commands
make help

# SSH into the backend container
make ssh-be

# SSH into the MySQL container
make ssh-mysql

# Run API in development environment
make run

# Generate Swagger documentation
make swagger

# Format source code
make fmt

# Create a new migration
make migrate-create

# Run migrations up
make migrate-up

# Revert migration (down)
make migrate-down

# Run shell command in container
make shell

# Run seed master
make seed-master

# Generate mocks for testing
make generate-mock

# Run tests with coverage
make test

# Run linting checks
make lint

# Auto-fix linting issues where possible
make lint-fix
```

### Examples

```bash
# Create a new migration
make migrate-create
# Enter migration name: create_new_feature_table

# Run migrations
make migrate-up

# Seed master data
make seed-master
```

## IV. Getting Started

1. Start services with Docker Compose:

2. Run migrations and seed the database:

```bash
make migrate-up
make seed-master
```

3. Access:
   - API: http://localhost:3011
   - Swagger UI: http://localhost:3011/swagger/index.html
