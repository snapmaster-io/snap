package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/snapmaster-io/snap/pkg/api"
	"github.com/snapmaster-io/snap/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

// connectionsCmd represents the connections command
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

// disconnectToolCmd represents the disconnect tool subcommand
var disconnectToolCmd = &cobra.Command{
	Use:   "disconnect [tool name]",
	Short: "Disconnect a tool and remove all credential sets associated with it",
	Long:  `Disconnect a tool and remove all credential sets associated with it.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve tool as the first argument
		tool := args[0]
		data := make(map[string]interface{})
		data["action"] = "remove"
		data["provider"] = tool
		processConnectionCommand("/connections", tool, data)
	},
}

// getConnectionCmd represents the get connection subcommand
var getConnectionCmd = &cobra.Command{
	Use:   "get [connection name]",
	Short: "Get credential sets associated with a connection",
	Long:  `Get credential sets associated with a connection.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve connection as the first argument
		connection := args[0]

		// execute the API call
		path := fmt.Sprintf("/entities/%s", connection)
		response, err := api.Get(path)
		if err != nil {
			utils.PrintErrorMessage("could not retrieve data", err)
			os.Exit(1)
		}

		format, err := rootCmd.PersistentFlags().GetString("format")
		if format == "json" {
			printJSON(response)
			return
		}

		if format == "table" {
			printCredentialsTable(response, connection)
			return
		}

		// other action - return the raw response
		printRawResponse(response)
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
			utils.PrintErrorMessage("could not retrieve data", err)
			os.Exit(1)
		}

		format, err := rootCmd.PersistentFlags().GetString("format")
		if format == "json" {
			printJSON(response)
			return
		}

		if format == "table" {
			printConnectionsTable(response)
			return
		}

		// unknown format - return the raw response
		printRawResponse(response)
	},
}

func init() {
	rootCmd.AddCommand(connectionsCmd)
	connectionsCmd.AddCommand(disconnectToolCmd)
	connectionsCmd.AddCommand(getConnectionCmd)
	connectionsCmd.AddCommand(listConnectionsCmd)
}

func processConnectionCommand(path string, connection string, data map[string]interface{}) {
	action := data["action"]

	payload, err := json.Marshal(data)
	if err != nil {
		utils.PrintErrorMessage("could not serialize payload into JSON", err)
		os.Exit(1)
	}

	// execute the API call
	response, err := api.Post(path, payload)
	if err != nil {
		utils.PrintErrorMessage("could not retrieve data", err)
		os.Exit(1)
	}

	format, err := rootCmd.PersistentFlags().GetString("format")
	if format == "json" {
		printJSON(response)
		return
	}

	if action == "remove" {
		// if credential sets were returned, display them
		num := gjson.GetBytes(response, "#").Int()
		if num > 0 {
			utils.PrintMessage(fmt.Sprintf("successfully removed credential-set %s from tool %s", data["id"], connection))
			printCredentialsTable(response, connection)
			return
		}

		printStatus(response)
		return
	}

	// other action - return the raw response
	printRawResponse(response)
}
