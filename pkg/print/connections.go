package print

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/snapmaster-io/snap/pkg/utils"
)

// ConnectionsTable prints out the connected providers as a table
func ConnectionsTable(response []byte) {
	var toolsResponse ToolsResponse
	json.Unmarshal(response, &toolsResponse)

	if toolsResponse.Status == "error" {
		utils.PrintStatus(toolsResponse.Status, toolsResponse.Message)
		return
	}

	tools := toolsResponse.Data

	// write out the table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Provider"})
	for _, tool := range tools {
		connected := tool.Connected != ""
		if connected {
			t.AppendRow(table.Row{tool.Provider})
		}
	}
	t.SetStyle(tableStyle)
	t.Style().Title.Align = text.AlignCenter
	t.Render()
}

// CredentialsTable prints out the credentials of a connection as a table
func CredentialsTable(response []byte, connection string) {
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
