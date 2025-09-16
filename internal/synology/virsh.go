package synology

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// parseVirshList parses the output of 'virsh list --all'
func parseVirshList(output string) ([]VM, error) {
	var vms []VM

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 3 {
		// No VMs found (just header lines)
		return vms, nil
	}

	// Skip header lines (usually first 2 lines)
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// Parse line format: " Id   Name   State"
		// Example: " 1    test-vm    running"
		// Example: " -    stopped-vm shut off"
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		vm := VM{
			Name:   fields[1],
			Status: strings.Join(fields[2:], " "), // Handle multi-word states like "shut off"
		}

		vms = append(vms, vm)
	}

	return vms, nil
}

// getVMInfo gets detailed information about a specific VM using virsh
func (c *Client) getVMInfo(vmName string) (*VM, error) {
	// Get basic domain info
	output, err := c.ExecuteCommand(fmt.Sprintf("/usr/local/bin/virsh dominfo %s", vmName))
	if err != nil {
		return nil, fmt.Errorf("failed to get VM info: %w", err)
	}

	vm := &VM{Name: vmName}

	// Parse dominfo output
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "State":
			vm.Status = value
		case "CPU(s)":
			if cpu, err := strconv.Atoi(value); err == nil {
				vm.CPU = cpu
			}
		case "Max memory":
			// Parse memory like "2097152 KiB"
			if memStr := strings.Fields(value); len(memStr) > 0 {
				if mem, err := strconv.ParseFloat(memStr[0], 64); err == nil {
					vm.Memory = int(mem / 1024) // Convert KiB to MB
				}
			}
		}
	}

	// Try to get IP address
	if ip, err := c.getVMIPAddress(vmName); err == nil {
		vm.IPAddress = ip
	}

	return vm, nil
}

// getVMIPAddress attempts to get the IP address of a VM
func (c *Client) getVMIPAddress(vmName string) (string, error) {
	// Try to get IP from domifaddr
	output, err := c.ExecuteCommand(fmt.Sprintf("/usr/local/bin/virsh domifaddr %s", vmName))
	if err != nil {
		return "", err
	}

	// Parse output to extract IP
	lines := strings.Split(output, "\n")
	ipRegex := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)`)

	for _, line := range lines {
		if matches := ipRegex.FindStringSubmatch(line); len(matches) > 1 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("no IP address found")
}

// executeVirshCommand executes a virsh command
func (c *Client) executeVirshCommand(args string) error {
	cmd := fmt.Sprintf("/usr/local/bin/virsh %s", args)
	_, err := c.ExecuteCommand(cmd)
	return err
}