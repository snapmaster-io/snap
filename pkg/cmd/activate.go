package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/snapmaster-io/snap/pkg/api"
	"github.com/snapmaster-io/snap/pkg/utils"
	"github.com/spf13/cobra"
)

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate [snap ID]",
	Short: "Activate a snap",
	Long: `Activate a snap.
	
	If only the snap ID is passed in, the command will prompt for parameters.
	
	If the parameter file was provided, those parameter values will be used to activate the snap.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		snapID := args[0]
		paramsFile, err := cmd.Flags().GetString("params-file")
		if err != nil {
			utils.PrintError(fmt.Sprintf("couldn't read params-file %s\nerror: %s\n", paramsFile, err))
			os.Exit(1)
		}

		utils.PrintMessage(fmt.Sprintf("activating snap %s", snapID))

		var params []map[string]string

		// if no params file supplied, need to get the snap definition and prompt for parameters
		if paramsFile == "" {
			params = getSnapParameters(snapID, "snaps", "data.parameters")
			inputParameters(params)
		}

		// make the POST call to the API
		processActivateCommand(snapID, "activate", params)
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
	activateCmd.Flags().StringP("params-file", "p", "", "a yaml file that defines snap parameter values")
}

func getSnapParameters(snapID string, path string, jsonPath string) []map[string]string {
	urlpath := fmt.Sprintf("/%s/%s", path, snapID)

	params := getParameterDescriptions(urlpath, jsonPath)
	return params
}

func processActivateCommand(snapID string, action string, params []map[string]string) {
	path := "/activesnaps"

	// set up the data map
	data := make(map[string]interface{})
	data["action"] = action
	data["snapId"] = snapID
	data["params"] = params

	payload, err := json.Marshal(data)
	if err != nil {
		utils.PrintError(fmt.Sprintf("could not serialize payload into JSON\nerror: %s\n", err))
		os.Exit(1)
	}

	// execute the API call
	response, err := api.Post(path, payload)
	if err != nil {
		utils.PrintError(fmt.Sprintf("could not retrieve data\nerror: %s\n", err))
		os.Exit(1)
	}

	format, err := rootCmd.PersistentFlags().GetString("format")
	if format == "json" {
		printJSON(response)
		return
	}

	printActiveSnapStatus(response)
}

func obtainSnapParameters(snapID string, path string, jsonPath string) []map[string]string {
	params := getSnapParameters(snapID, path, jsonPath)

	// create a new reader from stdin
	reader := bufio.NewReader(os.Stdin)

	// get values for each parameter and store them in the same map
	for i, param := range params {
		fmt.Printf("%s (%s): ", param["name"], param["description"])
		text, _ := reader.ReadString('\n')

		// store the parameter value without the last character (\n)
		params[i]["value"] = text[:len(text)-1]
	}

	return params
}
