package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/snapmaster-io/snap/pkg/api"
	"github.com/snapmaster-io/snap/pkg/utils"
	"github.com/tidwall/gjson"
)

// getParameterDescriptions retrieves the definitions via the API call and creates
// a slice of maps, each containing the name and description of a parameter
func getParameterDescriptions(path string, jsonPath string) []map[string]string {
	// execute the API call
	response, err := api.Get(path)
	if err != nil {
		utils.PrintErrorMessage(fmt.Sprintf("could not retrieve snap %s", path), err)
		os.Exit(1)
	}

	var status map[string]string
	json.Unmarshal(response, &status)
	if status["status"] != "success" {
		utils.PrintStatus(status["status"], status["message"])
		os.Exit(1)
	}

	// json.Unmarshal doesn't do very well with nested arrays / maps in json
	// gjson is a bit better but still a bit limited... so need to iterate over results and create a new []map
	// get an array of names and descriptions
	responseString := string(response)
	names := gjson.Get(responseString, fmt.Sprintf("%s.#.name", jsonPath)).Array()
	descriptions := gjson.Get(responseString, fmt.Sprintf("%s.#.description", jsonPath)).Array()
	types := gjson.Get(responseString, fmt.Sprintf("%s.#.type", jsonPath)).Array()

	// create a slice of maps which will contain parameter names and descriptions
	params := make([]map[string]string, len(names))
	for i, name := range names {
		params[i] = make(map[string]string)
		params[i]["name"] = name.String()
		if len(descriptions) > i {
			params[i]["description"] = descriptions[i].String()
		} else {
			params[i]["description"] = ""
		}
		if len(types) > i {
			params[i]["type"] = types[i].String()
		}
	}

	return params
}

// inputParameters expects a slice of maps containing name and description keys
// it stores inputted values in the value key
func inputParameters(params []map[string]string) {
	// create a new reader from stdin
	reader := bufio.NewReader(os.Stdin)

	// get values for each parameter and store them in the same map
	for i, param := range params {
		fmt.Printf("%s (%s): ", param["name"], param["description"])
		text, _ := reader.ReadString('\n')

		// store the parameter value without the last character (\n)
		params[i]["value"] = text[:len(text)-1]
	}
}

// readParametersFromFile expects a slice of maps containing name and description keys
// it stores the credential name in the "name" parameter and retrieves the contents of
// the filename into any other parameter
func readParametersFromFile(params []map[string]string, credentialName string, credentialsFile string) {
	contents, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		utils.PrintErrorMessage(fmt.Sprintf("could not read credentials file %s", credentialsFile), err)
		os.Exit(1)
	}

	for _, param := range params {
		if param["type"] == "name" {
			param["value"] = credentialName
		} else {
			param["value"] = string(contents)
		}
	}
}
