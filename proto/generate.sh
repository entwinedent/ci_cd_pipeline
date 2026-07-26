#!/bin/bash
# Generate protobuf code for all languages

set -e

PROTO_DIR="proto"
OUTPUT_DIR="."

echo "Generating protobuf code..."

# Go
echo "Generating Go code..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ${PROTO_DIR}/data_store.proto

# Rust
echo "Generating Rust code..."
protoc --rust_out=. ${PROTO_DIR}/data_store.proto

# Python
echo "Generating Python code..."
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. ${PROTO_DIR}/data_store.proto

echo "Protobuf code generation complete!"
