package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/scttfrdmn/syno-vm/internal/synology"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List virtual machines",
	Long:  `List all virtual machines on the Synology NAS.`,
	RunE:  runList,
}

var (
	listAll bool
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Show all VMs including stopped ones")
}

func runList(cmd *cobra.Command, args []string) error {
	client, err := synology.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	vms, err := client.ListVMs()
	if err != nil {
		return fmt.Errorf("failed to list VMs: %w", err)
	}

	if len(vms) == 0 {
		fmt.Println("No virtual machines found.")
		return nil
	}

	// Print header
	fmt.Printf("%-20s %-15s %-10s %-15s\n", "NAME", "STATUS", "CPU", "MEMORY")
	fmt.Println("------------------------------------------------------------")

	// Print VMs
	for _, vm := range vms {
		if !listAll && vm.Status == "stopped" {
			continue
		}
		fmt.Printf("%-20s %-15s %-10s %-15s\n",
			vm.Name,
			vm.Status,
			fmt.Sprintf("%d cores", vm.CPU),
			fmt.Sprintf("%d MB", vm.Memory))
	}

	return nil
}