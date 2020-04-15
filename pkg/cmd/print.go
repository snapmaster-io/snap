package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"bytes"

	"github.com/jedib0t/go-pretty/table"
)

// ActiveSnap defines the fields to print for an activeSnap
type ActiveSnap struct {
	ActiveSnapID string `json:"activeSnapId"`
	SnapID string `json:"snapID"`
	State string `json:"state"`
	Provider string `json:"provider"`
	Activated bool `json:"activated"`
	ExecutionCounter int `json:"executionCounter"`
}

// ActiveSnapStatus defines the fields to unmarshal from a pause/resume operation
type ActiveSnapStatus struct {
	Message string `json:"message"`
	ActiveSnap ActiveSnap `json:"activeSnap"`
}

// Snap defines the fields to print for a snap
type Snap struct {
	SnapID string `json:"snapId"`
	Description string `json:"description"`
	Provider string `json:"provider"`
	Private bool `json:"private"`
}

// SnapStatus defines the fields to unmarshal from a create/fork/publish/unpublish operation
type SnapStatus struct {
	Message string `json:"message"`
	Snap Snap `json:"snap"`
}

// what style to use for all tables
var tableStyle = table.StyleColoredCyanWhiteOnBlack

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
		t.AppendRow(table.Row{field, value})
	}
	t.SetStyle(tableStyle)
	t.Render()
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
		t.AppendRow(table.Row{field, value})
	}
	t.SetStyle(tableStyle)
	t.Render()
}

func printActiveSnapsTable(response []byte) {
	var snaps []map[string]string
	json.Unmarshal(response, &snaps)

	// write out the table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Active Snaps")
	t.AppendHeader(table.Row{"Active Snap ID", "Snap ID", "State", "Trigger", "Executions", "Errors"})
	for _, snap := range snaps {
		t.AppendRow(table.Row{snap["activeSnapId"], snap["snapId"], snap["state"], snap["trigger"], snap["executionCounter"], snap["errorCounter"]})
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
		t.AppendRow(table.Row{snap["snapId"], snap["description"], snap["trigger"]})
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
