package integration

import (
	"os"
	"testing"

	"github.com/scttfrdmn/syno-vm/internal/synology"
	"github.com/spf13/viper"
)

// TestIntegration runs integration tests against a real Synology NAS
// These tests require environment variables to be set:
// - SYNO_VM_TEST_HOST: Synology NAS hostname/IP
// - SYNO_VM_TEST_USERNAME: SSH username
// - SYNO_VM_TEST_KEYFILE: SSH private key file path (optional)
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// Check for required environment variables
	host := os.Getenv("SYNO_VM_TEST_HOST")
	username := os.Getenv("SYNO_VM_TEST_USERNAME")

	if host == "" || username == "" {
		t.Skip("Integration tests require SYNO_VM_TEST_HOST and SYNO_VM_TEST_USERNAME environment variables")
	}

	// Set up viper configuration for testing
	viper.Reset()
	viper.Set("host", host)
	viper.Set("username", username)
	viper.Set("port", 22)
	viper.Set("timeout", 30)

	if keyfile := os.Getenv("SYNO_VM_TEST_KEYFILE"); keyfile != "" {
		viper.Set("keyfile", keyfile)
	}

	t.Run("CreateClient", testCreateClient)
	t.Run("ListVMs", testListVMs)
	t.Run("VMLifecycle", testVMLifecycle)
	t.Run("Templates", testTemplates)
}

func testCreateClient(t *testing.T) {
	client, err := synology.NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if client == nil {
		t.Fatal("Client is nil")
	}
}

func testListVMs(t *testing.T) {
	client, err := synology.NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	vms, err := client.ListVMs()
	if err != nil {
		t.Fatalf("Failed to list VMs: %v", err)
	}

	// VMs list can be empty, that's okay
	t.Logf("Found %d VMs", len(vms))

	for _, vm := range vms {
		if vm.Name == "" {
			t.Error("VM has empty name")
		}
		t.Logf("VM: %s, Status: %s, CPU: %d, Memory: %d", vm.Name, vm.Status, vm.CPU, vm.Memory)
	}
}

func testVMLifecycle(t *testing.T) {
	if os.Getenv("SYNO_VM_ENABLE_LIFECYCLE_TEST") != "true" {
		t.Skip("VM lifecycle tests require SYNO_VM_ENABLE_LIFECYCLE_TEST=true")
	}

	client, err := synology.NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	testVMName := "test-integration-vm"

	// Clean up any existing test VM
	defer func() {
		_ = client.DeleteVM(testVMName) // Best effort cleanup
	}()

	t.Run("CreateVM", func(t *testing.T) {
		config := synology.VMConfig{
			Name:   testVMName,
			CPU:    1,
			Memory: 1024,
		}

		err := client.CreateVM(config)
		if err != nil {
			t.Fatalf("Failed to create VM: %v", err)
		}
	})

	t.Run("GetVMStatus", func(t *testing.T) {
		vm, err := client.GetVMStatus(testVMName)
		if err != nil {
			t.Fatalf("Failed to get VM status: %v", err)
		}

		if vm.Name != testVMName {
			t.Errorf("Expected VM name %s, got %s", testVMName, vm.Name)
		}
	})

	t.Run("StartVM", func(t *testing.T) {
		err := client.StartVM(testVMName)
		if err != nil {
			t.Fatalf("Failed to start VM: %v", err)
		}
	})

	t.Run("StopVM", func(t *testing.T) {
		err := client.StopVM(testVMName)
		if err != nil {
			t.Fatalf("Failed to stop VM: %v", err)
		}
	})

	t.Run("DeleteVM", func(t *testing.T) {
		err := client.DeleteVM(testVMName)
		if err != nil {
			t.Fatalf("Failed to delete VM: %v", err)
		}
	})
}

func testTemplates(t *testing.T) {
	client, err := synology.NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	templates, err := client.ListTemplates()
	if err != nil {
		t.Fatalf("Failed to list templates: %v", err)
	}

	t.Logf("Found %d templates", len(templates))

	for _, template := range templates {
		if template.Name == "" {
			t.Error("Template has empty name")
		}
		t.Logf("Template: %s, Description: %s, OS: %s", template.Name, template.Description, template.OS)
	}
}

// BenchmarkListVMs benchmarks the ListVMs operation
func BenchmarkListVMs(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmarks in short mode")
	}

	host := os.Getenv("SYNO_VM_TEST_HOST")
	username := os.Getenv("SYNO_VM_TEST_USERNAME")

	if host == "" || username == "" {
		b.Skip("Benchmarks require SYNO_VM_TEST_HOST and SYNO_VM_TEST_USERNAME environment variables")
	}

	viper.Reset()
	viper.Set("host", host)
	viper.Set("username", username)
	viper.Set("port", 22)
	viper.Set("timeout", 30)

	if keyfile := os.Getenv("SYNO_VM_TEST_KEYFILE"); keyfile != "" {
		viper.Set("keyfile", keyfile)
	}

	client, err := synology.NewClient()
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := client.ListVMs()
		if err != nil {
			b.Fatalf("Failed to list VMs: %v", err)
		}
	}
}