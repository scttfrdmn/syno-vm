package mock

import (
	"fmt"

	"github.com/scttfrdmn/syno-vm/internal/synology"
)

// MockClient is a mock implementation of the Synology client for testing
type MockClient struct {
	VMs       []synology.VM
	Templates []synology.Template
	Connected bool
	Fail      map[string]bool // Map of method names that should fail
}

// NewMockClient creates a new mock client with sample data
func NewMockClient() *MockClient {
	return &MockClient{
		VMs: []synology.VM{
			{
				Name:      "test-vm-1",
				Status:    "running",
				CPU:       2,
				Memory:    2048,
				Storage:   "20GB",
				IPAddress: "192.168.1.100",
			},
			{
				Name:      "test-vm-2",
				Status:    "stopped",
				CPU:       4,
				Memory:    4096,
				Storage:   "40GB",
				IPAddress: "",
			},
		},
		Templates: []synology.Template{
			{
				Name:        "ubuntu-20.04",
				Description: "Ubuntu 20.04 LTS",
				OS:          "Linux",
			},
			{
				Name:        "windows-10",
				Description: "Windows 10 Pro",
				OS:          "Windows",
			},
		},
		Connected: true,
		Fail:      make(map[string]bool),
	}
}

// Connect simulates connecting to the Synology NAS
func (m *MockClient) Connect() error {
	if m.Fail["Connect"] {
		return fmt.Errorf("mock connection failed")
	}
	m.Connected = true
	return nil
}

// Disconnect simulates disconnecting from the Synology NAS
func (m *MockClient) Disconnect() error {
	if m.Fail["Disconnect"] {
		return fmt.Errorf("mock disconnect failed")
	}
	m.Connected = false
	return nil
}

// ListVMs returns the mock VM list
func (m *MockClient) ListVMs() ([]synology.VM, error) {
	if m.Fail["ListVMs"] {
		return nil, fmt.Errorf("mock ListVMs failed")
	}
	return m.VMs, nil
}

// StartVM simulates starting a VM
func (m *MockClient) StartVM(vmName string) error {
	if m.Fail["StartVM"] {
		return fmt.Errorf("mock StartVM failed")
	}

	for i, vm := range m.VMs {
		if vm.Name == vmName {
			m.VMs[i].Status = "running"
			return nil
		}
	}

	return fmt.Errorf("VM not found: %s", vmName)
}

// StopVM simulates stopping a VM
func (m *MockClient) StopVM(vmName string) error {
	if m.Fail["StopVM"] {
		return fmt.Errorf("mock StopVM failed")
	}

	for i, vm := range m.VMs {
		if vm.Name == vmName {
			m.VMs[i].Status = "stopped"
			m.VMs[i].IPAddress = ""
			return nil
		}
	}

	return fmt.Errorf("VM not found: %s", vmName)
}

// RestartVM simulates restarting a VM
func (m *MockClient) RestartVM(vmName string) error {
	if m.Fail["RestartVM"] {
		return fmt.Errorf("mock RestartVM failed")
	}

	for i, vm := range m.VMs {
		if vm.Name == vmName {
			m.VMs[i].Status = "running"
			return nil
		}
	}

	return fmt.Errorf("VM not found: %s", vmName)
}

// GetVMStatus returns the status of a specific VM
func (m *MockClient) GetVMStatus(vmName string) (*synology.VM, error) {
	if m.Fail["GetVMStatus"] {
		return nil, fmt.Errorf("mock GetVMStatus failed")
	}

	for _, vm := range m.VMs {
		if vm.Name == vmName {
			return &vm, nil
		}
	}

	return nil, fmt.Errorf("VM not found: %s", vmName)
}

// CreateVM simulates creating a new VM
func (m *MockClient) CreateVM(config synology.VMConfig) error {
	if m.Fail["CreateVM"] {
		return fmt.Errorf("mock CreateVM failed")
	}

	// Check if VM already exists
	for _, vm := range m.VMs {
		if vm.Name == config.Name {
			return fmt.Errorf("VM already exists: %s", config.Name)
		}
	}

	// Add new VM to mock list
	newVM := synology.VM{
		Name:    config.Name,
		Status:  "stopped",
		CPU:     config.CPU,
		Memory:  config.Memory,
		Storage: config.Storage,
	}

	m.VMs = append(m.VMs, newVM)
	return nil
}

// DeleteVM simulates deleting a VM
func (m *MockClient) DeleteVM(vmName string) error {
	if m.Fail["DeleteVM"] {
		return fmt.Errorf("mock DeleteVM failed")
	}

	for i, vm := range m.VMs {
		if vm.Name == vmName {
			// Remove VM from slice
			m.VMs = append(m.VMs[:i], m.VMs[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("VM not found: %s", vmName)
}

// ListTemplates returns the mock template list
func (m *MockClient) ListTemplates() ([]synology.Template, error) {
	if m.Fail["ListTemplates"] {
		return nil, fmt.Errorf("mock ListTemplates failed")
	}
	return m.Templates, nil
}

// CreateTemplate simulates creating a new template
func (m *MockClient) CreateTemplate(templateName, vmName string) error {
	if m.Fail["CreateTemplate"] {
		return fmt.Errorf("mock CreateTemplate failed")
	}

	// Check if VM exists
	vmExists := false
	for _, vm := range m.VMs {
		if vm.Name == vmName {
			vmExists = true
			break
		}
	}

	if !vmExists {
		return fmt.Errorf("source VM not found: %s", vmName)
	}

	// Check if template already exists
	for _, template := range m.Templates {
		if template.Name == templateName {
			return fmt.Errorf("template already exists: %s", templateName)
		}
	}

	// Add new template
	newTemplate := synology.Template{
		Name:        templateName,
		Description: fmt.Sprintf("Template created from %s", vmName),
		OS:          "Linux", // Default for mock
	}

	m.Templates = append(m.Templates, newTemplate)
	return nil
}

// DeleteTemplate simulates deleting a template
func (m *MockClient) DeleteTemplate(templateName string) error {
	if m.Fail["DeleteTemplate"] {
		return fmt.Errorf("mock DeleteTemplate failed")
	}

	for i, template := range m.Templates {
		if template.Name == templateName {
			// Remove template from slice
			m.Templates = append(m.Templates[:i], m.Templates[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("template not found: %s", templateName)
}

// SetFailure configures the mock to fail specific method calls
func (m *MockClient) SetFailure(method string, shouldFail bool) {
	m.Fail[method] = shouldFail
}

// ResetFailures clears all failure configurations
func (m *MockClient) ResetFailures() {
	m.Fail = make(map[string]bool)
}