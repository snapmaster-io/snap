package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/snapmaster-io/snap/pkg/api"
	"github.com/snapmaster-io/snap/pkg/print"
	"github.com/snapmaster-io/snap/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect [tool name] [credential-set name] [credential file name]",
	Short: "Connect a tool",
	Long: `Connect a tool.
	
	If only the tool name is passed in, the command will prompt for credential information.
	
	If a credential-set name and credential file name are provided, the command will create 
	a default connection as well as a named credential-set with those parameters.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tool := args[0]

		// retrieve the parameter definitions from the connections API
		jsonPath := fmt.Sprintf("#(provider==%s).definition.connection.connectionInfo", tool)
		credentials := getParameterDescriptions("/connections", jsonPath)

		utils.PrintMessage(fmt.Sprintf("connecting %s", tool))

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
		processConnectCommand(tool, "/connections", credentials)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}

func processConnectCommand(tool string, path string, params []map[string]string) {
	// set up the data map
	data := make(map[string]interface{})
	data["action"] = "add"
	data["provider"] = tool
	data["connectionInfo"] = params

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
		print.JSON(response)
		return
	}

	// if credential sets were returned, display them
	num := gjson.GetBytes(response, "#").Int()
	if num > 0 {
		utils.PrintMessage(fmt.Sprintf("connected %s and stored credentials", tool))
		print.CredentialsTable(response, tool)
		return
	}

	print.Status(response)
}
