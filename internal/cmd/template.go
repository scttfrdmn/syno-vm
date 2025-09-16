package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/scttfrdmn/syno-vm/internal/synology"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage VM templates",
	Long:  `Manage virtual machine templates.`,
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available VM templates",
	Long:  `List all available VM templates.`,
	RunE:  runTemplateList,
}

var templateCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new VM template",
	Long:  `Create a new VM template from an existing VM.`,
	RunE:  runTemplateCreate,
}

var templateDeleteCmd = &cobra.Command{
	Use:   "delete <template-name>",
	Short: "Delete a VM template",
	Long:  `Delete a VM template by name.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runTemplateDelete,
}

var (
	templateName   string
	templateFromVM string
)

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateCreateCmd)
	templateCmd.AddCommand(templateDeleteCmd)

	templateCreateCmd.Flags().StringVar(&templateName, "name", "", "Name of the template (required)")
	templateCreateCmd.Flags().StringVar(&templateFromVM, "from-vm", "", "Create template from existing VM (required)")
	templateCreateCmd.MarkFlagRequired("name")   // nolint:errcheck // CLI setup
	templateCreateCmd.MarkFlagRequired("from-vm") // nolint:errcheck // CLI setup
}

func runTemplateList(cmd *cobra.Command, args []string) error {
	client, err := synology.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	templates, err := client.ListTemplates()
	if err != nil {
		return fmt.Errorf("failed to list templates: %w", err)
	}

	if len(templates) == 0 {
		fmt.Println("No templates found.")
		return nil
	}

	fmt.Printf("%-20s %-20s %-15s\n", "NAME", "DESCRIPTION", "OS")
	fmt.Println("-------------------------------------------------------")

	for _, template := range templates {
		fmt.Printf("%-20s %-20s %-15s\n",
			template.Name,
			template.Description,
			template.OS)
	}

	return nil
}

func runTemplateCreate(cmd *cobra.Command, args []string) error {
	client, err := synology.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	fmt.Printf("Creating template '%s' from VM '%s'\n", templateName, templateFromVM)

	if err := client.CreateTemplate(templateName, templateFromVM); err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	fmt.Printf("Template %s created successfully\n", templateName)
	return nil
}

func runTemplateDelete(cmd *cobra.Command, args []string) error {
	templateName := args[0]

	client, err := synology.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	fmt.Printf("Deleting template: %s\n", templateName)

	if err := client.DeleteTemplate(templateName); err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	fmt.Printf("Template %s deleted successfully\n", templateName)
	return nil
}