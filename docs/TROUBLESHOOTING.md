# Troubleshooting Guide

This guide helps resolve common issues when using syno-vm.

## Connection Issues

### SSH Connection Refused

**Error**: `connection refused` or `network is unreachable`

**Causes**:
- SSH service not enabled on Synology NAS
- Incorrect hostname/IP address
- Network connectivity issues
- Firewall blocking SSH

**Solutions**:
1. Enable SSH on your NAS:
   - DSM → Control Panel → Terminal & SNMP → Enable SSH service
2. Verify hostname/IP:
   ```bash
   ping your-synology.local
   # or
   ping 192.168.1.100
   ```
3. Test SSH manually:
   ```bash
   ssh admin@your-synology.local
   ```
4. Check firewall settings in DSM

### SSH Authentication Failed

**Error**: `permission denied (publickey,password)`

**Causes**:
- Incorrect username/password
- SSH key not properly configured
- SSH key permissions incorrect

**Solutions**:
1. Verify username (typically "admin"):
   ```bash
   syno-vm config set --username admin
   ```
2. Check SSH key permissions:
   ```bash
   chmod 600 ~/.ssh/id_rsa
   chmod 644 ~/.ssh/id_rsa.pub
   ```
3. Ensure public key is on NAS:
   ```bash
   ssh-copy-id admin@your-synology.local
   ```
4. Test with password authentication first

### SSH Timeout

**Error**: `connection timeout`

**Causes**:
- Network latency
- NAS overloaded
- Incorrect timeout setting

**Solutions**:
1. Increase timeout:
   ```bash
   syno-vm config set --timeout 60
   ```
2. Check NAS system resources in DSM
3. Try during less busy times

## API Issues

### synowebapi Command Not Found

**Error**: `command not found: synowebapi`

**Causes**:
- Virtual Machine Manager not installed
- Outdated DSM version
- Package not started

**Solutions**:
1. Install VMM package:
   - DSM → Package Center → Install "Virtual Machine Manager"
2. Update DSM to latest version
3. Restart VMM package:
   - Package Center → Virtual Machine Manager → Stop/Start

### API Permission Denied

**Error**: `insufficient user privilege`

**Causes**:
- User not in administrators group
- VMM permissions not granted

**Solutions**:
1. Ensure user is administrator:
   - DSM → Control Panel → User & Group → User → Edit user → User Groups → administrators
2. Grant VMM permissions:
   - DSM → Control Panel → User & Group → User → Edit user → Applications → Virtual Machine Manager

### Invalid API Response

**Error**: `failed to parse API response`

**Causes**:
- VMM service not running
- Corrupted API response
- Version mismatch

**Solutions**:
1. Restart VMM service:
   - Package Center → Virtual Machine Manager → Stop/Start
2. Check DSM logs for errors
3. Update VMM package to latest version

## VM Management Issues

### VM Creation Failed

**Error**: `failed to create VM`

**Causes**:
- Insufficient resources
- No Btrfs volume available
- VMM not properly configured

**Solutions**:
1. Check available resources:
   - DSM → Resource Monitor
2. Ensure Btrfs volume exists:
   - DSM → Storage Manager → Volume
3. Configure vSwitch:
   - Control Panel → Network → Network Interface → Manage → Open vSwitch Settings

### VM Won't Start

**Error**: `failed to start VM`

**Causes**:
- Insufficient RAM
- CPU doesn't support virtualization
- VM configuration issues

**Solutions**:
1. Check CPU virtualization support:
   - Intel VT-x or AMD-V required
2. Verify RAM availability:
   - DSM → Resource Monitor
3. Check VM configuration in VMM interface
4. Review VMM logs

### VM Slow Performance

**Causes**:
- Insufficient allocated resources
- Host system overloaded
- Storage performance issues

**Solutions**:
1. Increase VM resources:
   ```bash
   # Edit VM in VMM interface to increase CPU/RAM
   ```
2. Monitor host performance:
   - DSM → Resource Monitor
3. Use SSD for VM storage if available
4. Reduce concurrent VM count

## Configuration Issues

### Config File Not Found

**Error**: `config file not found`

**Solutions**:
1. Create config directory:
   ```bash
   mkdir -p ~/.syno-vm
   ```
2. Set initial configuration:
   ```bash
   syno-vm config set --host your-synology.local --username admin
   ```

### Invalid Configuration Values

**Error**: `invalid configuration`

**Solutions**:
1. Check config file format:
   ```bash
   cat ~/.syno-vm/config.yaml
   ```
2. Reset configuration:
   ```bash
   rm ~/.syno-vm/config.yaml
   syno-vm config set --host your-nas.local --username admin
   ```

## Network Issues

### VM No Network Access

**Causes**:
- vSwitch not configured
- Network adapter not connected
- DHCP issues

**Solutions**:
1. Enable vSwitch:
   - Control Panel → Network → Network Interface → Manage → Open vSwitch Settings
2. Check VM network settings in VMM
3. Verify DHCP server configuration
4. Try static IP configuration

### Cannot Connect to VM

**Causes**:
- VM firewall blocking connections
- Network configuration issues
- VM not fully booted

**Solutions**:
1. Check VM console in VMM interface
2. Wait for VM to fully boot
3. Verify VM network configuration
4. Check VM firewall settings

## Hardware Issues

### Virtualization Not Supported

**Error**: `hardware virtualization not available`

**Causes**:
- CPU doesn't support VT-x/AMD-V
- Virtualization disabled in BIOS
- Running on unsupported hardware

**Solutions**:
1. Check CPU specifications for virtualization support
2. Enable virtualization in BIOS/UEFI
3. Use compatible Synology model

### Insufficient Memory

**Error**: `not enough memory`

**Solutions**:
1. Increase NAS RAM
2. Reduce VM memory allocation
3. Close unnecessary applications on NAS
4. Stop other VMs

## Debugging Tips

### Enable Verbose Logging

```bash
syno-vm -v command
```

### Check SSH Connection

```bash
ssh -v admin@your-synology.local
```

### Monitor NAS Resources

```bash
# On NAS via SSH
top
free -m
df -h
```

### Check VMM Logs

1. SSH to NAS
2. Check logs in `/var/log/` directory
3. Look for VMM-related errors

## Getting Help

If you're still experiencing issues:

1. Check the [GitHub Issues](https://github.com/scttfrdmn/syno-vm/issues)
2. Create a new issue with:
   - syno-vm version (`syno-vm --version`)
   - DSM version
   - VMM package version
   - Complete error messages
   - Steps to reproduce

## Common Error Codes

| Code | Description | Solution |
|------|-------------|----------|
| 101 | Invalid parameter | Check command syntax |
| 102 | Unknown API | Update VMM package |
| 103 | Unknown method | Check API version |
| 104 | Not supported | Feature not available |
| 105 | Timeout | Increase timeout setting |
| 106 | Insufficient privileges | Check user permissions |