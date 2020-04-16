package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/snapmaster-io/snap/pkg/api"
	"github.com/spf13/cobra"
)

// snapsCmd represents the snaps command
var snapsCmd = &cobra.Command{
	Use:   "snaps [subcommand]",
	Short: "Manage the user's snaps",
	Long:  `Manage the user's snaps.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

// createSnapCmd represents the create snap subcommand
var createSnapCmd = &cobra.Command{
	Use:   "create [definition-file.yaml]",
	Short: "Create a snap in the user's namespace from a yaml definition file",
	Long:  `Create a snap in the user's namespace from a yaml definition file.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve yaml file as the first argument
		snapFile := args[0]

		// read the file contents
		contents, err := ioutil.ReadFile(snapFile)
		if err != nil {
			fmt.Printf("snap: could not read snap definition file %s\nerror: %s\n", snapFile, err)
			os.Exit(1)
		}

		data := make(map[string]interface{})
		data["action"] = "create"
		data["definition"] = string(contents)
		processSnapCommand(data)
	},
}

// deleteSnapCmd represents the delete snap subcommand
var deleteSnapCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a snap from the user's namespace",
	Long:  `Delete a snap from the user's namespace.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve snapID as the first argument
		snapID := args[0]
		data := make(map[string]interface{})
		data["action"] = "delete"
		data["snapId"] = snapID
		processSnapCommand(data)
	},
}

// forkSnapCmd represents the fork snap subcommand
var forkSnapCmd = &cobra.Command{
	Use:   "fork",
	Short: "Forks a public snap into the user's namespace",
	Long:  `Forks a public snap into the user's namespace.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve snapID as the first argument
		snapID := args[0]
		data := make(map[string]interface{})
		data["action"] = "fork"
		data["snapId"] = snapID
		processSnapCommand(data)
	},
}

// getSnapCmd represents the get snap subcommand
var getSnapCmd = &cobra.Command{
	Use:   "get [snap ID]",
	Short: "Get a description of a snap",
	Long:  `Get a description of a snap.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve snapID as the first argument
		snapID := args[0]
		path := fmt.Sprintf("/snaps/%s", snapID)

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

		var snap map[string]string
		json.Unmarshal(response, &snap)
		text := snap["text"]
		fmt.Printf(text)
	},
}

// listSnapsCmd represents the list snaps subcommand
var listSnapsCmd = &cobra.Command{
	Use:   "list",
	Short: "List the user's snaps",
	Long:  `List the user's snaps.`,
	Run: func(cmd *cobra.Command, args []string) {

		// execute the API call
		response, err := api.Get("/snaps")
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
			printSnapsTable(response)
			return
		}

		// unknown format - return the raw response
		fmt.Printf("Raw response:\n%s\n", string(response))
	},
}

// publishSnapCmd represents the publish snap subcommand
var publishSnapCmd = &cobra.Command{
	Use:   "publish",
	Short: "Makes a user's snap public and discoverable by others in the gallery",
	Long:  `Makes a user's snap public and discoverable by others in the gallery.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve snapID as the first argument
		snapID := args[0]
		data := make(map[string]interface{})
		data["action"] = "edit"
		data["snapId"] = snapID
		data["private"] = false
		processSnapCommand(data)
	},
}

// unpublishSnapCmd represents the publish snap subcommand
var unpublishSnapCmd = &cobra.Command{
	Use:   "unpublish",
	Short: "Makes a user's snap private",
	Long:  `Makes a user's snap private.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve snapID as the first argument
		snapID := args[0]
		data := make(map[string]interface{})
		data["action"] = "edit"
		data["snapId"] = snapID
		data["private"] = true
		processSnapCommand(data)
	},
}

func init() {
	rootCmd.AddCommand(snapsCmd)
	snapsCmd.AddCommand(createSnapCmd)
	snapsCmd.AddCommand(deleteSnapCmd)
	snapsCmd.AddCommand(forkSnapCmd)
	snapsCmd.AddCommand(getSnapCmd)
	snapsCmd.AddCommand(listSnapsCmd)
	snapsCmd.AddCommand(publishSnapCmd)
	snapsCmd.AddCommand(unpublishSnapCmd)
}

func processSnapCommand(data map[string]interface{}) {
	path := "/snaps"
	action := data["action"]
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

	if action == "delete" {
		printStatus(response)
	} else {
		printSnapStatus(response)
	}
}
