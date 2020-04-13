package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

// Get calls the API at the relative path, and returns the data retrieved or an error
func Get(path string) ([]map[string]interface{}, error) {
	// retrieve access token
	accessToken := viper.GetString("AccessToken")
	if len(accessToken) < 1 {
		fmt.Println("snap: login required before executing this command")
		os.Exit(1)
	}

	// retrieve API URL
	apiURL := viper.GetString("APIURL")
	if len(apiURL) < 1 {
		fmt.Println("snap: API URL required but not found")
		os.Exit(1)
	}

	// create the request and execute it
	req, _ := http.NewRequest("GET", url)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer: %s", accessToken))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// process the response
	defer res.Body.Close()
	var responseData []map[string]interface{}
	body, _ := ioutil.ReadAll(res.Body)

	// unmarshal the json into a string map
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Printf("snap: error unmarshalling JSON: %s\n", err)
		os.Exit(1)
	}

	return responseData, nil
}
