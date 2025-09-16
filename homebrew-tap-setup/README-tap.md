# homebrew-syno-vm

Homebrew tap for syno-vm - A CLI tool for managing Synology Virtual Machine Manager

## Installation

```bash
brew tap scttfrdmn/syno-vm
brew install syno-vm
```

Or install directly:

```bash
brew install scttfrdmn/syno-vm/syno-vm
```

## Usage

After installation, you can use syno-vm:

```bash
# Configure connection to your Synology NAS
syno-vm config set --host your-synology.local --username admin

# List virtual machines
syno-vm list

# Start a VM
syno-vm start vm-name

# Stop a VM
syno-vm stop vm-name

# Get VM status
syno-vm status vm-name
```

## Documentation

For full documentation, visit the main project: https://github.com/scttfrdmn/syno-vm

## Support

- Issues: https://github.com/scttfrdmn/syno-vm/issues
- Discussions: https://github.com/scttfrdmn/syno-vm/discussions