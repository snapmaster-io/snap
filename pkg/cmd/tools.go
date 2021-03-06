package cmd

import (
	"fmt"
	"os"

	"github.com/snapmaster-io/snap/pkg/api"
	"github.com/snapmaster-io/snap/pkg/print"
	"github.com/snapmaster-io/snap/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

// toolsCmd represents the tools command
var toolsCmd = &cobra.Command{
	Use:   "tools [subcommand]",
	Short: "Interact with SnapMaster tools",
	Long:  `Interact with SnapMaster tools.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

// getToolCmd represents the get tool subcommand
var getToolCmd = &cobra.Command{
	Use:   "get [tool]",
	Short: "Get a description of a tool",
	Long:  `Get a description of a tool.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve tool as the first argument
		tool := args[0]

		// execute the API call
		response, err := api.Get("/connections")
		if err != nil {
			utils.PrintErrorMessage("could not retrieve data", err)
			os.Exit(1)
		}

		// select the entry that matches the provider name
		toolDescription := gjson.GetBytes(response, fmt.Sprintf("data.#(provider==%s)|@pretty", tool))
		if toolDescription.Raw == "" {
			utils.PrintError(fmt.Sprintf("tool %s not found", tool))
			return
		}

		format, err := rootCmd.PersistentFlags().GetString("format")
		if format == "json" {
			print.JSONString(toolDescription.Raw)
			return
		}

		// select the entry that matches the provider name
		description := gjson.GetBytes(response, fmt.Sprintf("data.#(provider==%s).definition.text", tool)).Value()
		// print the tool description
		utils.PrintYAML(description.(string))
	},
}

// listToolsCmd represents the list tools subcommand
var listToolsCmd = &cobra.Command{
	Use:   "list",
	Short: "List tools in the SnapMaster tools library",
	Long:  `List the user's connections.`,
	Run: func(cmd *cobra.Command, args []string) {

		// execute the API call
		response, err := api.Get("/connections")
		if err != nil {
			utils.PrintErrorMessage("could not retrieve data", err)
			os.Exit(1)
		}

		format, err := rootCmd.PersistentFlags().GetString("format")
		if format == "json" {
			print.JSON(response)
			return
		}

		if format == "table" {
			print.ToolsTable(response)
			return
		}

		// unknown format - return the raw response
		print.RawResponse(response)
	},
}

func init() {
	rootCmd.AddCommand(toolsCmd)
	toolsCmd.AddCommand(getToolCmd)
	toolsCmd.AddCommand(listToolsCmd)
}
