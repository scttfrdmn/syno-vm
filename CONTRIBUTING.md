# Contributing to syno-vm

Thank you for your interest in contributing to syno-vm! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contributing Process](#contributing-process)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)
- [Release Process](#release-process)

## Code of Conduct

This project adheres to a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Access to a Synology NAS with VMM for integration testing (optional)
- Docker (optional, for container builds)

### Development Tools

Install the following tools for the best development experience:

```bash
# golangci-lint for code linting
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# gosec for security scanning
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# goimports for import formatting
go install golang.org/x/tools/cmd/goimports@latest
```

## Development Setup

1. **Fork the repository** on GitHub

2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/syno-vm.git
   cd syno-vm
   ```

3. **Add the upstream remote**:
   ```bash
   git remote add upstream https://github.com/scttfrdmn/syno-vm.git
   ```

4. **Install dependencies**:
   ```bash
   make deps
   ```

5. **Build the project**:
   ```bash
   make build
   ```

6. **Run tests**:
   ```bash
   make test
   ```

## Contributing Process

### 1. Create an Issue

Before starting work, create or find an existing issue that describes:
- The problem you're solving
- The proposed solution
- Any design considerations

### 2. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### 3. Make Your Changes

- Write clear, concise code
- Follow the coding standards below
- Add tests for new functionality
- Update documentation as needed

### 4. Test Your Changes

```bash
# Run unit tests
make test

# Run linting
make lint

# Run security scan
make security

# Generate coverage report
make coverage
```

### 5. Commit Your Changes

Use conventional commit messages:

```bash
git commit -m "feat: add new VM template management"
git commit -m "fix: resolve SSH connection timeout issue"
git commit -m "docs: update API documentation"
```

Commit message types:
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test additions or changes
- `chore`: Maintenance tasks

### 6. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Create a pull request with:
- Clear title and description
- Reference to related issues
- Screenshots (if UI changes)
- Test results

## Coding Standards

### Go Code Style

- Follow standard Go formatting (`gofmt`, `goimports`)
- Use meaningful variable and function names
- Keep functions small and focused
- Add comments for exported functions and types
- Handle errors appropriately

### Project Structure

```
syno-vm/
├── cmd/                    # Application entry points
├── internal/              # Private application code
│   ├── cmd/              # CLI commands
│   └── synology/         # Synology client library
├── test/                  # Test utilities and integration tests
│   ├── integration/      # Integration tests
│   └── mock/            # Mock implementations
├── docs/                  # Documentation
├── examples/             # Example configurations
└── scripts/              # Build and utility scripts
```

### Error Handling

- Always handle errors appropriately
- Use meaningful error messages
- Wrap errors with context when needed

```go
if err != nil {
    return fmt.Errorf("failed to connect to Synology NAS: %w", err)
}
```

### Logging

- Use structured logging where appropriate
- Include relevant context in log messages
- Use appropriate log levels

## Testing

### Unit Tests

- Write tests for all new functionality
- Use table-driven tests where appropriate
- Mock external dependencies
- Aim for high test coverage (>80%)

```go
func TestVMConfig_Validate(t *testing.T) {
    tests := []struct {
        name      string
        config    VMConfig
        wantError bool
    }{
        // Test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation...
        })
    }
}
```

### Integration Tests

- Test against real Synology NAS when possible
- Use environment variables for configuration
- Make tests skippable if resources unavailable

```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration tests in short mode")
    }

    host := os.Getenv("SYNO_VM_TEST_HOST")
    if host == "" {
        t.Skip("Integration tests require SYNO_VM_TEST_HOST")
    }

    // Test implementation...
}
```

### Running Tests

```bash
# Unit tests only
make test

# Integration tests (requires Synology NAS)
SYNO_VM_TEST_HOST=your-nas.local \
SYNO_VM_TEST_USERNAME=admin \
SYNO_VM_TEST_KEYFILE=~/.ssh/id_rsa \
make integration-test

# All tests with coverage
make coverage
```

## Documentation

### Code Documentation

- Document all exported functions and types
- Use Go doc conventions
- Provide examples where helpful

```go
// CreateVM creates a new virtual machine with the specified configuration.
// It returns an error if the VM already exists or if the configuration is invalid.
//
// Example:
//   config := VMConfig{Name: "test-vm", CPU: 2, Memory: 2048}
//   err := client.CreateVM(config)
func (c *Client) CreateVM(config VMConfig) error {
    // Implementation...
}
```

### User Documentation

- Update README.md for user-facing changes
- Add examples to docs/ directory
- Update API documentation as needed

### Commit Documentation

- Include documentation updates in the same commit as code changes
- Reference documentation in pull requests

## Release Process

### Versioning

We use [Semantic Versioning](https://semver.org/):
- `MAJOR.MINOR.PATCH`
- Major: Breaking changes
- Minor: New features (backward compatible)
- Patch: Bug fixes (backward compatible)

### Release Checklist

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create release tag: `git tag -a v1.2.3 -m "Release v1.2.3"`
4. Push tag: `git push origin v1.2.3`
5. GitHub Actions will automatically build and publish release

## Getting Help

- **Issues**: Use GitHub Issues for bugs and feature requests
- **Discussions**: Use GitHub Discussions for questions and ideas
- **Security**: Report security issues privately via email

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Special recognition for significant contributions

Thank you for contributing to syno-vm!