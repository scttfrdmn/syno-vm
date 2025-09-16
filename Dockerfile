# Build stage
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates for fetching dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o syno-vm \
    cmd/syno-vm/main.go

# Final stage
FROM alpine:3.18

# Install ca-certificates and ssh client
RUN apk --no-cache add ca-certificates openssh-client

# Create non-root user
RUN addgroup -g 1001 -S synovm && \
    adduser -u 1001 -S synovm -G synovm

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/syno-vm .

# Copy configuration examples
COPY --from=builder /build/examples/ ./examples/

# Change ownership to non-root user
RUN chown -R synovm:synovm /app

# Switch to non-root user
USER synovm

# Create config directory
RUN mkdir -p /home/synovm/.syno-vm

# Expose any necessary ports (none for this CLI tool)

# Set entrypoint
ENTRYPOINT ["./syno-vm"]

# Default command
CMD ["--help"]

# Labels for metadata
LABEL org.opencontainers.image.title="syno-vm" \
      org.opencontainers.image.description="A CLI tool for managing Synology Virtual Machine Manager" \
      org.opencontainers.image.vendor="Scott Friedman" \
      org.opencontainers.image.source="https://github.com/scttfrdmn/syno-vm" \
      org.opencontainers.image.documentation="https://github.com/scttfrdmn/syno-vm/blob/main/README.md" \
      org.opencontainers.image.licenses="MIT"