package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/snapmaster-io/snap/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the login command
var initCmd = &cobra.Command{
	Use:   "init [API server URL] [Client ID] [Auth Domain]",
	Short: "Initialize the snap CLI environment",
	Long: `Initialize the snap CLI environment.

If no arguments are specified, initializes the snap CLI to the public SnapMaster service.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			viper.Set("APIURL", args[0])
		}
		if len(args) > 1 {
			viper.Set("ClientID", args[1])
		}
		if len(args) > 2 {
			viper.Set("AuthDomain", args[2])
		}

		// Create the config file in case the path hasn't been created yet
		filename, err := config.WriteConfigFile("config.json", []byte(""))
		if err != nil {
			log.Fatal("could not write config file")
		}

		// use viper to write the config to the file
		err = viper.WriteConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Printf("snap: created config file in %s\n", filename)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")
	// initCmd.Flags().StringP("bar", "b", "", "Bar")
}
