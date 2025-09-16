#!/bin/bash

# syno-vm installation script
# This script downloads and installs the latest version of syno-vm

set -e

# Configuration
REPO="scttfrdmn/syno-vm"
BINARY_NAME="syno-vm"
INSTALL_DIR="/usr/local/bin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)

    case $os in
        darwin)
            OS="darwin"
            ;;
        linux)
            OS="linux"
            ;;
        *)
            log_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac

    case $arch in
        x86_64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            log_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac

    log_info "Detected platform: $OS/$ARCH"
}

# Get the latest release version
get_latest_version() {
    log_info "Fetching latest version..."

    if command -v curl >/dev/null 2>&1; then
        VERSION=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        VERSION=$(wget -qO- https://api.github.com/repos/$REPO/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        log_error "curl or wget is required to download syno-vm"
        exit 1
    fi

    if [ -z "$VERSION" ]; then
        log_error "Failed to get latest version"
        exit 1
    fi

    log_info "Latest version: $VERSION"
}

# Download and install binary
install_binary() {
    local binary_name="$BINARY_NAME-$VERSION-$OS-$ARCH"
    local download_url="https://github.com/$REPO/releases/download/$VERSION/$binary_name"
    local temp_file="/tmp/$binary_name"

    log_info "Downloading $binary_name..."

    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$temp_file" "$download_url"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$temp_file" "$download_url"
    fi

    if [ ! -f "$temp_file" ]; then
        log_error "Failed to download binary"
        exit 1
    fi

    # Make binary executable
    chmod +x "$temp_file"

    # Install binary
    log_info "Installing to $INSTALL_DIR/$BINARY_NAME..."

    if [ -w "$INSTALL_DIR" ]; then
        mv "$temp_file" "$INSTALL_DIR/$BINARY_NAME"
    else
        sudo mv "$temp_file" "$INSTALL_DIR/$BINARY_NAME"
    fi

    log_info "Installation complete!"
}

# Create config directory
create_config_dir() {
    local config_dir="$HOME/.syno-vm"

    if [ ! -d "$config_dir" ]; then
        log_info "Creating config directory: $config_dir"
        mkdir -p "$config_dir"
    fi
}

# Verify installation
verify_installation() {
    if command -v $BINARY_NAME >/dev/null 2>&1; then
        local installed_version=$($BINARY_NAME --version 2>/dev/null | grep -o 'v[0-9]*\.[0-9]*\.[0-9]*' | head -1)
        log_info "syno-vm $installed_version installed successfully!"
        log_info "Run '$BINARY_NAME --help' to get started."
    else
        log_error "Installation verification failed. $BINARY_NAME not found in PATH."
        log_warn "You may need to add $INSTALL_DIR to your PATH or restart your shell."
        exit 1
    fi
}

# Show next steps
show_next_steps() {
    echo ""
    log_info "Next steps:"
    echo "  1. Configure your Synology NAS connection:"
    echo "     $BINARY_NAME config set --host your-synology.local --username admin"
    echo ""
    echo "  2. Test the connection:"
    echo "     $BINARY_NAME list"
    echo ""
    echo "  3. For more information, visit:"
    echo "     https://github.com/$REPO"
}

# Main installation process
main() {
    log_info "Installing syno-vm..."

    # Check if running as root (not recommended)
    if [ "$EUID" -eq 0 ]; then
        log_warn "Running as root is not recommended"
    fi

    # Detect platform
    detect_platform

    # Get latest version
    get_latest_version

    # Install binary
    install_binary

    # Create config directory
    create_config_dir

    # Verify installation
    verify_installation

    # Show next steps
    show_next_steps
}

# Check for help flag
if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo "syno-vm installation script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "This script downloads and installs the latest version of syno-vm."
    echo ""
    echo "Options:"
    echo "  -h, --help    Show this help message"
    echo ""
    echo "Environment variables:"
    echo "  INSTALL_DIR   Installation directory (default: /usr/local/bin)"
    echo ""
    exit 0
fi

# Run main installation
main