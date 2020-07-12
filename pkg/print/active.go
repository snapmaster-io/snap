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

// ActiveSnap defines the fields to unmarshal for an active snap
type ActiveSnap struct {
	ActiveSnapID     string `json:"activeSnapId"`
	SnapID           string `json:"snapID"`
	State            string `json:"state"`
	Provider         string `json:"provider"`
	Activated        int64  `json:"activated"`
	ExecutionCounter int    `json:"executionCounter"`
	ErrorCounter     int    `json:"errorCounter"`
}

// ActiveSnapActionsLog defines the fields to unmarshal for action logs
type ActiveSnapActionsLog struct {
	Provider string                     `json:"provider"`
	Action   string                     `json:"action"`
	State    string                     `json:"state"`
	Output   ActiveSnapActionsLogOutput `json:"output"`
}

// ActiveSnapActionsLogOutput defines the fields to unmarshal for action logs
type ActiveSnapActionsLogOutput struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Error  string `json:"error"`
}

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

// ActiveSnapsResponse defines the fields to unmarshal from getting all active snaps
type ActiveSnapsResponse struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    []ActiveSnap `json:"data"`
}

// ActiveSnapResponse defines the fields to unmarshal from a delete/pause/resume operation
type ActiveSnapResponse struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Data    ActiveSnap `json:"data"`
}

// ActiveSnapTable prints the active snap response as a table
func ActiveSnapTable(response []byte) {
	// unmarshal into the ActiveSnapResponse struct, to flatten the property set
	var activeSnapResponse ActiveSnapResponse
	json.Unmarshal(response, &activeSnapResponse)

	if activeSnapResponse.Status != "success" {
		utils.PrintStatus(activeSnapResponse.Status, activeSnapResponse.Message)
		return
	}

	activeSnap := activeSnapResponse.Data

	// re-marshal and unmarshal into a map, which can be iterated over as a {name, value} pair
	intermediateEntity, _ := json.Marshal(activeSnap)
	var entity map[string]interface{}
	json.Unmarshal(intermediateEntity, &entity)

	// TODO: sort / alphabetize the keys

	// write out the table of properties
	t := table.NewWriter()
	t.SetTitle(fmt.Sprintf("Active Snap %s", activeSnap.ActiveSnapID))
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Field", "Value"})
	for field, value := range entity {
		if field == "activated" {
			value = time.Unix(int64(value.(float64))/1000, 0)
		}
		t.AppendRow(table.Row{field, value})
	}
	t.SetStyle(tableStyle)
	t.Style().Title.Align = text.AlignCenter
	t.Render()
}

// ActiveSnapLogsTable prints out the active snap logs response as a table
func ActiveSnapLogsTable(response []byte) {
	// unmarshal into the ActiveSnapLogsResponse struct, to flatten the property set
	var activeSnapLogsResponse ActiveSnapLogsResponse
	json.Unmarshal(response, &activeSnapLogsResponse)

	// check for errors
	utils.PrintStatus(activeSnapLogsResponse.Status, activeSnapLogsResponse.Message)
	if activeSnapLogsResponse.Status == "error" {
		return
	}

	// extract the active snap logs
	activeSnapLogs := activeSnapLogsResponse.Data

	// check for no rows
	if len(activeSnapLogs) < 1 {
		utils.PrintError("no logs found for this snap")
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
	/*
		// unmarshal into the ActiveSnapLogsResponse struct, to flatten the property set
		var activeSnapLogsResponse ActiveSnapLogsResponse
		json.Unmarshal(response, &activeSnapLogsResponse)

		// check for errors
		utils.PrintStatus(activeSnapLogsResponse.Status, activeSnapLogsResponse.Message)
		if activeSnapLogsResponse.Status == "error" {
			return
		}

		// extract the active snap logs
		activeSnapLogs := activeSnapLogsResponse.Data

		// check for no rows
		if len(activeSnapLogs) < 1 {
			utils.PrintError("no logs found for this snap")
			return
		}
	*/
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
		utils.PrintError(fmt.Sprintf("log ID %s not found for active Snap ID %s\n", logID, activeSnapID))
		os.Exit(1)
	}

	// write out general information
	t := table.NewWriter()
	t.SetTitle("Action log details")
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Snap ID", "Active Snap ID", "Log ID"})
	t.AppendRow(table.Row{snapID, activeSnapID, logID})
	t.SetStyle(tableStyle)
	t.Style().Title.Align = text.AlignCenter
	t.Render()

	fmt.Println("\nAction details:")

	// write out the table of properties
	for _, action := range logEntry.Actions {
		fmt.Println()
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Provider", "Action", "State"})
		t.AppendRow(table.Row{action.Provider, action.Action, action.State})
		t.SetStyle(actionTableStyle)
		t.Render()

		// write out stdout
		t = table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Stdout"})
		t.SetStyle(tableStyle)
		t.Render()
		fmt.Printf("%s\n", action.Output.Stdout)

		// write out stderr
		t = table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Stderr"})
		t.SetStyle(tableStyle)
		t.Render()
		fmt.Printf("%s\n", action.Output.Stderr)
	}
}

