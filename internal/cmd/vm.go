package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/scttfrdmn/syno-vm/internal/synology"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start <vm-name>",
	Short: "Start a virtual machine",
	Long:  `Start a virtual machine by name.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runStart,
}

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop <vm-name>",
	Short: "Stop a virtual machine",
	Long:  `Stop a virtual machine by name.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runStop,
}

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart <vm-name>",
	Short: "Restart a virtual machine",
	Long:  `Restart a virtual machine by name.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runRestart,
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status <vm-name>",
	Short: "Show virtual machine status",
	Long:  `Show detailed status information for a virtual machine.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runStatus,
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <vm-name>",
	Short: "Delete a virtual machine",
	Long:  `Delete a virtual machine by name. This action is irreversible.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDelete,
}

var (
	force bool
)

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVarP(&force, "force", "f", false, "Force delete without confirmation")
}

func runStart(cmd *cobra.Command, args []string) error {
	vmName := args[0]

	client, err := synology.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	fmt.Printf("Starting VM: %s\n", vmName)

	if err := client.StartVM(vmName); err != nil {
		return fmt.Errorf("failed to start VM: %w", err)
	}

	fmt.Printf("VM %s started successfully\n", vmName)
	return nil
}

func runStop(cmd *cobra.Command, args []string) error {
	vmName := args[0]

	client, err := synology.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	fmt.Printf("Stopping VM: %s\n", vmName)

	if err := client.StopVM(vmName); err != nil {
		return fmt.Errorf("failed to stop VM: %w", err)
	}

	fmt.Printf("VM %s stopped successfully\n", vmName)
	return nil
}

func runRestart(cmd *cobra.Command, args []string) error {
	vmName := args[0]

	client, err := synology.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	fmt.Printf("Restarting VM: %s\n", vmName)

	if err := client.RestartVM(vmName); err != nil {
		return fmt.Errorf("failed to restart VM: %w", err)
	}

	fmt.Printf("VM %s restarted successfully\n", vmName)
	return nil
}

func runStatus(cmd *cobra.Command, args []string) error {
	vmName := args[0]

	client, err := synology.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	vm, err := client.GetVMStatus(vmName)
	if err != nil {
		return fmt.Errorf("failed to get VM status: %w", err)
	}

	fmt.Printf("Virtual Machine: %s\n", vm.Name)
	fmt.Printf("Status: %s\n", vm.Status)
	fmt.Printf("CPU Cores: %d\n", vm.CPU)
	fmt.Printf("Memory: %d MB\n", vm.Memory)
	fmt.Printf("Storage: %s\n", vm.Storage)
	if vm.IPAddress != "" {
		fmt.Printf("IP Address: %s\n", vm.IPAddress)
	}

	return nil
}

func runDelete(cmd *cobra.Command, args []string) error {
	vmName := args[0]

	if !force {
		fmt.Printf("Are you sure you want to delete VM '%s'? This action cannot be undone. (y/N): ", vmName)
		var response string
		_, _ = fmt.Scanln(&response) // Ignore input errors for confirmation
		if response != "y" && response != "Y" && response != "yes" {
			fmt.Println("Delete cancelled.")
			return nil
		}
	}

	client, err := synology.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	fmt.Printf("Deleting VM: %s\n", vmName)

	if err := client.DeleteVM(vmName); err != nil {
		return fmt.Errorf("failed to delete VM: %w", err)
	}

	fmt.Printf("VM %s deleted successfully\n", vmName)
	return nil
}