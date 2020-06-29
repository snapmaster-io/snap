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

// configSetDevCmd represents the config set "dev" command
var configSetDevCmd = &cobra.Command{
	Use:   "dev",
	Short: "Set config information to dev environment",
	Long:  `Set config information to dev environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("ClientID", "f9BSuAhmF8dmUtJWZyjAVJbGJWQMKsMW")
		viper.Set("APIURL", "https://dev.snapmaster.io")
		viper.Set("AuthDomain", "snapmaster-dev.auth0.com")

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

// configSetProdCmd represents the config set "prod" command
var configSetProdCmd = &cobra.Command{
	Use:   "prod",
	Short: "Set config information to production environment",
	Long:  `Set config information to production environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("ClientID", "O4e0z2Ky5DSvjzw3N5YLgtrz1GGltkOb")
		viper.Set("APIURL", "https://www.snapmaster.io")
		viper.Set("AuthDomain", "snapmaster.auth0.com")

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
	configSetCmd.AddCommand(configSetDevCmd)
	configSetCmd.AddCommand(configSetProdCmd)

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
