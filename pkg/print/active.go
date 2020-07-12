package print

import (
	"encoding/json"
	"fmt"
	"os"
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
