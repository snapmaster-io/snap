package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/snapmaster-io/snap/pkg/api"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

// connectionsCmd represents the tools command
var connectionsCmd = &cobra.Command{
	Use:   "connections [subcommand]",
	Short: "Manage connections to tools",
	Long:  `Manage connections to tools.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

// connectToolCmd represents the connect tool subcommand
var connectToolCmd = &cobra.Command{
	Use:   "connect [tool name]",
	Short: "Connect a tool",
	Long:  `Connect a tool.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve tool name as the first argument
		tool := args[0]

		data := make(map[string]interface{})
		data["action"] = "add"
		data["provider"] = tool
		processConnectionCommand(data)
	},
}

// disconnectToolCmd represents the disconnect tool subcommand
var disconnectToolCmd = &cobra.Command{
	Use:   "disconnect [tool name]",
	Short: "Disconnect a tool",
	Long:  `Disconnect a tool.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve tool as the first argument
		tool := args[0]
		data := make(map[string]interface{})
		data["action"] = "remove"
		data["provider"] = tool
		processConnectionCommand(data)
	},
}

// getConnectionCmd represents the get connection subcommand
var getConnectionCmd = &cobra.Command{
	Use:   "get [connection name]",
	Short: "Get a description of a connection",
	Long:  `Get a description of a connection.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve snapID as the first argument
		connection := args[0]

		// execute the API call
		response, err := api.Get("/connections")
		if err != nil {
			fmt.Printf("snap: could not retrieve data: %s", err)
			os.Exit(1)
		}

		format, err := rootCmd.PersistentFlags().GetString("format")
		if format == "json" {
			// select the entry that matches the provider name
			toolDescription := gjson.GetBytes(response, fmt.Sprintf("#(provider==%s)|@pretty", connection)).Raw
			// print the tool description
			fmt.Print(toolDescription)
			return
		}

		// select the entry that matches the provider name
		toolDescription := gjson.GetBytes(response, fmt.Sprintf("#(provider==%s).definition", connection)).Raw
		// print the tool description
		fmt.Print(toolDescription)
	},
}

// listConnectionsCmd represents the list tools subcommand
var listConnectionsCmd = &cobra.Command{
	Use:   "list",
	Short: "List the user's connections",
	Long:  `List the user's connections.`,
	Run: func(cmd *cobra.Command, args []string) {

		// execute the API call
		response, err := api.Get("/connections")
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
			printToolsTable(response)
			return
		}

		// unknown format - return the raw response
		fmt.Printf("Raw response:\n%s\n", string(response))
	},
}

func init() {
	rootCmd.AddCommand(connectionsCmd)
	rootCmd.AddCommand(connectToolCmd)
	connectionsCmd.AddCommand(disconnectToolCmd)
	connectionsCmd.AddCommand(getConnectionCmd)
	connectionsCmd.AddCommand(listConnectionsCmd)
}

func processConnectionCommand(data map[string]interface{}) {
	path := "/connections"
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

	if action == "remove" {
		printStatus(response)
	} else {
		printSnapStatus(response)
	}
}
