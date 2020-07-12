package print

import (
	"encoding/json"

	"github.com/snapmaster-io/snap/pkg/utils"
)

// Status extract the status and prints out the error message, if the status is an error
func Status(response []byte) {
	// unmarshal into the SnapStatus struct, to get "Status" and
	// flatten the property set of the response
	var status map[string]string
	json.Unmarshal(response, &status)

	utils.PrintStatus(status["status"], status["message"])
}
