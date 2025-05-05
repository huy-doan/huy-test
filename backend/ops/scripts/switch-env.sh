#!/bin/bash

# Colors for better readability
COLOR_RESET="\033[0m"
COLOR_GREEN="\033[32m"
COLOR_YELLOW="\033[33m"
COLOR_BLUE="\033[34m"

# Get the backend directory path (parent of ops/scripts)
BACKEND_CONTAINER=${BACKEND_CONTAINER:-"makeshop_payment_backend_1"}

# Function to display usage information
usage() {
  echo -e "${COLOR_BLUE}Usage:${COLOR_RESET} $0 [dev|test]"
  echo ""
  echo -e "${COLOR_BLUE}Options:${COLOR_RESET}"
  echo -e "  dev    Switch to development environment"
  echo -e "  test   Switch to test environment"
  echo ""
  exit 1
}

# Check if environment parameter is provided
if [ $# -ne 1 ]; then
  usage
fi

# Switch environment based on parameter
if [ "$1" == "test" ]; then
  echo -e "${COLOR_GREEN}Switching to test environment${COLOR_RESET}"
  docker exec -it ${BACKEND_CONTAINER} sh -c "export $(grep -v '^#' /app/ops/development/.test.env | xargs)"
  echo -e "${COLOR_YELLOW}Environment switched to TEST${COLOR_RESET}"
  echo -e "${COLOR_YELLOW}Current database configuration:${COLOR_RESET}"
  docker exec -it ${BACKEND_CONTAINER} sh -c "printenv | grep DB_"
elif [ "$1" == "dev" ]; then
  echo -e "${COLOR_GREEN}Switching to development environment${COLOR_RESET}"
  docker exec -it ${BACKEND_CONTAINER} sh -c "export $(grep -v '^#' /app/ops/development/.env | xargs)"
  echo -e "${COLOR_YELLOW}Environment switched to DEV${COLOR_RESET}"
  echo -e "${COLOR_YELLOW}Current database configuration:${COLOR_RESET}"
    docker exec -it ${BACKEND_CONTAINER} sh -c "printenv | grep DB_"
else
  echo -e "${COLOR_YELLOW}Invalid environment: $1${COLOR_RESET}"
  usage
fi