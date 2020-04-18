package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get and set config information",
	Long:  `Get and set config information.`,
	Run: func(cmd *cobra.Command, args []string) {
		printConfig()
	},
}

// configGetCmd represents the config get command
var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Print out config information",
	Long:  `Print out config information.`,
	Run: func(cmd *cobra.Command, args []string) {
		printConfig()
	},
}

// configSetCmd represents the config set command
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set config information",
	Long:  `Set config information based on the flags provided.`,
	Run: func(cmd *cobra.Command, args []string) {
		// use viper to write the config to the file
		err := viper.WriteConfig()
		if err != nil {
			fmt.Printf("snap: could not write config file\nerror: %s\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("snap: updated config\n")
		}

		printConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)

	configSetCmd.Flags().StringP("api-url", "", "", "API URL (defaults to https://dev.snapmaster.io)")
	configSetCmd.Flags().StringP("client-id", "", "", "Auth0 Client ID (required for any non-default API URL)")
	configSetCmd.Flags().StringP("auth-domain", "", "", "Auth0 Auth Domain (defaults to snapmaster-dev.auth0.com)")

	viper.BindPFlag("APIURL", configSetCmd.Flags().Lookup("api-url"))
	viper.BindPFlag("ClientID", configSetCmd.Flags().Lookup("client-id"))
	viper.BindPFlag("AuthDomain", configSetCmd.Flags().Lookup("auth-domain"))
}

func printConfig() {
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
