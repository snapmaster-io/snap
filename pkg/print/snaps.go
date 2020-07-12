package print

import (
	"encoding/json"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/snapmaster-io/snap/pkg/utils"
)

// Snap defines the fields to unmarshal for a snap
type Snap struct {
	SnapID      string `json:"snapId"`
	Description string `json:"description"`
	Provider    string `json:"provider"`
	Private     bool   `json:"private"`
}

// SnapDefinition defines the text field to unmarshal for a snap's YAML definition
type SnapDefinition struct {
	Text string `json:"text"`
}

// SnapDefinitionResponse defines the fields to unmarshal for a SnapDefinition response
type SnapDefinitionResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    SnapDefinition `json:"data"`
}

// SnapResponse defines the fields to unmarshal from a create/fork/publish/unpublish operation
type SnapResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    Snap   `json:"data"`
}

// SnapsResponse defines the fields to unmarshal from a gallery list / snaps list operation
type SnapsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []Snap `json:"data"`
}

// SnapDefinitionYaml processes the snap response and prints out the definition as YAML
func SnapDefinitionYaml(response []byte) {
	// unmarshal into the SnapResponse struct, to get "Status" and
	// flatten the property set of the Snap
	var snapResponse SnapDefinitionResponse
	json.Unmarshal(response, &snapResponse)

	if snapResponse.Status == "error" {
		utils.PrintStatus(snapResponse.Status, snapResponse.Message)
		return
	}

	// since the operation was successful, get the snap
	snap := snapResponse.Data
	utils.PrintYAML(snap.Text)
}

// SnapStatusTable prints out the snap status response as a table
func SnapStatusTable(response []byte) {
	// unmarshal into the SnapResponse struct, to get "Status" and
	// flatten the property set of the Snap
	var snapResponse SnapResponse
	json.Unmarshal(response, &snapResponse)

	utils.PrintStatus(snapResponse.Status, snapResponse.Message)
	if snapResponse.Status == "error" {
		return
	}

	// since the operation was successful, get the snap
	snap := snapResponse.Data

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

// SnapsTable prints out the snaps in the response as a table
func SnapsTable(response []byte) {
	var snapsResponse SnapsResponse
	json.Unmarshal(response, &snapsResponse)

	if snapsResponse.Status == "error" {
		utils.PrintStatus(snapsResponse.Status, snapsResponse.Message)
		return
	}

	snaps := snapsResponse.Data

	// write out the table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Snaps")
	t.AppendHeader(table.Row{"Snap ID", "Description", "Trigger"})
	for _, snap := range snaps {
		//t.AppendRow(table.Row{snap["snapId"], snap["description"], snap["provider"]})
		t.AppendRow(table.Row{snap.SnapID, snap.Description, snap.Provider})
	}
	t.SetStyle(tableStyle)
	t.Style().Title.Align = text.AlignCenter
	t.Render()
}
