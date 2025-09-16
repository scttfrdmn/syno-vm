package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/scttfrdmn/syno-vm/internal/synology"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new virtual machine",
	Long:  `Create a new virtual machine with specified configuration.`,
	RunE:  runCreate,
}

var (
	createName     string
	createTemplate string
	createCPU      int
	createMemory   int
	createStorage  string
)

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&createName, "name", "", "Name of the virtual machine (required)")
	createCmd.Flags().StringVar(&createTemplate, "template", "", "Template to use for VM creation")
	createCmd.Flags().IntVar(&createCPU, "cpu", 2, "Number of CPU cores")
	createCmd.Flags().IntVar(&createMemory, "memory", 2048, "Memory in MB")
	createCmd.Flags().StringVar(&createStorage, "storage", "", "Storage configuration")

	createCmd.MarkFlagRequired("name") // nolint:errcheck // CLI flag setup
}

func runCreate(cmd *cobra.Command, args []string) error {
	if createName == "" {
		return fmt.Errorf("VM name is required")
	}

	client, err := synology.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	vmConfig := synology.VMConfig{
		Name:     createName,
		Template: createTemplate,
		CPU:      createCPU,
		Memory:   createMemory,
		Storage:  createStorage,
	}

	fmt.Printf("Creating VM: %s\n", createName)
	fmt.Printf("  CPU: %d cores\n", createCPU)
	fmt.Printf("  Memory: %d MB\n", createMemory)
	if createTemplate != "" {
		fmt.Printf("  Template: %s\n", createTemplate)
	}

	if err := client.CreateVM(vmConfig); err != nil {
		return fmt.Errorf("failed to create VM: %w", err)
	}

	fmt.Printf("VM %s created successfully\n", createName)
	return nil
}