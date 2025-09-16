package synology

import (
	"bytes"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

// Client represents a Synology VMM client
type Client struct {
	host       string
	username   string
	port       int
	keyfile    string
	timeout    time.Duration
	sshClient  *ssh.Client
}

// VM represents a virtual machine
type VM struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	CPU       int    `json:"cpu"`
	Memory    int    `json:"memory"`
	Storage   string `json:"storage"`
	IPAddress string `json:"ip_address,omitempty"`
}

// VMConfig represents VM configuration for creation
type VMConfig struct {
	Name     string
	Template string
	CPU      int
	Memory   int
	Storage  string
}

// Validate validates the VM configuration
func (c VMConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("VM name is required")
	}
	if c.CPU <= 0 {
		return fmt.Errorf("CPU must be greater than 0")
	}
	if c.Memory <= 0 {
		return fmt.Errorf("memory must be greater than 0")
	}
	return nil
}

// Template represents a VM template
type Template struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	OS          string `json:"os"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   *APIError              `json:"error,omitempty"`
}

// APIError represents an API error
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewClient creates a new Synology VMM client
func NewClient() (*Client, error) {
	host := viper.GetString("host")
	username := viper.GetString("username")
	port := viper.GetInt("port")
	keyfile := viper.GetString("keyfile")
	timeout := viper.GetInt("timeout")

	if host == "" {
		return nil, fmt.Errorf("host not configured. Run 'syno-vm config set --host <hostname>'")
	}

	if username == "" {
		return nil, fmt.Errorf("username not configured. Run 'syno-vm config set --username <username>'")
	}

	client := &Client{
		host:     host,
		username: username,
		port:     port,
		keyfile:  keyfile,
		timeout:  time.Duration(timeout) * time.Second,
	}

	return client, nil
}

// Connect establishes an SSH connection to the Synology NAS
func (c *Client) Connect() error {
	if c.sshClient != nil {
		return nil // Already connected
	}

	config := &ssh.ClientConfig{
		User:            c.username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         c.timeout,
	}

	// Configure authentication methods
	var authMethods []ssh.AuthMethod

	// Try ssh-agent first
	if sshAgent, err := getSSHAgent(); err == nil {
		authMethods = append(authMethods, ssh.PublicKeysCallback(sshAgent.Signers))
	}

	// If keyfile is specified, add it as well
	if c.keyfile != "" {
		if key, err := readPrivateKey(c.keyfile); err == nil {
			authMethods = append(authMethods, ssh.PublicKeys(key))
		}
	}

	if len(authMethods) == 0 {
		return fmt.Errorf("no SSH authentication methods available. Please ensure ssh-agent is running or configure a keyfile")
	}

	config.Auth = authMethods

	address := fmt.Sprintf("%s:%d", c.host, c.port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return fmt.Errorf("failed to connect to SSH server: %w", err)
	}

	c.sshClient = client
	return nil
}

// Disconnect closes the SSH connection
func (c *Client) Disconnect() error {
	if c.sshClient != nil {
		return c.sshClient.Close()
	}
	return nil
}

// ExecuteCommand executes a command on the Synology NAS via SSH
func (c *Client) ExecuteCommand(command string) (string, error) {
	if err := c.Connect(); err != nil {
		return "", err
	}

	session, err := c.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer func() { _ = session.Close() }() // Ensure session cleanup

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	if err := session.Run(command); err != nil {
		return "", fmt.Errorf("command failed: %s, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}


// ListVMs lists all virtual machines using virsh
func (c *Client) ListVMs() ([]VM, error) {
	// Use virsh to list all VMs
	output, err := c.ExecuteCommand("/usr/local/bin/virsh list --all")
	if err != nil {
		return nil, fmt.Errorf("failed to list VMs: %w", err)
	}

	return parseVirshList(output)
}

// StartVM starts a virtual machine using virsh
func (c *Client) StartVM(vmName string) error {
	return c.executeVirshCommand(fmt.Sprintf("start %s", vmName))
}

// StopVM stops a virtual machine using virsh
func (c *Client) StopVM(vmName string) error {
	return c.executeVirshCommand(fmt.Sprintf("shutdown %s", vmName))
}

// RestartVM restarts a virtual machine using virsh
func (c *Client) RestartVM(vmName string) error {
	return c.executeVirshCommand(fmt.Sprintf("reboot %s", vmName))
}

// GetVMStatus gets the status of a specific virtual machine using virsh
func (c *Client) GetVMStatus(vmName string) (*VM, error) {
	return c.getVMInfo(vmName)
}

// CreateVM creates a new virtual machine using virsh
func (c *Client) CreateVM(config VMConfig) error {
	// VM creation via virsh requires XML configuration
	// This is a complex operation that would need a proper XML template
	// For now, return an error indicating this needs VMM GUI or more complex implementation
	return fmt.Errorf("VM creation via virsh requires XML configuration - please use VMM interface for VM creation")
}

// DeleteVM deletes a virtual machine using virsh
func (c *Client) DeleteVM(vmName string) error {
	// First undefine the domain (this removes it completely)
	return c.executeVirshCommand(fmt.Sprintf("undefine %s", vmName))
}

// ListTemplates lists available VM templates
func (c *Client) ListTemplates() ([]Template, error) {
	// VMM templates are typically stored as VM snapshots or images
	// For now, we'll return an empty list since template management
	// requires more complex VMM-specific operations
	var templates []Template

	// Note: Template management in VMM typically involves:
	// - Creating VMs with specific configurations
	// - Converting VMs to templates via the VMM interface
	// This functionality would need VMM-specific APIs or file system access

	return templates, nil
}

// CreateTemplate creates a new VM template
func (c *Client) CreateTemplate(templateName, vmName string) error {
	// Template creation in VMM typically requires the VMM interface
	// This would involve creating snapshots or exporting VM configurations
	return fmt.Errorf("template creation requires VMM interface - not implemented via virsh")
}

// DeleteTemplate deletes a VM template
func (c *Client) DeleteTemplate(templateName string) error {
	// Template deletion in VMM typically requires the VMM interface
	return fmt.Errorf("template deletion requires VMM interface - not implemented via virsh")
}