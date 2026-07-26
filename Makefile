.PHONY: help build test test-unit test-integration test-chaos test-pact clean docker-build docker-push kind-setup kind-delete deploy

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s%s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all services
	@echo "Building Go service..."
	cd services/go-api-gateway && go build -o api-gateway ./cmd/main.go
	@echo "Building Rust service..."
	cd services/rust-data-store && cargo build --release
	@echo "Building Python service..."
	@echo "Python service builds at runtime"

test: ## Run all tests
	@echo "Running all tests..."
	$(MAKE) test-unit
	$(MAKE) test-integration
	$(MAKE) test-pact

test-unit: ## Run unit tests with coverage
	@echo "Running Go unit tests with coverage..."
	cd services/go-api-gateway && go test -cover -coverprofile=coverage.out ./...
	@echo "Running Rust unit tests with coverage..."
	cd services/rust-data-store && cargo test --release
	@echo "Running Python unit tests with coverage..."
	cd services/python-telemetry && pytest --cov=src --cov-report=html --cov-report=term

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	cd tests/integration && pytest -v

test-chaos: ## Run chaos engineering tests
	@echo "Running chaos engineering tests..."
	cd tests/chaos && pytest -v

test-pact: ## Run Pact contract tests
	@echo "Running Pact contract tests..."
	cd tests/pact && pytest -v

test-performance: ## Run k6 load performance scripts
	@echo "Running k6 load performance tests..."
	k6 run tests/load/load-test.js
	k6 run tests/load/scenarios/spike-test.js

clean: ## Clean build artifacts
	@echo "Cleaning Go build artifacts..."
	cd services/go-api-gateway && go clean
	@echo "Cleaning Rust build artifacts..."
	cd services/rust-data-store && cargo clean
	@echo "Cleaning Python cache..."
	cd services/python-telemetry && find . -type d -name __pycache__ -exec rm -rf {} + 2>/dev/null || true

docker-build: ## Build Docker images
	@echo "Building Docker images..."
	docker-compose build

docker-push: ## Push Docker images to registry
	@echo "Pushing Docker images to registry..."
	docker push ghcr.io/username/ci-cd-pipeline/go-api-gateway:latest
	docker push ghcr.io/username/ci-cd-pipeline/rust-data-store:latest
	docker push ghcr.io/username/ci-cd-pipeline/python-telemetry:latest

kind-setup: ## Setup Kind cluster
	@echo "Setting up Kind cluster..."
	bash scripts/setup.sh

kind-delete: ## Delete Kind cluster
	@echo "Deleting Kind cluster..."
	kind delete cluster --name ci-cd-pipeline

deploy: ## Deploy to Kind cluster
	@echo "Deploying to Kind cluster..."
	bash scripts/load-images.sh
	kubectl apply -f k8s/base/go-gateway/
	kubectl apply -f k8s/base/rust-store/
	kubectl apply -f k8s/base/python-telemetry/

logs: ## Show logs from all services
	@echo "Showing logs from Go API Gateway..."
	kubectl logs -f deployment/go-api-gateway

logs-rust: ## Show logs from Rust data store
	@echo "Showing logs from Rust data store..."
	kubectl logs -f deployment/rust-data-store

logs-python: ## Show logs from Python telemetry
	@echo "Showing logs from Python telemetry..."
	kubectl logs -f deployment/python-telemetry

status: ## Show status of all deployments
	@echo "Deployment status:"
	kubectl get deployments
	@echo ""
	@echo "Pod status:"
	kubectl get pods
	@echo ""
	@echo "Service status:"
	kubectl get services

port-forward: ## Port forward Go API gateway to localhost:8080
	@echo "Port forwarding Go API gateway to localhost:8080..."
	kubectl port-forward svc/go-api-gateway 8080:8080

install-argocd: ## Install Argo CD
	@echo "Installing Argo CD..."
	bash scripts/install-argocd.sh

test-coverage: ## Generate combined coverage report
	@echo "Generating coverage reports..."
	cd services/go-api-gateway && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage reports generated:"
	@echo "  Go: services/go-api-gateway/coverage.html"
	@echo "  Python: services/python-telemetry/htmlcov/index.html"

test-smoke: ## Run smoke tests without full infrastructure
	@echo "Running smoke tests..."
	cd services/go-api-gateway && go test -run TestSmoke ./...
	cd services/rust-data-store && cargo test --test smoke
	cd services/python-telemetry && pytest -k smoke
