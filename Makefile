# syno-vm Makefile

BINARY_NAME=syno-vm
VERSION?=0.1.0
BUILD_TIME=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(GIT_COMMIT) -X main.date=$(BUILD_TIME) -X main.builtBy=makefile"

# Platforms
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

.PHONY: build clean test deps help install uninstall release

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME) cmd/syno-vm/main.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/
	rm -rf dist/

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Install locally
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp bin/$(BINARY_NAME) /usr/local/bin/
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)

# Uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

# Build for multiple platforms
release: clean
	@echo "Building for multiple platforms..."
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d'/' -f1); \
		GOARCH=$$(echo $$platform | cut -d'/' -f2); \
		OUTPUT_NAME=$(BINARY_NAME)-$(VERSION)-$$GOOS-$$GOARCH; \
		if [ $$GOOS = "windows" ]; then \
			OUTPUT_NAME=$$OUTPUT_NAME.exe; \
		fi; \
		echo "Building $$OUTPUT_NAME..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH $(GOBUILD) $(LDFLAGS) -o dist/$$OUTPUT_NAME cmd/syno-vm/main.go; \
	done

# Development build (with race detector and debug info)
dev:
	@echo "Building development version..."
	$(GOBUILD) -race -o bin/$(BINARY_NAME)-dev cmd/syno-vm/main.go

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Security scan
security:
	@echo "Running security scan..."
	gosec ./...

# Run integration tests
integration-test:
	@echo "Running integration tests..."
	$(GOTEST) -v -tags=integration ./test/integration/...

# Generate coverage report
coverage:
	@echo "Generating coverage report..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  build           - Build the binary"
	@echo "  clean           - Clean build artifacts"
	@echo "  test            - Run tests"
	@echo "  integration-test- Run integration tests"
	@echo "  coverage        - Generate coverage report"
	@echo "  bench           - Run benchmarks"
	@echo "  deps            - Download dependencies"
	@echo "  install         - Install binary to /usr/local/bin"
	@echo "  uninstall       - Remove binary from /usr/local/bin"
	@echo "  release         - Build for multiple platforms"
	@echo "  dev             - Build development version with debug info"
	@echo "  fmt             - Format code"
	@echo "  lint            - Lint code (requires golangci-lint)"
	@echo "  security        - Run security scan (requires gosec)"
	@echo "  help            - Show this help"