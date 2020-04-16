package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/snapmaster-io/snap/pkg/api"
	"github.com/spf13/cobra"
)

// activeSnapsCmd represents the snaps command
var activeSnapsCmd = &cobra.Command{
	Use:   "active [subcommand]",
	Short: "Manage the user's active snaps",
	Long:  `Manage the user's active snaps.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

// deactivateSnapCmd represents the deactivate snap subcommand
var deactivateSnapCmd = &cobra.Command{
	Use:   "deactivate [active snap ID]",
	Short: "Deactivate a snap",
	Long: `Deactivate a snap.
	
	Note that once an active snap is deactivated, ALL LOGS ARE DELETED.
	
	If you want to stop the active snap from triggering, use the pause subcommand.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve activeSnapID as the first argument
		activeSnapID := args[0]
		processActiveCommand(activeSnapID, "deactivate")
	},
}

// getActiveSnapCmd represents the get active snap subcommand
var getActiveSnapCmd = &cobra.Command{
	Use:   "get [active snap ID]",
	Short: "Get the state of an active snap",
	Long:  `Get the state of an active snap.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve activeSnapID as the first argument
		activeSnapID := args[0]
		path := fmt.Sprintf("/activesnaps/%s", activeSnapID)

		// execute the API call
		response, err := api.Get(path)
		if err != nil {
			fmt.Printf("snap: could not retrieve data: %s", err)
			os.Exit(1)
		}

		format, err := rootCmd.PersistentFlags().GetString("format")
		if format == "json" {
			printJSON(response)
			return
		}

		printActiveSnap(response)
	},
}

// getActiveSnapLogsCmd represents the get active snap logs subcommand
var getActiveSnapLogsCmd = &cobra.Command{
	Use: `logs [active snap ID] [flags]
  snap active logs [active snap ID] details [log ID] --format={stdout, stderr} (default "stdout")`,
	Short: "Get the logs of an activated snap",
	Long: `Get the logs of an activated snap.

snap active logs [active snap ID] will return a table (or json) with all log entries for this active snap.

You can obtain a log ID for a specific log entry using this command.

snap active logs [active snap ID] details [log ID] will return the output for each action - either stdout or stderr.
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve activeSnapID as the first argument
		activeSnapID := args[0]
		logID := ""

		// validate second form of command (logs [activeSnapID] details [logID])
		if len(args) > 2 {
			if args[1] == "details" {
				logID = args[2]
			} else {
				cmd.Help()
				os.Exit(1)
			}
		}

		// execute the API call
		path := fmt.Sprintf("/logs/%s", activeSnapID)
		response, err := api.Get(path)
		if err != nil {
			fmt.Printf("snap: could not retrieve data: %s", err)
			os.Exit(1)
		}

		format, err := rootCmd.PersistentFlags().GetString("format")
		if format == "json" {
			printJSON(response)
			return
		}

		if len(args) > 2 && logID != "" {
			printActiveSnapLogDetails(response, logID, format)
		} else {
			printActiveSnapLogs(response)
		}
	},
}

// listActiveSnapsCmd represents the list snaps subcommand
var listActiveSnapsCmd = &cobra.Command{
	Use:   "list",
	Short: "List the user's active snaps",
	Long:  `List the user's active snaps.`,
	Run: func(cmd *cobra.Command, args []string) {

		// execute the API call
		response, err := api.Get("/activesnaps")
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
			printActiveSnapsTable(response)
			return
		}

		// unknown format - return the raw response
		fmt.Printf("Raw response:\n%s\n", string(response))
	},
}

// pauseActiveSnapCmd represents the pause active snap subcommand
var pauseActiveSnapCmd = &cobra.Command{
	Use:   "pause [active snap ID]",
	Short: "Pause an active snap",
	Long:  `Pause an active snap.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve activeSnapID as the first argument
		activeSnapID := args[0]
		processActiveCommand(activeSnapID, "pause")
	},
}

// resumeActiveSnapCmd represents the resume active snap subcommand
var resumeActiveSnapCmd = &cobra.Command{
	Use:   "resume [active snap ID]",
	Short: "Resume an active snap",
	Long:  `Resume an active snap.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve activeSnapID as the first argument
		activeSnapID := args[0]
		processActiveCommand(activeSnapID, "resume")
	},
}

func init() {
	rootCmd.AddCommand(activeSnapsCmd)
	activeSnapsCmd.AddCommand(deactivateSnapCmd)
	activeSnapsCmd.AddCommand(getActiveSnapCmd)
	activeSnapsCmd.AddCommand(getActiveSnapLogsCmd)
	activeSnapsCmd.AddCommand(resumeActiveSnapCmd)
	activeSnapsCmd.AddCommand(listActiveSnapsCmd)
	activeSnapsCmd.AddCommand(pauseActiveSnapCmd)
	activeSnapsCmd.AddCommand(resumeActiveSnapCmd)
}

func processActiveCommand(activeSnapID string, action string) {
	path := "/activesnaps"

	data := make(map[string]string)
	data["action"] = action
	data["snapId"] = activeSnapID
	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("snap: could not serialize payload into JSON: %s\n", err)
		os.Exit(1)
	}

	// execute the API call
	response, err := api.Post(path, payload)
	if err != nil {
		fmt.Printf("snap: could not retrieve data: %s\n", err)
		os.Exit(1)
	}

	format, err := rootCmd.PersistentFlags().GetString("format")
	if format == "json" {
		printJSON(response)
		return
	}

	if action != "deactivate" {
		printActiveSnapStatus(response)
	} else {
		printStatus(response)
	}
}

func processGetLogDetailsCommand(activeSnapID string, logID string) {
	// execute the API call
	path := fmt.Sprintf("/logs/%s/%s", activeSnapID, logID)
	response, err := api.Get(path)
	if err != nil {
		fmt.Printf("snap: could not retrieve data: %s", err)
		os.Exit(1)
	}

	format, err := rootCmd.PersistentFlags().GetString("format")
	if format == "json" {
		printJSON(response)
		return
	}

	printActiveSnapLogDetails(response, logID, format)
}