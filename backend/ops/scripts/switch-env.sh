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
  
  # Create a more persistent approach - modify Docker container environment
  docker exec -it ${BACKEND_CONTAINER} sh -c '
    # Create a clean environment file without comments
    grep -v "^#" /app/ops/development/.test.env | grep -v "^$" > /etc/test-environment
    
    # Remove any existing environment file symlink
    rm -f /etc/environment.sh
    
    # Create a script that will be sourced by all subsequent shell sessions
    echo "#!/bin/sh" > /etc/environment.sh
    echo "# Auto-generated environment variables for TEST environment" >> /etc/environment.sh
    
    # Add each variable from the clean file to the environment script
    while IFS= read -r line; do
        # Skip empty lines or lines without an equal sign
        if [[ -z "$line" || "$line" != *=* ]]; then
            continue
        fi

        # Extract variable name and value
        varname=$(echo "$line" | cut -d= -f1)  # Remove any leading/trailing spaces
        varvalue=$(echo "$line" | cut -d= -f2-)

        echo "Debug: $varname=$varvalue"

        # Write the export statement to /etc/environment.sh
        echo "export $varname=\"$varvalue\"" >> /etc/environment.sh

        # Also export the variable in the current shell
        export "$varname=$varvalue"
    done < /etc/test-environment
    
    # Make the script executable
    chmod +x /etc/environment.sh
    
    # Create a profile.d script to ensure the environment is loaded for all users
    echo "#!/bin/sh" > /etc/profile.d/env.sh
    echo ". /etc/environment.sh" >> /etc/profile.d/env.sh
    chmod +x /etc/profile.d/env.sh
    
    # Create symlink to /etc/environment for compatibility
    ln -sf /etc/environment.sh /etc/environment
    
    # Source the environment script in the current shell
    . /etc/environment.sh
    
    # Print confirmation
    echo "Test environment variables loaded successfully"
  '
  
  # Verify that the environment was set correctly
  echo -e "${COLOR_YELLOW}Environment switched to TEST${COLOR_RESET}"
  echo -e "${COLOR_YELLOW}Current database configuration:${COLOR_RESET}"
  docker exec -it ${BACKEND_CONTAINER} sh -c "source /etc/environment.sh && printenv | grep DB_"

elif [ "$1" == "dev" ]; then
  echo -e "${COLOR_GREEN}Switching to development environment${COLOR_RESET}"
  
  # Create a more persistent approach - modify Docker container environment
  docker exec -it ${BACKEND_CONTAINER} sh -c '
    # Create a clean environment file without comments
    grep -v "^#" /app/ops/development/.env | grep -v "^$" > /etc/dev-environment
    
    # Remove any existing environment file symlink
    rm -f /etc/environment.sh
    
    # Create a script that will be sourced by all subsequent shell sessions
    echo "#!/bin/sh" > /etc/environment.sh
    echo "# Auto-generated environment variables for DEV environment" >> /etc/environment.sh
    
    # Add each variable from the clean file to the environment script
    while IFS= read -r line; do
      # Skip empty lines or lines without an equal sign
        if [[ -z "$line" || "$line" != *=* ]]; then
            continue
        fi

        # Extract variable name and value
        varname=$(echo "$line" | cut -d= -f1)  # Remove any leading/trailing spaces
        varvalue=$(echo "$line" | cut -d= -f2-)

        echo "Debug: $varname=$varvalue"

        # Write the export statement to /etc/environment.sh
        echo "export $varname=\"$varvalue\"" >> /etc/environment.sh

        # Also export the variable in the current shell
        export "$varname=$varvalue"
    done < /etc/dev-environment
    
    # Make the script executable
    chmod +x /etc/environment.sh
    
    # Create a profile.d script to ensure the environment is loaded for all users
    echo "#!/bin/sh" > /etc/profile.d/env.sh
    echo ". /etc/environment.sh" >> /etc/profile.d/env.sh
    chmod +x /etc/profile.d/env.sh
    
    # Create symlink to /etc/environment for compatibility
    ln -sf /etc/environment.sh /etc/environment
    
    # Source the environment script in the current shell
    . /etc/environment.sh
    
    # Print confirmation
    echo "Development environment variables loaded successfully"
  '
  
  # Verify that the environment was set correctly
  echo -e "${COLOR_YELLOW}Environment switched to DEV${COLOR_RESET}"
  echo -e "${COLOR_YELLOW}Current database configuration:${COLOR_RESET}"
  docker exec -it ${BACKEND_CONTAINER} sh -c "source /etc/environment.sh && printenv | grep DB_"
  
else
  echo -e "${COLOR_YELLOW}Invalid environment: $1${COLOR_RESET}"
  usage
fi
