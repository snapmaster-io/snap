package cmd

import (
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

// getActiveSnapCmd represents the get active snap subcommand
var getActiveSnapCmd = &cobra.Command{
	Use:   "get [active snap ID]",
	Short: "Get the state of an active snap",
	Long:  `Get a description of a snap.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve snapID as the first argument
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

func init() {
	rootCmd.AddCommand(activeSnapsCmd)
	activeSnapsCmd.AddCommand(listActiveSnapsCmd)
	activeSnapsCmd.AddCommand(getActiveSnapCmd)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// snapsCmd.PersistentFlags().String("foo", "", "A help for foo")
	// snapsCmd.Flags().StringP("username", "u", "", "Username")
}
