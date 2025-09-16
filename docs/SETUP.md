# Setup Guide for syno-vm

This guide will help you set up syno-vm to manage virtual machines on your Synology NAS.

## Prerequisites

Before using syno-vm, ensure your Synology NAS meets these requirements:

### Hardware Requirements
- Synology NAS with Intel or AMD CPU supporting virtualization (Intel VT-x or AMD-V)
- At least 4GB RAM (8GB recommended)
- At least one Btrfs volume for VM storage

### Software Requirements
- DSM 7.x or later
- Virtual Machine Manager (VMM) package installed
- SSH access enabled

## Synology NAS Configuration

### 1. Enable SSH Access

1. Log in to DSM as an administrator
2. Go to **Control Panel** > **Terminal & SNMP**
3. Check **Enable SSH service**
4. Set the SSH port (default: 22)
5. Click **Apply**

### 2. Install Virtual Machine Manager

1. Open **Package Center**
2. Search for "Virtual Machine Manager"
3. Click **Install**
4. Wait for installation to complete

### 3. Configure Virtual Switch

VMM requires vSwitch to be enabled:

1. Go to **Control Panel** > **Network** > **Network Interface**
2. Click **Manage** > **Open vSwitch Settings**
3. Check **Enable Open vSwitch**
4. Click **Apply**

### 4. Create Btrfs Volume (if needed)

VMs require a Btrfs volume for storage:

1. Go to **Storage Manager** > **Volume**
2. If you don't have a Btrfs volume, create one:
   - Click **Create** > **Create Volume**
   - Choose **Btrfs** as the file system
   - Follow the wizard to complete setup

## SSH Key Setup (Recommended)

For secure, password-less authentication, set up SSH key authentication:

### 1. Generate SSH Key Pair (if you don't have one)

On your local machine:

```bash
ssh-keygen -t rsa -b 4096 -C "your_email@example.com"
```

### 2. Copy Public Key to Synology NAS

```bash
ssh-copy-id admin@your-synology.local
```

Or manually:

1. Copy the contents of `~/.ssh/id_rsa.pub`
2. SSH to your NAS: `ssh admin@your-synology.local`
3. Create the SSH directory: `mkdir -p ~/.ssh`
4. Add your key: `echo "your_public_key_here" >> ~/.ssh/authorized_keys`
5. Set permissions: `chmod 600 ~/.ssh/authorized_keys`

## syno-vm Installation

### Option 1: Download Pre-built Binary

1. Go to the [releases page](https://github.com/scttfrdmn/syno-vm/releases)
2. Download the appropriate binary for your platform
3. Make it executable: `chmod +x syno-vm`
4. Move to a directory in your PATH: `sudo mv syno-vm /usr/local/bin/`

### Option 2: Build from Source

```bash
git clone https://github.com/scttfrdmn/syno-vm
cd syno-vm
make build
sudo make install
```

### Option 3: Using Go

```bash
go install github.com/scttfrdmn/syno-vm/cmd/syno-vm@latest
```

## Configuration

### 1. Initial Configuration

Configure your Synology NAS connection:

```bash
syno-vm config set --host your-synology.local --username admin
```

If using SSH key authentication:

```bash
syno-vm config set --keyfile ~/.ssh/id_rsa
```

### 2. Test Connection

Verify your configuration:

```bash
syno-vm list
```

If successful, you should see a list of VMs (or a message indicating no VMs are found).

## Troubleshooting

### Connection Issues

**Problem**: "connection refused" or timeout errors

**Solutions**:
1. Verify SSH is enabled on your NAS
2. Check the SSH port (default: 22)
3. Ensure firewall allows SSH connections
4. Try connecting manually: `ssh admin@your-synology.local`

### Authentication Issues

**Problem**: "permission denied" errors

**Solutions**:
1. Verify username is correct (typically "admin")
2. Check SSH key permissions: `chmod 600 ~/.ssh/id_rsa`
3. Ensure public key is in `~/.ssh/authorized_keys` on the NAS
4. Try password authentication first to verify connectivity

### API Issues

**Problem**: "command not found: synowebapi"

**Solutions**:
1. Ensure VMM package is installed
2. Update to the latest DSM version
3. Check that VMM is running: log in to DSM and open VMM

### VMM Not Working

**Problem**: VMM interface shows errors or VMs don't start

**Solutions**:
1. Ensure CPU supports virtualization (Intel VT-x or AMD-V)
2. Verify sufficient RAM is available
3. Check that vSwitch is enabled
4. Ensure at least one Btrfs volume exists

## Next Steps

Once setup is complete, you can:

1. **List VMs**: `syno-vm list`
2. **Create a VM**: `syno-vm create --name test-vm --cpu 2 --memory 2048`
3. **Start a VM**: `syno-vm start test-vm`
4. **Check status**: `syno-vm status test-vm`

For more commands, run:

```bash
syno-vm --help
```