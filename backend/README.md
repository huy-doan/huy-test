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

The project provides multiple task commands to simplify development tasks:

```bash
# Display a list of available commands
task

# SSH into the backend container
task ssh-go

# SSH into the MySQL container
task ssh-mysql

# Run API in development environment
task run

# Generate Swagger documentation
task swagger

# Format source code
task fmt

# Create a new migration
task migrate-create

# Run migrations up
task migrate-up

# Revert migration (down)
task migrate-down

# Run shell command in container
task shell -- <name_of_batch>

# Run seed master
task seed-master

# Generate mocks for testing
task generate-mock

# Run tests with coverage
task test

# Run linting checks
task lint

# Auto-fix linting issues where possible
task lint-fix

# modernize 
task modernize-fix

# Swagger docs generate 
task docs-generate

```

### Examples

```bash
# Create a new migration
task migrate-create
# Enter migration name: create_new_feature_table

# Run migrations
task migrate-up

# Seed master data
task seed-master
```

## IV. Getting Started

1. Start services with Docker Compose:

2. Run migrations and seed the database:

```bash
task migrate-up
task seed-master
```

3. Access:
   - API: http://localhost:3011
   - Swagger UI: http://localhost:3011/swagger/index.html
