# Synology VMM API Reference

This document outlines the Synology Virtual Machine Manager (VMM) API endpoints used by syno-vm.

## Overview

syno-vm uses the Synology Web API through the `synowebapi` command-line tool available on DSM. All API calls are made via SSH to the Synology NAS.

## Authentication

API calls require authentication through the `synowebapi` tool, which handles session management automatically when executed on the NAS.

## Base API Structure

```bash
synowebapi --exec api=<API_NAME> method=<METHOD> version=<VERSION> [parameters]
```

## Guest Management APIs

### SYNO.Virtualization.API.Guest.Action

Power management operations for virtual machines.

**Power On VM:**
```bash
synowebapi --exec api=SYNO.Virtualization.API.Guest.Action version=1 method=poweron runner=admin guest_name="vm-name"
```

**Power Off VM:**
```bash
synowebapi --exec api=SYNO.Virtualization.API.Guest.Action version=1 method=poweroff runner=admin guest_name="vm-name"
```

**Restart VM:**
```bash
synowebapi --exec api=SYNO.Virtualization.API.Guest.Action version=1 method=restart runner=admin guest_name="vm-name"
```

### SYNO.Virtualization.API.Guest.Info

Get information about virtual machines.

**List All VMs:**
```bash
synowebapi --exec api=SYNO.Virtualization.API.Guest.Info version=1 method=list runner=admin
```

**Get VM Details:**
```bash
synowebapi --exec api=SYNO.Virtualization.API.Guest.Info version=1 method=get runner=admin guest_name="vm-name"
```

### SYNO.Virtualization.API.Guest

VM lifecycle management.

**Create VM:**
```bash
synowebapi --exec api=SYNO.Virtualization.API.Guest version=1 method=create runner=admin guest_name="vm-name" cpu=2 memory=2048
```

**Delete VM:**
```bash
synowebapi --exec api=SYNO.Virtualization.API.Guest version=1 method=delete runner=admin guest_name="vm-name"
```

## Template Management APIs

### SYNO.Virtualization.API.Template

Manage VM templates.

**List Templates:**
```bash
synowebapi --exec api=SYNO.Virtualization.API.Template version=1 method=list runner=admin
```

**Create Template:**
```bash
synowebapi --exec api=SYNO.Virtualization.API.Template version=1 method=create runner=admin template_name="template-name" source_vm="vm-name"
```

**Delete Template:**
```bash
synowebapi --exec api=SYNO.Virtualization.API.Template version=1 method=delete runner=admin template_name="template-name"
```

## Host Information APIs

### SYNO.Virtualization.API.Host

Get host system information.

**Get Host Info:**
```bash
synowebapi --exec api=SYNO.Virtualization.API.Host version=1 method=get runner=admin
```

## Response Format

All API responses are returned in JSON format:

```json
{
  "success": true,
  "data": {
    // Response data
  },
  "error": {
    "code": 0,
    "message": "Success"
  }
}
```

### Error Responses

When an API call fails, the response will contain error information:

```json
{
  "success": false,
  "error": {
    "code": 101,
    "message": "Invalid parameter"
  }
}
```

## Common Parameters

- `runner`: The user executing the operation (typically "admin")
- `guest_name`: The name of the virtual machine
- `version`: API version (typically "1")

## Limitations

- Authentication is handled by the `synowebapi` tool
- All operations require administrator privileges
- API availability depends on VMM package installation and version
- Some advanced features may require specific DSM versions

## Notes

This API reference is based on observed behavior and may not be complete. For official documentation, refer to the Synology Virtual Machine Manager API Guide provided by Synology.