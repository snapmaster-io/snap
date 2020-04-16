package cmd

import (
	"fmt"
	"os"

	"github.com/snapmaster-io/snap/pkg/api"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get all logs for the current user",
	Long:  `Get all logs for the current user.`,
	Run: func(cmd *cobra.Command, args []string) {
		// execute the API call
		response, err := api.Get("/logs")
		if err != nil {
			fmt.Printf("snap: could not retrieve data: %s", err)
			os.Exit(1)
		}

		format, err := rootCmd.PersistentFlags().GetString("format")
		if format == "json" {
			printJSON(response)
			return
		}

		if format == "table" {
			printActiveSnapLogs(response)
			return
		}

		// unknown format - return the raw response
		fmt.Printf("Raw response:\n%s\n", string(response))
	},
}

// getActiveSnapLogsCmd represents the get active snap logs subcommand
var logDetailsCmd = &cobra.Command{
	Use:   "details [log ID]",
	Short: "Get the details of a log entry",
	Long:  `Get the details of a log entry.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve activeSnapID as the first argument
		logID := args[0]

		// execute the API call
		response, err := api.Get("/logs")
		if err != nil {
			fmt.Printf("snap: could not retrieve data: %s", err)
			os.Exit(1)
		}

		format, err := rootCmd.PersistentFlags().GetString("format")
		if format == "json" {
			printJSON(response)
			return
		}

		if format == "table" {
			printActiveSnapLogDetails(response, logID, format)
			return
		}

		// unknown format - return the raw response
		fmt.Printf("Raw response:\n%s\n", string(response))
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.AddCommand(logDetailsCmd)
}
