package cmd

import (
	"fmt"
	"os"

	"github.com/snapmaster-io/snap/pkg/api"
	"github.com/snapmaster-io/snap/pkg/print"
	"github.com/snapmaster-io/snap/pkg/utils"
	"github.com/spf13/cobra"
)

// credentialsCmd represents the credential set command
var credentialsCmd = &cobra.Command{
	Use:   "credential-set [subcommand]",
	Short: "Manage tool credential-sets",
	Long:  `Manage tool credential-sets.`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

// credentialsAddCmd represents the add credential set add command
var credentialsAddCmd = &cobra.Command{
	Use:   "add [tool name] [credential-set-name] [credential file name]",
	Short: "Add a credential set and associate it with this tool",
	Long: `Add a credential set and associate it with this tool.
	
	If only the tool name is passed in, the command will prompt for credential information.
	
	If a credential-set name and credential file name are provided, the command will create 
	a named credential-set with those parameters.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve tool as the first argument
		tool := args[0]

		// retrieve the parameter definitions from the API
		jsonPath := fmt.Sprintf("data.#(provider==%s).definition.connection.connectionInfo", tool)
		credentials := getParameterDescriptions("/connections", jsonPath)

		utils.PrintMessage(fmt.Sprintf("adding credential-set for %s", tool))

		// if no credentials file supplied, prompt for parameters
		if len(args) == 1 {
			// input the parameters and store their values in the credentials slice of maps
			inputParameters(credentials)
		} else {
			if len(args) != 3 {
				cmd.Help()
				os.Exit(1)
			}

			// populate parameter values based on command-line arguments
			credentialName := args[1]
			credentialsFile := args[2]
			readParametersFromFile(credentials, credentialName, credentialsFile)
		}

		// make the POST call to the API
		path := fmt.Sprintf("/entities/%s", tool)
		processConnectCommand(tool, path, credentials)
	},
}

// credentialsListCmd represents the credential set list command
var credentialsListCmd = &cobra.Command{
	Use:   "list [tool name]",
	Short: "List credential sets associated with this tool",
	Long:  `List credential sets associated with this tool.`,
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
			print.JSON(response)
			return
		}

		if format == "table" {
			print.CredentialsTable(response, connection)
			return
		}

		// other action - return the raw response
		print.RawResponse(response)
	},
}

// credentialsRemoveCmd represents the credential set remove command
var credentialsRemoveCmd = &cobra.Command{
	Use:   "remove [tool name] [credential-set-name]",
	Short: "Removes a credential set associated with this tool",
	Long: `Removes a credential set associated with this tool.
	
	Both the tool name and credential-set name arguments are required.`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// retrieve arguments
		tool := args[0]
		credentials := args[1]

		data := make(map[string]interface{})
		data["action"] = "remove"
		data["id"] = credentials

		// make the POST call to the API
		path := fmt.Sprintf("/entities/%s", tool)
		processConnectionCommand(path, tool, data)
	},
}

func init() {
	connectionsCmd.AddCommand(credentialsCmd)
	credentialsCmd.AddCommand(credentialsAddCmd)
	credentialsCmd.AddCommand(credentialsListCmd)
	credentialsCmd.AddCommand(credentialsRemoveCmd)
}
