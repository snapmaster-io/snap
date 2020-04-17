package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/table"
)

// ActiveSnap defines the fields to print for an activeSnap
type ActiveSnap struct {
	ActiveSnapID     string `json:"activeSnapId"`
	SnapID           string `json:"snapID"`
	State            string `json:"state"`
	Provider         string `json:"provider"`
	Activated        int64  `json:"activated"`
	ExecutionCounter int    `json:"executionCounter"`
	ErrorCounter     int    `json:"errorCounter"`
}

// ActiveSnapActionsLog defines the fields to print for action logs
type ActiveSnapActionsLog struct {
	Provider string                     `json:"provider"`
	State    string                     `json:"state"`
	Output   ActiveSnapActionsLogOutput `json:"output"`
}

// ActiveSnapActionsLogOutput defines the fields to print for action logs
type ActiveSnapActionsLogOutput struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Error  string `json:"error"`
}

// ActiveSnapLog defines the fields to print for an activeSnap's logs
type ActiveSnapLog struct {
	LogID        int64                  `json:"timestamp"`
	ActiveSnapID string                 `json:"activeSnapId"`
	SnapID       string                 `json:"snapID"`
	State        string                 `json:"state"`
	Trigger      string                 `json:"trigger"`
	Actions      []ActiveSnapActionsLog `json:"actions"`
}

// ActiveSnapStatus defines the fields to unmarshal from a pause/resume operation
type ActiveSnapStatus struct {
	Message    string     `json:"message"`
	ActiveSnap ActiveSnap `json:"activeSnap"`
}

// Snap defines the fields to print for a snap
type Snap struct {
	SnapID      string `json:"snapId"`
	Description string `json:"description"`
	Provider    string `json:"provider"`
	Private     bool   `json:"private"`
}

// SnapStatus defines the fields to unmarshal from a create/fork/publish/unpublish operation
type SnapStatus struct {
	Message string `json:"message"`
	Snap    Snap   `json:"snap"`
}

// what style to use for all tables
var tableStyle = table.StyleColoredCyanWhiteOnBlack
var actionTableStyle = table.StyleColoredBright

func printJSON(response []byte) {
	// pretty-print the json
	output := &bytes.Buffer{}
	err := json.Indent(output, response, "", "  ")
	if err != nil {
		fmt.Println("snap: could not format response as json")
		fmt.Println(string(response))
		os.Exit(1)
	}
	fmt.Println(output.String())
}

func printActiveSnap(response []byte) {
	// unmarshal into the ActiveSnap struct, to flatten the property set
	var activeSnap ActiveSnap
	json.Unmarshal(response, &activeSnap)

	// re-marshal and unmarshal into a map, which can be iterated over as a {name, value} pair
	intermediateEntity, _ := json.Marshal(activeSnap)
	var entity map[string]interface{}
	json.Unmarshal(intermediateEntity, &entity)

	// TODO: sort / alphabetize the keys

	// write out the table of properties
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Field", "Value"})
	for field, value := range entity {
		if field == "activated" {
			value = time.Unix(value.(int64)/1000, 0)
		}
		t.AppendRow(table.Row{field, value})
	}
	t.SetStyle(tableStyle)
	t.Render()
}

func printActiveSnapLogs(response []byte) {
	// unmarshal into the ActiveSnap struct, to flatten the property set
	var activeSnapLogs []ActiveSnapLog
	json.Unmarshal(response, &activeSnapLogs)

	// check for no rows
	if len(activeSnapLogs) < 1 {
		fmt.Println("snap: no logs found for this snap")
		return
	}

	// grab the SnapIP, ActiveSnapID, and Trigger from the first record
	activeSnapInstance := activeSnapLogs[0]
	activeSnapID := activeSnapInstance.ActiveSnapID
	snapID := activeSnapInstance.SnapID
	trigger := activeSnapInstance.Trigger

	// write out the table of properties
	t := table.NewWriter()
	t.SetTitle(fmt.Sprintf(
		"Logs for Snap ID %s\nActive Snap ID %s, triggered by %s",
		snapID, activeSnapID, trigger))
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Log ID", "Timestamp", "State"})
	for _, logEntry := range activeSnapLogs {
		timestamp := time.Unix(logEntry.LogID/1000, 0)
		t.AppendRow(table.Row{logEntry.LogID, timestamp, logEntry.State})
	}
	t.SetStyle(tableStyle)
	t.Render()
}

func printActiveSnapLogDetails(response []byte, logID string, format string) {
	// unmarshal into the ActiveSnap struct, to flatten the property set
	var activeSnapLogs []ActiveSnapLog
	json.Unmarshal(response, &activeSnapLogs)

	// check for no rows
	if len(activeSnapLogs) < 1 {
		fmt.Println("snap: no logs found for this snap")
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
		fmt.Printf("snap: log ID %s not found for active Snap ID %s\n", logID, activeSnapID)
		os.Exit(1)
	}

	// write out general information
	t := table.NewWriter()
	t.SetTitle("Action log details")
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Snap ID", "Active Snap ID", "Log ID"})
	t.AppendRow(table.Row{snapID, activeSnapID, logID})
	t.SetStyle(tableStyle)
	t.Render()

	fmt.Println("\nAction details:")

	// write out the table of properties
	for _, action := range logEntry.Actions {
		fmt.Println()
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Provider", "State"})
		t.AppendRow(table.Row{action.Provider, action.State})
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

func printActiveSnapStatus(response []byte) {
	// unmarshal into the ActiveSnapStatus struct, to get "Message" and
	// flatten the property set of the ActiveSnap
	var activeSnapStatus ActiveSnapStatus
	json.Unmarshal(response, &activeSnapStatus)

	fmt.Printf("snap: operation status: %s\n\n", activeSnapStatus.Message)
	activeSnap := activeSnapStatus.ActiveSnap

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
	t.Render()
}

func printActiveSnapsTable(response []byte) {
	//var activeSnaps []map[string]string
	var activeSnaps []ActiveSnap
	json.Unmarshal(response, &activeSnaps)

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
	t.Render()
}

func printSnapStatus(response []byte) {
	// unmarshal into the SnapStatus struct, to get "Message" and
	// flatten the property set of the Snap
	var snapStatus SnapStatus
	json.Unmarshal(response, &snapStatus)

	fmt.Printf("snap: operation status: %s\n\n", snapStatus.Message)
	if snapStatus.Message == "error" {
		return
	}

	// since the operation was successful, get the snap
	snap := snapStatus.Snap

	// re-marshal and unmarshal into a map, which can be iterated over as a {name, value} pair
	intermediateEntity, _ := json.Marshal(snap)
	var entity map[string]interface{}
	json.Unmarshal(intermediateEntity, &entity)

	// TODO: sort / alphabetize the keys

	// write out the table of properties
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Snap Values")
	t.AppendHeader(table.Row{"Field", "Value"})
	for field, value := range entity {
		t.AppendRow(table.Row{field, value})
	}
	t.SetStyle(tableStyle)
	t.Render()
}

func printSnapsTable(response []byte) {
	var snaps []map[string]string
	json.Unmarshal(response, &snaps)

	// write out the table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Snaps")
	t.AppendHeader(table.Row{"Snap ID", "Description", "Trigger"})
	for _, snap := range snaps {
		t.AppendRow(table.Row{snap["snapId"], snap["description"], snap["provider"]})
	}
	t.SetStyle(tableStyle)
	t.Render()
}

func printStatus(response []byte) {
	var status map[string]string
	json.Unmarshal(response, &status)

	// print the message field as the operation status
	fmt.Printf("snap: operation status: %s\n", string(status["message"]))
}
