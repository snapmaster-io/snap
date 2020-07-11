package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/snapmaster-io/snap/pkg/utils"
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
	Action   string                     `json:"action"`
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
	Event        string                 `json:"event"`
	Actions      []ActiveSnapActionsLog `json:"actions"`
}

// ActiveSnapLogsStatus defines the fields to unmarshal from a pause/resume operation
type ActiveSnapLogsStatus struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    []ActiveSnapLog `json:"data"`
}

// ActiveSnapsStatus defines the fields to unmarshal from getting all active snaps
type ActiveSnapsStatus struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    []ActiveSnap `json:"data"`
}

// ActiveSnapStatus defines the fields to unmarshal from a pause/resume operation
type ActiveSnapStatus struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Data    ActiveSnap `json:"data"`
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
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    Snap   `json:"data"`
}

// SnapsStatus defines the fields to unmarshal from a gallery list / snaps list operation
type SnapsStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []Snap `json:"data"`
}

// what style to use for all tables
var tableStyle = table.StyleColoredCyanWhiteOnBlack
var actionTableStyle = table.StyleColoredBright

func printJSON(response []byte) {
	// pretty-print the json
	utils.PrintJSON(response)
}

func printJSONString(response string) {
	// pretty-print the json
	utils.PrintJSON([]byte(response))
}

func printRawResponse(response []byte) {
	printRawResponse(response)
}

func printActiveSnap(response []byte) {
	// unmarshal into the ActiveSnap struct, to flatten the property set
	var activeSnapStatus ActiveSnapStatus
	json.Unmarshal(response, &activeSnapStatus)

	if activeSnapStatus.Status != "success" {
		utils.PrintStatus(activeSnapStatus.Status, activeSnapStatus.Message)
		return
	}

	activeSnap := activeSnapStatus.Data

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

func printActiveSnapLogs(response []byte) {
	// unmarshal into the ActiveSnapLogsStatus struct, to flatten the property set
	var activeSnapLogs ActiveSnapLogsStatus
	json.Unmarshal(response, &activeSnapLogs)

	// check for errors
	utils.PrintStatus(activeSnapLogs.Status, activeSnapLogs.Message)
	if activeSnapLogs.Status == "error" {
		return
	}

	// check for no rows
	if len(activeSnapLogs.Data) < 1 {
		utils.PrintError("no logs found for this snap")
		return
	}

	// grab the SnapIP, ActiveSnapID, and Trigger from the first record
	activeSnapInstance := activeSnapLogs.Data[0]
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
	for _, logEntry := range activeSnapLogs.Data {
		timestamp := time.Unix(logEntry.LogID/1000, 0)
		t.AppendRow(table.Row{logEntry.LogID, timestamp, logEntry.State})
	}
	t.SetStyle(tableStyle)
	t.Style().Title.Align = text.AlignCenter
	t.Render()
}

func printActiveSnapLogDetails(response []byte, logID string, format string) {
	// unmarshal into the ActiveSnapLogsStatus struct, to flatten the property set
	var activeSnapLogs ActiveSnapLogsStatus
	json.Unmarshal(response, &activeSnapLogs)

	// check for errors
	utils.PrintStatus(activeSnapLogs.Status, activeSnapLogs.Message)
	if activeSnapLogs.Status == "error" {
		return
	}

	// check for no rows
	if len(activeSnapLogs.Data) < 1 {
		utils.PrintError("no logs found for this snap")
		return
	}

	// grab the SnapIP, ActiveSnapID, and Trigger from the first record
	activeSnapInstance := activeSnapLogs.Data[0]
	activeSnapID := activeSnapInstance.ActiveSnapID
	snapID := activeSnapInstance.SnapID

	var logEntry ActiveSnapLog

	// find the entry with the right logID
	found := false
	for k, v := range activeSnapLogs.Data {
		logIDasInt64, _ := strconv.ParseInt(logID, 10, 64)
		if v.LogID == logIDasInt64 {
			logEntry = activeSnapLogs.Data[k]
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

func printActiveSnapStatus(response []byte) {
	// unmarshal into the ActiveSnapStatus struct, to get "Status" and
	// flatten the property set of the ActiveSnap
	var activeSnapStatus ActiveSnapStatus
	json.Unmarshal(response, &activeSnapStatus)

	utils.PrintStatus(activeSnapStatus.Status, activeSnapStatus.Message)

	// if the status indicates an error, there is no active snap to display
	if activeSnapStatus.Status != "success" {
		return
	}

	activeSnap := activeSnapStatus.Data

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

func printActiveSnapsTable(response []byte) {
	var activeSnapsStatus ActiveSnapsStatus
	json.Unmarshal(response, &activeSnapsStatus)

	if activeSnapsStatus.Status != "success" {
		utils.PrintStatus(activeSnapsStatus.Status, activeSnapsStatus.Message)
		return
	}

	activeSnaps := activeSnapsStatus.Data

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

func printConnectionsTable(response []byte) {
	var tools []map[string]string
	json.Unmarshal(response, &tools)

	// write out the table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Provider"})
	for _, tool := range tools {
		connected := tool["connected"] != ""
		if connected {
			t.AppendRow(table.Row{tool["provider"]})
		}
	}
	t.SetStyle(tableStyle)
	t.Style().Title.Align = text.AlignCenter
	t.Render()
}

func printCredentialsTable(response []byte, connection string) {
	var credentials []map[string]string
	json.Unmarshal(response, &credentials)

	// write out the table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle(fmt.Sprintf("Credential sets for %s connection", connection))
	t.AppendHeader(table.Row{"Credential set name"})
	for _, credential := range credentials {
		t.AppendRow(table.Row{credential["__id"]})
	}
	t.SetStyle(tableStyle)
	t.Style().Title.Align = text.AlignCenter
	t.Render()
}

func printSnapStatus(response []byte) {
	// unmarshal into the SnapStatus struct, to get "Status" and
	// flatten the property set of the Snap
	var snapStatus SnapStatus
	json.Unmarshal(response, &snapStatus)

	utils.PrintStatus(snapStatus.Status, snapStatus.Message)
	if snapStatus.Status == "error" {
		return
	}

	// since the operation was successful, get the snap
	snap := snapStatus.Data

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
	t.Style().Title.Align = text.AlignCenter
	t.Render()
}

func printSnapsTable(response []byte) {
	var snaps SnapsStatus
	json.Unmarshal(response, &snaps)

	// write out the table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Snaps")
	t.AppendHeader(table.Row{"Snap ID", "Description", "Trigger"})
	for _, snap := range snaps.Data {
		//t.AppendRow(table.Row{snap["snapId"], snap["description"], snap["provider"]})
		t.AppendRow(table.Row{snap.SnapID, snap.Description, snap.Provider})
	}
	t.SetStyle(tableStyle)
	t.Style().Title.Align = text.AlignCenter
	t.Render()
}

func printStatus(response []byte) {
	// unmarshal into the SnapStatus struct, to get "Status" and
	// flatten the property set of the Snap
	var snapStatus SnapStatus
	json.Unmarshal(response, &snapStatus)

	utils.PrintStatus(snapStatus.Status, snapStatus.Message)
}

func printToolsTable(response []byte) {
	var tools []map[string]string
	json.Unmarshal(response, &tools)

	// write out the table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Tools Library")
	t.AppendHeader(table.Row{"Provider", "Type", "Connected?"})
	for _, tool := range tools {
		connected := tool["connected"] != ""
		t.AppendRow(table.Row{tool["provider"], tool["type"], connected})
	}
	t.SetStyle(tableStyle)
	t.Style().Title.Align = text.AlignCenter
	t.Render()
}
