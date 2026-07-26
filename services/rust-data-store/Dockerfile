# Build stage
FROM rust:latest AS builder

WORKDIR /app

# Copy Cargo files and source
COPY Cargo.toml ./
COPY src ./src
RUN cargo generate-lockfile

# Build release binary
RUN cargo build --release

# Final stage
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy binary from builder
COPY --from=builder /app/target/release/rust-data-store /rust-data-store

# Expose HTTP port
EXPOSE 50051

# Run binary
ENTRYPOINT ["/rust-data-store"]