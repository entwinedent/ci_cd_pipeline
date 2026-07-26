# Protocol Buffers Definitions

This directory contains the Protocol Buffers (protobuf) definitions for the gRPC services in the CI/CD pipeline.

## Services

### Data Store Service

The Rust Data Store provides a gRPC interface for key-value storage operations.

**File**: `data_store.proto`

**Methods**:
- `Set` - Store a key-value pair with optional TTL
- `Get` - Retrieve a value by key
- `Delete` - Remove a key from the store
- `HealthCheck` - Check service health

## Compilation

### Go

```bash
# Install protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate Go code
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/data_store.proto
```

### Rust

```bash
# Install protoc and protoc-gen-rust
cargo install protobuf-codegen

# Generate Rust code
protoc --rust_out=. proto/data_store.proto
```

### Python

```bash
# Install grpcio-tools
pip install grpcio-tools

# Generate Python code
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. proto/data_store.proto
```

## Usage

### Go Client

```go
import (
    "github.com/username/ci-cd-pipeline/proto/datastore"
    "google.golang.org/grpc"
)

conn, err := grpc.Dial("rust-data-store:50051", grpc.WithInsecure())
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

client := datastore.NewDataStoreServiceClient(conn)

resp, err := client.Set(ctx, &datastore.SetRequest{
    Key:   []byte("test"),
    Value: []byte("hello"),
})
```

### Rust Server

```rust
use datastore::data_store_server::{DataStoreService, DataStoreServiceServer};
use tonic::{Request, Response, Status};

struct DataStore;

impl DataStoreService for DataStore {
    async fn set(&self, request: Request<SetRequest>) -> Result<Response<SetResponse>, Status> {
        // Implementation
    }
}
```

### Python Client

```python
import grpc
from proto import data_store_pb2, data_store_pb2_grpc

channel = grpc.insecure_channel('rust-data-store:50051')
stub = data_store_pb2_grpc.DataStoreServiceStub(channel)

response = stub.Set(data_store_pb2.SetRequest(
    key=b'test',
    value=b'hello'
))
```

## Testing

### Unit Tests

Test protobuf message serialization:

```bash
# Go
go test ./proto/...

# Rust
cargo test --package proto

# Python
pytest tests/proto/
```

## Schema Evolution

When modifying protobuf definitions:

1. Follow protobuf best practices
2. Maintain backward compatibility
3. Update all language bindings
4. Test cross-language compatibility
5. Document breaking changes

## Documentation

- [Protocol Buffers Documentation](https://developers.google.com/protocol-buffers)
- [gRPC Documentation](https://grpc.io/docs/)
- [Language Guides](https://grpc.io/docs/languages/)
