#!/bin/bash

# Development setup script for vybium-crypto
# Sets up the complete development environment

set -e

echo "🚀 Setting up vybium-crypto development environment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    print_error "Not in a Go project directory. Please run from vybium-crypto root."
    exit 1
fi

# Check Go version
print_status "Checking Go version..."
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
print_success "Go version: $GO_VERSION"

# Install Go dependencies
print_status "Installing Go dependencies..."
go mod download
go mod verify
print_success "Go dependencies installed"

# Install development tools
print_status "Installing development tools..."

# golangci-lint
if ! command -v golangci-lint &> /dev/null; then
    print_status "Installing golangci-lint..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
    print_success "golangci-lint installed"
else
    print_success "golangci-lint already installed"
fi

# gosec
if ! command -v gosec &> /dev/null; then
    print_status "Installing gosec..."
    go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest || print_warning "gosec installation failed, continuing without it"
    print_success "gosec installed"
else
    print_success "gosec already installed"
fi

# staticcheck
if ! command -v staticcheck &> /dev/null; then
    print_status "Installing staticcheck..."
    go install honnef.co/go/tools/cmd/staticcheck@latest
    print_success "staticcheck installed"
else
    print_success "staticcheck already installed"
fi

# pre-commit
if ! command -v pre-commit &> /dev/null; then
    print_status "Installing pre-commit..."
    if command -v pip &> /dev/null; then
        pip install pre-commit
    elif command -v pip3 &> /dev/null; then
        pip3 install pre-commit
    else
        print_warning "pip not found, please install pre-commit manually"
    fi
    print_success "pre-commit installed"
else
    print_success "pre-commit already installed"
fi

# Install pre-commit hooks
print_status "Installing pre-commit hooks..."
if command -v pre-commit &> /dev/null; then
    pre-commit install
    print_success "Pre-commit hooks installed"
else
    print_warning "pre-commit not available, skipping hook installation"
fi

# Create necessary directories
print_status "Creating necessary directories..."
mkdir -p docs scripts
print_success "Directories created"

# Run initial checks
print_status "Running initial checks..."

# Format code
print_status "Formatting code..."
gofmt -s -w ./pkg/vybium-crypto/...
if command -v goimports &> /dev/null; then
    goimports -w ./pkg/vybium-crypto/...
fi
print_success "Code formatted"

# Run tests
print_status "Running tests..."
go test -v ./pkg/vybium-crypto/...
print_success "Tests passed"

# Run linters
print_status "Running linters..."
golangci-lint run --timeout=5m
print_success "Linting passed"

# Generate documentation
print_status "Generating documentation..."
go doc -all ./pkg/vybium-crypto/... > docs/api.md
print_success "Documentation generated"

print_success "🎉 Development environment setup completed!"
echo ""
echo "📋 Available commands:"
echo "  make help          - Show all available commands"
echo "  make test          - Run tests"
echo "  make lint          - Run linters"
echo "  make format        - Format code"
echo "  make security       - Run security checks"
echo "  make ci            - Run full CI/CD pipeline"
echo "  make dev-setup     - Complete development setup"
echo "  make quick-check   - Quick development check"
echo "  make full-check    - Full quality check"
echo ""
echo "🔧 Pre-commit hooks are installed and will run automatically on git commit"
echo "🚀 You can now start developing!"
