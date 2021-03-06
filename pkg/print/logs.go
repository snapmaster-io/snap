package print

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/snapmaster-io/snap/pkg/utils"
)

// ActiveSnapActionsLog defines the fields to unmarshal for action logs
type ActiveSnapActionsLog struct {
	Provider string                     `json:"provider"`
	Action   string                     `json:"action"`
	State    string                     `json:"state"`
	Output   ActiveSnapActionsLogOutput `json:"output"`
}

// ActiveSnapActionsLogOutput defines the fields to unmarshal for action logs
type ActiveSnapActionsLogOutput struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

/*
type ActiveSnapActionsLogOutput struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Error  string `json:"error"`
}
*/

// ActiveSnapLog defines the fields to unmarshal for an active snap's logs
type ActiveSnapLog struct {
	LogID        int64                  `json:"timestamp"`
	ActiveSnapID string                 `json:"activeSnapId"`
	SnapID       string                 `json:"snapID"`
	State        string                 `json:"state"`
	Trigger      string                 `json:"trigger"`
	Event        string                 `json:"event"`
	Actions      []ActiveSnapActionsLog `json:"actions"`
}

// ActiveSnapLogsResponse defines the fields to unmarshal from getting all logs for an active snap
type ActiveSnapLogsResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    []ActiveSnapLog `json:"data"`
}

// ActiveSnapLogsTable prints out the active snap logs response as a table
func ActiveSnapLogsTable(response []byte) {
	// extract the logs and check for errors
	activeSnapLogs, err := extractLogs(response)
	if err != nil {
		return
	}

	// grab the SnapIP, ActiveSnapID, and Trigger from the first record
	activeSnapInstance := activeSnapLogs[0]
	activeSnapID := activeSnapInstance.ActiveSnapID
	snapID := activeSnapInstance.SnapID
	trigger := activeSnapInstance.Trigger
	event := activeSnapInstance.Event

	// write out the table of properties
	t := table.NewWriter()
	t.SetTitle(fmt.Sprintf(
		"Logs for Snap ID %s\nActive Snap ID %s, triggered by %s:%s",
		snapID, activeSnapID, trigger, event))
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Log ID", "Timestamp", "State"})
	for _, logEntry := range activeSnapLogs {
		timestamp := time.Unix(logEntry.LogID/1000, 0)
		t.AppendRow(table.Row{logEntry.LogID, timestamp, logEntry.State})
	}
	t.SetStyle(tableStyle)
	t.Style().Title.Align = text.AlignCenter
	t.Render()
}

// ActiveSnapLogDetails prints out the active log details response
func ActiveSnapLogDetails(response []byte, logID string, format string) {
	// extract the logs and check for errors
	activeSnapLogs, err := extractLogs(response)
	if err != nil {
		return
	}

	// grab the SnapIP, ActiveSnapID, and Trigger from the first record
	activeSnapInstance := activeSnapLogs[0]
	activeSnapID := activeSnapInstance.ActiveSnapID
	snapID := activeSnapInstance.SnapID

	var logEntry ActiveSnapLog

	// find the entry with the right logID
	found := false
	for k, v := range activeSnapLogs {
		logIDasInt64, _ := strconv.ParseInt(logID, 10, 64)
		if v.LogID == logIDasInt64 {
			logEntry = activeSnapLogs[k]
			found = true
		}
	}

	// check for log entry not found
	if found == false {
		utils.PrintError(fmt.Sprintf("log ID %s not found for active Snap ID %s", logID, activeSnapID))
		return
	}

	// write out general information
	t := table.NewWriter()
	t.SetTitle("Action log details")
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Snap ID", "Active Snap ID", "Log ID", "State"})
	t.AppendRow(table.Row{snapID, activeSnapID, logID, logEntry.State})
	t.SetStyle(tableStyle)
	t.Style().Title.Align = text.AlignCenter
	t.Render()

	fmt.Println("\nAction details:")

	// write out the table of properties
	for _, action := range logEntry.Actions {
		fmt.Println()
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)

		if action.Output.Status != "success" {
			t.AppendHeader(table.Row{"Provider", "Action", "State", "Message"})
			t.AppendRow(table.Row{action.Provider, action.Action, action.State, action.Output.Message})
		} else {
			t.AppendHeader(table.Row{"Provider", "Action", "State"})
			t.AppendRow(table.Row{action.Provider, action.Action, action.State})
		}
		t.SetStyle(actionTableStyle)
		t.Render()

		// print the output of the operation
		printOutput(action.Output)
	}
}

// do some common processing on the ActiveSnapLogResponse, and extract the logs
func extractLogs(response []byte) ([]ActiveSnapLog, error) {
	// unmarshal into the ActiveSnapLogsResponse struct, to flatten the property set
	var activeSnapLogsResponse ActiveSnapLogsResponse
	json.Unmarshal(response, &activeSnapLogsResponse)

	// check for errors
	if activeSnapLogsResponse.Status == "error" {
		utils.PrintStatus(activeSnapLogsResponse.Status, activeSnapLogsResponse.Message)
		return nil, errors.New("Error encountered")
	}

	// extract the active snap logs
	activeSnapLogs := activeSnapLogsResponse.Data

	// check for no rows
	if len(activeSnapLogs) < 1 {
		utils.PrintError("no logs found for this active snap")
		return nil, errors.New("No logs found for this active snap")
	}

	return activeSnapLogs, nil
}

// print the output for the action
func printOutput(output ActiveSnapActionsLogOutput) {
	data := output.Data
	if output.Status == "success" {
		// if the action supports the stdout/stderr protocol, print those out
		if data["stdout"] != nil || data["stderr"] != nil {
			if data["stdout"] != nil {
				// write out stdout
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"Stdout"})
				t.SetStyle(tableStyle)
				t.Render()
				fmt.Printf("%s\n", data["stdout"])
			}

			if data["stderr"] != nil {
				// write out stderr
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"Stderr"})
				t.SetStyle(tableStyle)
				t.Render()
				fmt.Printf("%s\n", data["stderr"])
			}
		} else {
			// print out the JSON response
			payload, err := json.Marshal(data)
			if err != nil {
				utils.PrintErrorMessage("could not serialize payload into JSON", err)
			} else {
				utils.PrintJSON(payload)
			}
		}
	}
}
