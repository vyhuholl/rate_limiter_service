.PHONY: all build test clean lint fmt help

# Default target
all: lint test build

# Build the project
build:
	go build -o bin/rate_limiter_service ./cmd/...

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-cov: test
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html

# Lint the code
lint:
	golangci-lint run ./...

# Format the code
fmt:
	gofmt -s -w .
	goimports -w .

# Check code formatting
fmt-check:
	gofmt -s -d .
	goimports -d .

# Run both lint and tests
check: lint test

# Install development dependencies
dev-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

# Show help
help:
	@echo "Available targets:"
	@echo "  all       - Run lint, test, and build"
	@echo "  build     - Build the project"
	@echo "  test      - Run tests with race detector and coverage"
	@echo "  test-cov  - Run tests and open coverage report"
	@echo "  clean     - Clean build artifacts"
	@echo "  lint      - Lint the code"
	@echo "  fmt       - Format the code"
	@echo "  fmt-check - Check code formatting"
	@echo "  check     - Run both lint and test"
	@echo "  dev-deps  - Install development dependencies"
	@echo "  help      - Show this help message"