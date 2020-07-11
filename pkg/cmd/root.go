package cmd

import (
	"fmt"
	"os"

	"github.com/snapmaster-io/snap/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "snap",
	Short: "SnapMaster CLI",
	Long: `
SnapMaster is a tool that manages and runs snaps.  Snaps are workflows which tie 
various dev and operational tools together.  Snaps define a trigger (an event such 
as a webhook) and a set of actions (anything that can be executed over a REST API).

snap is the SnapMaster CLI.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/snap/config.json)")
	rootCmd.PersistentFlags().StringP("format", "f", "table", "return output of command as one of {table, json}")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// set some default values
	viper.SetDefault("ClientID", "O4e0z2Ky5DSvjzw3N5YLgtrz1GGltkOb")
	viper.SetDefault("APIURL", "https://www.snapmaster.io")
	viper.SetDefault("AuthDomain", "snapmaster.auth0.com")
	viper.SetDefault("RedirectURL", "http://localhost:8085")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in $HOME/.config/snap/config
		viper.AddConfigPath(fmt.Sprintf("%s/.config/snap", home))
		viper.SetConfigName("config.json")
		viper.SetConfigType("json")
	}

	viper.SetEnvPrefix("snap")
	viper.AutomaticEnv() // read in environment variables that match, prefixed with "SNAP_"

	// read in a config file, if found
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
		} else {
			// Config file was found but another error was produced
			utils.PrintErrorMessage("config file was found, but another error occurred", err)
		}
	} else {
		// do not report non-error condition
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
