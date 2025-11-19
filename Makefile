# Makefile for vybium-crypto
# Provides convenient commands for local development
# Standardized for all Vybium projects

.PHONY: help install test test-race test-coverage benchmark lint format format-check security build clean ci pre-commit install-hooks dev-setup dev-deps tidy download vuln-check check

# Default target
help: ## Show this help message
	@echo "vybium-crypto Development Commands"
	@echo "================================="
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Installation
install: ## Install dependencies and tools
	@echo "Installing Go dependencies..."
	go mod download
	go mod verify
	@echo "Installing development tools..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2; \
	fi
	@if ! command -v gosec &> /dev/null; then \
		echo "Installing gosec..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest || echo "gosec installation failed, continuing without it"; \
	fi
	@if ! command -v staticcheck &> /dev/null; then \
		echo "Installing staticcheck..."; \
		go install honnef.co/go/tools/cmd/staticcheck@latest; \
	fi
	@if ! command -v pre-commit &> /dev/null; then \
		echo "Installing pre-commit..."; \
		pip install pre-commit; \
	fi

# Testing
test: ## Run tests
	@echo "Running tests..."
	go test -v ./pkg/vybium-crypto/...

test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	go test -v -race ./pkg/vybium-crypto/...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./pkg/vybium-crypto/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./pkg/vybium-crypto/... > benchmark_results.txt 2>&1
	@echo "Benchmark results saved to: benchmark_results.txt"

# Code quality
lint: ## Run linters
	@echo "Running golangci-lint..."
	golangci-lint run --timeout=5m
	@echo "Running staticcheck..."
	staticcheck ./pkg/vybium-crypto/...
	@echo "Running go vet..."
	go vet ./pkg/vybium-crypto/...

format: ## Format code
	@echo "Formatting code..."
	@if command -v gofumpt >/dev/null 2>&1; then \
		gofumpt -w .; \
	else \
		echo "Warning: gofumpt not found. Install with: go install mvdan.cc/gofumpt@latest"; \
		gofmt -s -w ./pkg/vybium-crypto/...; \
	fi
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w ./pkg/vybium-crypto/...; \
	else \
		echo "Warning: goimports not found. Install with: go install golang.org/x/tools/cmd/goimports@latest"; \
		go fmt ./pkg/vybium-crypto/...; \
	fi
	@echo "✓ Code formatted"

format-check: ## Check code formatting
	@echo "Checking code formatting..."
	@if command -v gofumpt >/dev/null 2>&1; then \
		unformatted=$$(gofumpt -l . | grep -v '^vendor/' | wc -l); \
		if [ $$unformatted -gt 0 ]; then \
			echo "The following files are not formatted:"; \
			gofumpt -l . | grep -v '^vendor/'; \
			echo ""; \
			echo "Run 'make format' to fix formatting"; \
			exit 1; \
		fi; \
	else \
		unformatted=$$(gofmt -l ./pkg/vybium-crypto/... | wc -l); \
		if [ $$unformatted -gt 0 ]; then \
			echo "The following files are not formatted:"; \
			gofmt -l ./pkg/vybium-crypto/...; \
			echo ""; \
			echo "Run 'make format' to fix formatting"; \
			exit 1; \
		fi; \
	fi
	@echo "✓ All files are formatted"

security: ## Run security checks
	@echo "Running gosec security scanner..."
	gosec -fmt sarif -out gosec.sarif ./pkg/vybium-crypto/...
	@echo "Security scan completed: gosec.sarif"

# Build
build: ## Build the project
	@echo "Building project..."
	go build ./pkg/vybium-crypto/...
	@echo "Build completed"

# Documentation
docs: ## Generate documentation
	@echo "Generating API documentation..."
	mkdir -p docs
	go doc -all ./pkg/vybium-crypto/... > docs/api.md
	@echo "API documentation generated: docs/api.md"

# Pre-commit hooks
install-hooks: ## Install pre-commit hooks
	@echo "Installing pre-commit hooks..."
	pre-commit install
	@echo "Pre-commit hooks installed"

pre-commit: ## Run pre-commit hooks on all files
	@echo "Running pre-commit hooks..."
	pre-commit run --all-files

# CI/CD
ci: ## Run full CI/CD pipeline locally
	@echo "Running local CI/CD pipeline..."
	./scripts/local-ci.sh

# Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@go clean -cache -testcache
	@rm -f coverage.out coverage.html benchmark_results.txt gosec.sarif
	@rm -rf bin/ reports/
	@echo "✓ Build artifacts cleaned"

# Development workflow
dev-setup: tidy download dev-deps install-hooks ## Complete development setup
	@echo "✓ Development environment ready!"

dev-deps: ## Install development dependencies
	@echo "Installing development dependencies..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest || echo "golangci-lint already installed or failed"
	@go install mvdan.cc/gofumpt@latest || echo "gofumpt already installed or failed"
	@go install golang.org/x/tools/cmd/goimports@latest || echo "goimports already installed or failed"
	@go install golang.org/x/vuln/cmd/govulncheck@latest || echo "govulncheck already installed or failed"
	@go install honnef.co/go/tools/cmd/staticcheck@latest || echo "staticcheck already installed or failed"
	@if ! command -v gosec &> /dev/null; then \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest || echo "gosec installation failed"; \
	fi
	@if ! command -v pre-commit &> /dev/null; then \
		echo "Installing pre-commit..."; \
		pip install pre-commit || echo "pre-commit installation failed"; \
	fi
	@echo "✓ Development dependencies installed"

tidy: ## Run go mod tidy
	@echo "Running go mod tidy..."
	@go mod tidy
	@echo "✓ Dependencies cleaned"

download: ## Download Go modules
	@echo "Downloading Go modules..."
	@go mod download
	@echo "✓ Modules downloaded"

dev: format lint test ## Quick dev cycle (format + lint + test)
	@echo "✓ Development cycle complete!"

# Security
vuln-check: ## Check for vulnerabilities
	@echo "Checking for vulnerabilities..."
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./pkg/vybium-crypto/...; \
	else \
		echo "Installing govulncheck..."; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
		govulncheck ./pkg/vybium-crypto/...; \
	fi

# Quick checks
check: format-check lint test ## Run all checks (format + lint + test)
	@echo "✓ All checks passed!"

# Full pipeline
full-check: format lint test-race test-coverage security benchmark ## Full quality check
	@echo "✓ Full quality check completed"
