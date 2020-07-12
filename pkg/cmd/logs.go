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

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get all logs for the current user",
	Long:  `Get all logs for the current user.`,
	Run: func(cmd *cobra.Command, args []string) {
		// execute the API call
		response, err := api.Get("/logs")
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
			print.ActiveSnapLogsTable(response)
			return
		}

		// unknown format - return the raw response
		print.RawResponse(response)
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
			utils.PrintErrorMessage("could not retrieve data", err)
			os.Exit(1)
		}

		format, err := rootCmd.PersistentFlags().GetString("format")
		if format == "json" {
			// select the entry that matches the log ID
			logEntry := gjson.GetBytes(response, fmt.Sprintf("data.#(timestamp==%s)|@pretty", logID)).Raw

			if logEntry == "" {
				utils.PrintError(fmt.Sprintf("log ID %s not found", logID))
				return
			}

			// print the log entry
			print.JSONString(logEntry)
			return
		}

		if format == "table" {
			print.ActiveSnapLogDetails(response, logID, format)
			return
		}

		// unknown format - return the raw response
		print.RawResponse(response)
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.AddCommand(logDetailsCmd)
}
