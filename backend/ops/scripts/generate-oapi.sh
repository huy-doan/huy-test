#!/bin/bash

# Set up environment variables
SPEC_DIR="./docs/api"
OUTPUT_DIR="./internal/pkg/api/generated"
PKG_NAME="generated"
BACKEND_CONTAINER=${BACKEND_CONTAINER:-"makeshop_payment_backend_1"}

# Create the output directory if it doesn't exist
mkdir -p $OUTPUT_DIR

# Create a temporary directory for the generated code
docker exec -it ${BACKEND_CONTAINER} oapi-codegen -package $PKG_NAME \
  -generate types,server,spec \
  -o "$OUTPUT_DIR/api.gen.go" \
  "$SPEC_DIR/swagger.yaml"

echo "API server code generated successfully!"
