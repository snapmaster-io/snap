package print

import "github.com/snapmaster-io/snap/pkg/utils"

// JSON pretty-prints a JSON response
func JSON(response []byte) {
	// pretty-print the json
	utils.PrintJSON(response)
}

// JSONString pretty-prints a string that contains JSON
func JSONString(response string) {
	// pretty-print the json
	utils.PrintJSON([]byte(response))
}
