package print

import (
	"encoding/json"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/snapmaster-io/snap/pkg/utils"
)

// Tool defines the fields to unmarshal for a tool
type Tool struct {
	Provider  string `json:"provider"`
	Type      string `json:"type"`
	Connected string `json:"connected"`
}

// ToolsResponse defines the fields to unmarshal from a get tools operation
type ToolsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []Tool `json:"data"`
}

// ToolsTable prints out the tools and their type and connection status as a table
func ToolsTable(response []byte) {
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
	t.SetTitle("Tools Library")
	t.AppendHeader(table.Row{"Provider", "Type", "Connected?"})
	for _, tool := range tools {
		connected := tool.Connected != ""
		t.AppendRow(table.Row{tool.Provider, tool.Type, connected})
	}
	t.SetStyle(tableStyle)
	t.Style().Title.Align = text.AlignCenter
	t.Render()
}
