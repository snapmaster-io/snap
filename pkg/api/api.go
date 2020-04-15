package api

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

// Get calls the API at the relative path, and returns the data retrieved or an error
func Get(path string) ([]byte, error) {
	return call(path, "GET", nil)
}

// Post calls the API at the relative path with the payload, and returns the data retrieved or an error
func Post(path string, payload []byte) ([]byte, error) {
	r := bytes.NewReader(payload)
	return call(path, "POST", r)
}

func call(path string, verb string, payload interface{ io.Reader }) ([]byte, error) {
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

	// construct the URL and request
	url := apiURL + path
	req, err := http.NewRequest(verb, url, payload)
	if err != nil {
		fmt.Printf("snap: could not create request with URL %s: %s\n", url, err)
		os.Exit(1)
	}

	// add headers and execute the request
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", accessToken))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("snap: could not execute HTTP request with URL %s\n", url, err)
		os.Exit(1)
	}

	// process the response
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("snap: error reading HTTP response from HTTP request against %s\n", url, err)
		os.Exit(1)
	}

	return contents, nil
}
