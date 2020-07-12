package print

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/viper"
)

// Config prints out the current configuration as a table
func Config() {
	configMap := map[string]string{
		"API URL":     viper.GetString("APIURL"),
		"Client ID":   viper.GetString("ClientID"),
		"Auth Domain": viper.GetString("AuthDomain"),
	}

	// write out the table of properties
	t := table.NewWriter()
	t.SetTitle("Config Values")
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Field", "Value"})
	for field, value := range configMap {
		t.AppendRow(table.Row{field, value})
	}
	t.SetStyle(tableStyle)
	t.Style().Title.Align = text.AlignCenter
	t.Render()
}
