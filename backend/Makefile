.PHONY: help run swagger fmt migrate-create migrate-up migrate-down shell lint lint-fix

# Variables
API_NAME = msp-api
API_BIN = ./bin/api
MAIN_FILE = ./src/main.go
LOCAL_PORT = 8080
DOCKER_PORT = 3010

# Tools
GO = go

# Database
MYSQL_DB = msp-db-dev
MYSQL_DB_TEST = msp-db-test
MYSQL_USER = root
MYSQL_PASSWORD = rootpw
MYSQL_HOST = mysql
MYSQL_PORT = 3306
MYSQL_CONTAINER = makeshop_payment_mysql_1

MYSQL_CONNECT = "$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DB)"
BACKEND_CONTAINER = makeshop_payment_backend_1

# Colors
COLOR_RESET = \033[0m
COLOR_GREEN = \033[32m
COLOR_YELLOW = \033[33m

help: ## Display available commands
	@echo "$(COLOR_GREEN)Available commands:$(COLOR_RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(COLOR_YELLOW)%-20s$(COLOR_RESET) %s\n", $$1, $$2}'

ssh-be: ## SSH to the container
	@echo "$(COLOR_GREEN)SSH to the BACKEND...$(COLOR_RESET)"
	docker exec -it $(BACKEND_CONTAINER) sh

ssh-mysql: ## SSH to the container
	@echo "$(COLOR_GREEN)SSH to the MYSQL...$(COLOR_RESET)"
	docker exec -it makeshop_payment_mysql_1 sh

run: ## Run the API locally
	@echo "$(COLOR_GREEN)Running $(API_NAME) on http://localhost:$(DOCKER_PORT)...$(COLOR_RESET)"
	docker exec -it $(BACKEND_CONTAINER) go run $(MAIN_FILE)

swagger: ## Generate Swagger documentation
	@echo "$(COLOR_GREEN)Generating Swagger documentation...$(COLOR_RESET)"
	docker exec -it $(BACKEND_CONTAINER) swag init -g src/main.go
	@echo "$(COLOR_GREEN)Swagger documentation generated!$(COLOR_RESET)"

fmt: ## Format code
	@echo "$(COLOR_GREEN)Formatting code...$(COLOR_RESET)"
	docker exec -it $(BACKEND_CONTAINER) $(GO) fmt ./...
	@echo "$(COLOR_GREEN)Formatting complete!$(COLOR_RESET)"

migrate-create: ## Run database migrations create
	@echo "$(COLOR_GREEN)Running create Migrations...$(COLOR_RESET)"
	@read -p "Enter migration name: " name; \
   	docker exec -it $(BACKEND_CONTAINER) goose -dir ./database/migrations create $$name sql
	@echo "$(COLOR_GREEN)Create Migrations complete!$(COLOR_RESET)"

migrate-up: ## Run database migrations up
	@echo "$(COLOR_GREEN)Create database if not exists$(COLOR_RESET)"
	docker exec -it $(MYSQL_CONTAINER) mysql -u$(MYSQL_USER) -p$(MYSQL_PASSWORD) -e 'CREATE DATABASE IF NOT EXISTS `$(MYSQL_DB)` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci'
	@echo "$(COLOR_GREEN)Running database migrations up...$(COLOR_RESET)"
	docker exec -it $(BACKEND_CONTAINER) goose -dir ./database/migrations mysql $(MYSQL_CONNECT) up
	@echo "$(COLOR_GREEN)Migrations complete!$(COLOR_RESET)"

migrate-down: ## Run database migrations down
	@echo "$(COLOR_GREEN)Running database migrations down...$(COLOR_RESET)"
	docker exec -it $(BACKEND_CONTAINER) goose -dir ./database/migrations mysql $(MYSQL_CONNECT) down
	@echo "$(COLOR_GREEN)Migrations complete!$(COLOR_RESET)"

seed-master: ## Run database seed-master
	@echo "$(COLOR_GREEN)Running seeds master...$(COLOR_RESET)"
	docker exec -it $(BACKEND_CONTAINER) goose --no-versioning -dir ./database/seeds/master mysql $(MYSQL_CONNECT) up
	@echo "$(COLOR_GREEN)Run seeds master complete!$(COLOR_RESET)"

shell: ## Run shell in the container
	@$(eval ARGS := $(filter-out $@,$(MAKECMDGOALS)))
	@echo "$(COLOR_GREEN)Running Shell...$(COLOR_RESET)"
	docker exec -it $(BACKEND_CONTAINER) go run main.go $(ARGS)
	@echo "$(COLOR_GREEN)Running Shell complete!$(COLOR_RESET)"

generate-mock: ## Generate mock for interfaces
	@echo "$(COLOR_GREEN)Generating mock...$(COLOR_RESET)"
	docker exec -it $(BACKEND_CONTAINER) sh -c "mockery --case=underscore --dir=src --output=src/tests/mocks/ --all"
	@echo "$(COLOR_GREEN)Mock generation complete!$(COLOR_RESET)"

lint: ## Run golangci-lint checks
	@echo "$(COLOR_GREEN)Running golangci-lint...$(COLOR_RESET)"
	docker exec -it $(BACKEND_CONTAINER) golangci-lint run ./...
	@echo "$(COLOR_GREEN)Lint check complete!$(COLOR_RESET)"

lint-fix: ## Fix lint issues automatically where possible
	@echo "$(COLOR_GREEN)Running golangci-lint with autofix...$(COLOR_RESET)"
	docker exec -it $(BACKEND_CONTAINER) golangci-lint run --fix ./...
	@echo "$(COLOR_GREEN)Lint fix complete!$(COLOR_RESET)"

COVERMODE=set

test: ## Run tests with coverage
	@chmod +x ./ops/scripts/switch-env.sh
	@chmod +x ./ops/scripts/test.sh
	@./ops/scripts/test.sh
