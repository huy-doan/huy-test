#!/bin/bash

# Set định nghĩa đường dẫn
SPEC_DIR="./docs/api"
OUTPUT_DIR="./docs/api/generated"
PKG_NAME="generated"
BACKEND_CONTAINER=${BACKEND_CONTAINER:-"makeshop_payment_backend_1"}

# Tạo thư mục nếu chưa tồn tại
mkdir -p $OUTPUT_DIR

# Tạo file server code từ OpenAPI specification
docker exec -it ${BACKEND_CONTAINER} oapi-codegen -package $PKG_NAME \
  -generate types,server,spec \
  -o "$OUTPUT_DIR/api.gen.go" \
  "$SPEC_DIR/swagger.yaml"

echo "API server code generated successfully!"