// ActiveSnapStatusTable prints out the active snap status response as a table
func ActiveSnapStatusTable(response []byte) {
	// unmarshal into the ActiveSnapResponse struct, to get "Status" and
	// flatten the property set of the ActiveSnap
	var activeSnapResponse ActiveSnapResponse
	json.Unmarshal(response, &activeSnapResponse)

	utils.PrintStatus(activeSnapResponse.Status, activeSnapResponse.Message)

	// if the status indicates an error, there is no active snap to display
	if activeSnapResponse.Status != "success" {
		return
	}

	activeSnap := activeSnapResponse.Data

	// re-marshal and unmarshal into a map, which can be iterated over as a {name, value} pair
	intermediateEntity, _ := json.Marshal(activeSnap)
	var entity map[string]interface{}
	json.Unmarshal(intermediateEntity, &entity)

	// TODO: sort / alphabetize the keys

	// write out the table of properties
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Active Snap Values")
	t.AppendHeader(table.Row{"Field", "Value"})
	for field, value := range entity {
		if field == "activated" {
			value = time.Unix(int64(value.(float64))/1000, 0)
		}
		t.AppendRow(table.Row{field, value})
	}
	t.SetStyle(tableStyle)
	t.Style().Title.Align = text.AlignCenter
	t.Render()
}

// ActiveSnapsTable prints out the active snaps response as a table
func ActiveSnapsTable(response []byte) {
	var activeSnapsResponse ActiveSnapsResponse
	json.Unmarshal(response, &activeSnapsResponse)

	if activeSnapsResponse.Status != "success" {
		utils.PrintStatus(activeSnapsResponse.Status, activeSnapsResponse.Message)
		return
	}

	activeSnaps := activeSnapsResponse.Data

	// write out the table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Active Snaps")
	t.AppendHeader(table.Row{"Active Snap ID", "Snap ID", "State", "Activated", "Trigger", "Executions", "Errors"})
	for _, a := range activeSnaps {
		activated := time.Unix(a.Activated/1000, 0)
		t.AppendRow(table.Row{a.ActiveSnapID, a.SnapID, a.State, activated, a.Provider, a.ExecutionCounter, a.ErrorCounter})
	}
	t.SetStyle(tableStyle)
	t.Style().Title.Align = text.AlignCenter
	t.Render()
}

// do some common processing on the ActiveSnapLogResponse, and extract the logs
func extractLogs(response []byte) ([]ActiveSnapLog, error) {
	// unmarshal into the ActiveSnapLogsResponse struct, to flatten the property set
	var activeSnapLogsResponse ActiveSnapLogsResponse
	json.Unmarshal(response, &activeSnapLogsResponse)

	// check for errors
	utils.PrintStatus(activeSnapLogsResponse.Status, activeSnapLogsResponse.Message)
	if activeSnapLogsResponse.Status == "error" {
		return nil, errors.New("Error encountered")
	}

	// extract the active snap logs
	activeSnapLogs := activeSnapLogsResponse.Data

	// check for no rows
	if len(activeSnapLogs) < 1 {
		utils.PrintError("no logs found for this snap")
		return nil, errors.New("No logs found for this snap")
	}

	return activeSnapLogs, nil
}
