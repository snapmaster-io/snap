package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
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
			printTable(response)
			return
		}

		// unknown format - return the raw response
		fmt.Printf("Raw response:\n%s\n", string(response))
	},
}

// listSnapsCmd represents the get snap subcommand
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

func init() {
	rootCmd.AddCommand(snapsCmd)
	snapsCmd.AddCommand(listSnapsCmd)
	snapsCmd.AddCommand(getSnapCmd)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// snapsCmd.PersistentFlags().String("foo", "", "A help for foo")
	// snapsCmd.Flags().StringP("username", "u", "", "Username")
}

func printJSON(response []byte) {
	// pretty-print the json
	output := &bytes.Buffer{}
	err := json.Indent(output, response, "", "  ")
	if err != nil {
		fmt.Printf("snap: could not format response as json")
		os.Exit(1)
	}
	fmt.Println(output.String())
}

func printTable(response []byte) {
	var snaps []map[string]string
	json.Unmarshal(response, &snaps)

	// write out the table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Snap ID", "Description", "Trigger"})
	for _, snap := range snaps {
		t.AppendRow(table.Row{snap["snapId"], snap["description"], snap["trigger"]})
	}
	t.Render()
}
