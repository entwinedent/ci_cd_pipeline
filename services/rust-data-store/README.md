# Rust Data Store

Ultra-low latency, thread-safe in-memory storage engine with gRPC interface and advanced concurrency patterns.

## Architecture

The Data Store provides:
- Thread-safe in-memory key-value storage
- gRPC interface for service communication
- Min-heap TTL eviction strategy
- DashMap for lock-free concurrent access
- Automatic garbage collection

## Concurrency Model

### DashMap vs Standard Locks

We chose **DashMap** over `RwLock<HashMap>` for several reasons:

**Performance Benefits**
- Lock-free reads with sharded architecture
- Eliminates lock contention in high-throughput scenarios
- Better cache locality with reduced false sharing
- O(1) average case for read-heavy workloads

**Memory Efficiency**
- Minimal overhead compared to lock-based alternatives
- Efficient memory allocation with arena-based storage
- Reduced memory fragmentation

**Implementation**
```rust
use dashmap::DashMap;

let store: DashMap<String, DataEntry> = DashMap::new();
// Lock-free reads
let value = store.get(&key).map(|entry| entry.value().clone());
// Optimistic writes
store.insert(key, entry);
```

## TTL Eviction Strategy

### Min-Heap Implementation

The data store uses a min-heap for efficient TTL-based eviction:

**Algorithm**
1. Each entry has an expiration timestamp
2. Entries are organized in a min-heap by expiration time
3. Background goroutine periodically checks heap
4. Expired entries are removed in O(log n) time

**Benefits**
- Efficient expiration checking
- Predictable memory usage
- No full garbage collection pauses
- Better performance than naive iteration

**Configuration**
- Default TTL: 5 minutes
- Cleanup interval: 30 seconds
- Max entries: 1,000,000

## gRPC Interface

### Service Definition

```protobuf
service DataStoreService {
  rpc Set(SetRequest) returns (SetResponse);
  rpc Get(GetRequest) returns (GetResponse);
  rpc Delete(DeleteRequest) returns (DeleteResponse);
  rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);
}
```

### Methods

**Set**
- Stores a key-value pair with optional TTL
- Returns success/error status
- Thread-safe concurrent writes

**Get**
- Retrieves value by key
- Returns not found if key expired
- Lock-free read operation

**Delete**
- Removes key from store
- Returns success/error status
- Atomic operation

**HealthCheck**
- Returns store statistics
- Entry count, memory usage
- Health status

## Configuration

### Environment Variables
- `PORT` - gRPC service port (default: 50051)
- `LOG_LEVEL` - Logging level (default: `info`)
- `DEFAULT_TTL_SECONDS` - Default TTL for entries (default: 300)
- `MAX_ENTRIES` - Maximum number of entries (default: 1000000)
- `CLEANUP_INTERVAL_SECONDS` - TTL cleanup interval (default: 30)

### Service Specifications
- **Port**: 50051
- **Protocol**: gRPC
- **Health**: gRPC HealthCheck service

## Development

### Build
```bash
cd services/rust-data-store
cargo build --release
```

### Run
```bash
cargo run
```

### Test
```bash
cargo test
cargo test --release
cargo test -- --nocapture
```

### Benchmark
```bash
cargo bench
```

### Docker Build
```bash
docker build -t rust-data-store:latest .
```

## Performance

- **Image Size**: < 25MB
- **gRPC Latency**: < 10ms (p95)
- **Throughput**: > 100,000 ops/sec
- **Memory**: < 100MB (1M entries)

## Dependencies

- `tokio` - Async runtime
- `dashmap` - Concurrent hash map
- `tonic` - gRPC framework
- `prost` - Protocol buffers
- `serde` - Serialization

## Security

- SPIFFE/SPIRE mTLS for gRPC connections
- Input validation on all operations
- Memory-safe Rust guarantees
- No unsafe code blocks
