#!/bin/bash

# Function to display colored messages
print_message() {
  local color=$1
  local message=$2
  
  case $color in
    "green") echo -e "\033[0;32m$message\033[0m" ;;
    "yellow") echo -e "\033[0;33m$message\033[0m" ;;
    "red") echo -e "\033[0;31m$message\033[0m" ;;
    *) echo "$message" ;;
  esac
}

# Get database connection info from environment variables
DB_HOST=${DB_HOST}
DB_PORT=${DB_PORT}
DB_USER=${DB_USER}
DB_PASSWORD=${DB_PASSWORD}
DB_NAME=${DB_NAME}

# Construct the connection string
CONNECTION_STRING="${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}"

# Define paths
MIGRATIONS_DIR="./config/migrations"
SEEDS_DIR="./config/seeds/master"

# Function to execute goose commands
run_goose() {
  local operation=$1
  local dir=$2
  local extra_args=$3
  local conn_string=$4
  
  print_message "yellow" "Running: goose $extra_args -dir $dir mysql \"$conn_string\" $operation"
  goose $extra_args -dir $dir mysql "$conn_string" $operation
  
  # Check if command was successful
  if [ $? -eq 0 ]; then
    print_message "green" "✅ Success: $operation completed for $dir"
  else
    print_message "red" "❌ Error: $operation failed for $dir"
    return 1
  fi
}

# Function to run migrations
run_migrations() {
  local operation=$1
  print_message "yellow" "===== DATABASE MIGRATIONS: $operation ====="
  run_goose "$operation" "$MIGRATIONS_DIR" "" "$CONNECTION_STRING"
}

# Function to run seeds
run_seeds() {
  local operation=$1
  print_message "yellow" "===== DATABASE SEEDS: $operation ====="
  run_goose "$operation" "$SEEDS_DIR" "--no-versioning" "$CONNECTION_STRING"
}

# Main execution logic
case "$1" in
  "migrate:up")
    run_migrations "up"
    ;;
  "migrate:down")
    run_migrations "down"
    ;;
  "seed:up")
    run_seeds "up"
    ;;
  "seed:down")
    run_seeds "down"
    ;;
  "status")
    print_message "yellow" "===== MIGRATION STATUS ====="
    goose -dir "$MIGRATIONS_DIR" mysql "$CONNECTION_STRING" status
    ;;
  *)
    echo "Usage: $0 {migrate:up|migrate:down|seed:up|seed:down|reset|status}"
    echo "  migrate:up   - Run all pending migrations"
    echo "  migrate:down - Rollback last migration"
    echo "  seed:up      - Apply all seed data"
    echo "  seed:down    - Remove seed data"
    echo "  status       - Show migration status"
    exit 1
    ;;
esac

exit 0
