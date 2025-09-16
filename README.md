# syno-vm

A command-line tool for managing virtual machines on Synology NAS devices with Virtual Machine Manager (VMM).

This project is adapted from [qnap-vm](https://github.com/scttfrdmn/qnap-vm) for Synology DSM 7.x+ systems.

## Features

- **VM Lifecycle Management**: Create, start, stop, restart, and delete virtual machines
- **SSH-based Connection**: Secure remote management via SSH
- **Synology VMM API Integration**: Native integration with Synology's Virtual Machine Manager
- **VM Template Support**: Create VMs from predefined templates
- **Storage Management**: Manage VM storage and volumes
- **Cross-platform Support**: Works on macOS, Linux, and Windows

## Prerequisites

- Synology NAS running DSM 7.x or later
- Virtual Machine Manager (VMM) package installed
- SSH access enabled on your Synology NAS
- At least 4GB RAM on the NAS
- Intel VT-x or AMD-V CPU support
- At least one Btrfs volume for VM storage

## Installation

### From GitHub Releases

Download the latest release for your platform from the [releases page](https://github.com/scttfrdmn/syno-vm/releases).

### From Source

```bash
go install github.com/scttfrdmn/syno-vm/cmd/syno-vm@latest
```

### Homebrew (coming soon)

```bash
brew install scttfrdmn/syno-vm/syno-vm
```

## Quick Start

1. Configure your connection:
   ```bash
   syno-vm config set --host your-synology.local --username admin
   ```

2. List your virtual machines:
   ```bash
   syno-vm list
   ```

3. Create a new VM:
   ```bash
   syno-vm create --name my-vm --template ubuntu-20.04
   ```

4. Start the VM:
   ```bash
   syno-vm start my-vm
   ```

5. Check VM status:
   ```bash
   syno-vm status my-vm
   ```

## Configuration

The tool uses a YAML configuration file stored at `~/.syno-vm/config.yaml`:

```yaml
host: "your-synology.local"
username: "admin"
port: 22
keyfile: "~/.ssh/id_rsa"
timeout: 30
```

## Commands

### Configuration
- `syno-vm config set` - Set configuration values
- `syno-vm config get` - Get configuration values
- `syno-vm config list` - List all configuration

### VM Management
- `syno-vm list` - List all virtual machines
- `syno-vm create` - Create a new virtual machine
- `syno-vm start <vm-name>` - Start a virtual machine
- `syno-vm stop <vm-name>` - Stop a virtual machine
- `syno-vm restart <vm-name>` - Restart a virtual machine
- `syno-vm delete <vm-name>` - Delete a virtual machine
- `syno-vm status <vm-name>` - Show VM status

### Templates
- `syno-vm template list` - List available VM templates
- `syno-vm template create` - Create a new template
- `syno-vm template delete` - Delete a template

## API Integration

This tool integrates with Synology's VMM API using the `synowebapi` command-line tool available on DSM. Key API endpoints used:

- `SYNO.Virtualization.API.Guest.Action` - VM power operations
- `SYNO.Virtualization.API.Guest.Info` - VM information
- `SYNO.Virtualization.API.Host` - Host information

## Development

### Building from Source

```bash
git clone https://github.com/scttfrdmn/syno-vm
cd syno-vm
go build -o bin/syno-vm cmd/syno-vm/main.go
```

### Running Tests

```bash
go test ./...
```

## Differences from qnap-vm

- Uses Synology VMM API instead of QNAP Virtualization Station
- Adapted for DSM 7.x+ specific features
- Uses `synowebapi` for API communication instead of libvirt
- Different storage and networking configuration

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues and feature requests, please use the [GitHub Issues](https://github.com/scttfrdmn/syno-vm/issues) page.