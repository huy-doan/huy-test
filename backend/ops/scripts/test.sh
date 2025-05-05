#!/bin/bash

# Colors for better readability
COLOR_RESET="\033[0m"
COLOR_GREEN="\033[32m"
COLOR_RED="\033[31m"

# MySQL connection settings
MYSQL_HOST=${DB_HOST:-"mysql"}
MYSQL_USER=${DB_USER:-"root"}
MYSQL_PASSWORD=${DB_PASSWORD:-"rootpw"}
MYSQL_DATABASE=${DB_NAME:-"msp-db-test"}
MYSQL_CONTAINER=${MYSQL_CONTAINER:-"makeshop_payment_mysql_1"}
BACKEND_CONTAINER=${BACKEND_CONTAINER:-"makeshop_payment_backend_1"}
MYSQL_CONNECT="${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:3306)/${MYSQL_DATABASE}"

# Get the backend directory path (parent of ops/scripts)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$(dirname "$SCRIPT_DIR")")"

# Set up error handling
set -e
trap 'echo -e "${COLOR_RED}Test failed${COLOR_RESET}"' ERR

echo -e "${COLOR_GREEN}Switching to test environment...${COLOR_RESET}"
${SCRIPT_DIR}/switch-env.sh test

echo -e "${COLOR_GREEN}Creating test database...${COLOR_RESET}"
docker exec -it ${MYSQL_CONTAINER} sh -c "mysql -u ${MYSQL_USER} -p${MYSQL_PASSWORD} -e 'CREATE DATABASE IF NOT EXISTS \`${MYSQL_DATABASE}\`'"

echo -e "${COLOR_GREEN}Running migrations...${COLOR_RESET}"
docker exec -it ${BACKEND_CONTAINER} goose -dir ./database/migrations mysql "${MYSQL_CONNECT}" up

echo -e "${COLOR_GREEN}Applying seed data...${COLOR_RESET}"
docker exec -it ${BACKEND_CONTAINER} goose --no-versioning -dir ./database/seeds/master mysql "${MYSQL_CONNECT}" up

echo -e "${COLOR_GREEN}Running tests...${COLOR_RESET}"
docker exec -it ${BACKEND_CONTAINER} sh -c "cd ./src/ && go test -v -covermode=set -coverprofile=coverage.out.tmp -coverpkg=./... ./tests/..."

echo -e "${COLOR_GREEN}Filtering coverage output...${COLOR_RESET}"
docker exec -it ${BACKEND_CONTAINER} sh -c "cd ./src/ && cat coverage.out.tmp | grep -v /tests > coverage.out"

echo -e "${COLOR_GREEN}Generating coverage report...${COLOR_RESET}"
docker exec -it ${BACKEND_CONTAINER} sh -c "cd ./src/ && go tool cover -html=coverage.out -o coverage.html"

echo -e "${COLOR_GREEN}Tests completed successfully!${COLOR_RESET}"

echo -e "${COLOR_GREEN}Switching back to development environment...${COLOR_RESET}"
${SCRIPT_DIR}/switch-env.sh dev

echo -e "${COLOR_GREEN}All done!${COLOR_RESET}"
