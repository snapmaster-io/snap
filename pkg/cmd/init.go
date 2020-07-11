package cmd

import (
	"fmt"
	"os"

	"github.com/snapmaster-io/snap/pkg/config"
	"github.com/snapmaster-io/snap/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the snap CLI environment",
	Long: `Initialize the snap CLI environment.

If no flags are specified, initializes the snap CLI to point to the public SnapMaster service.`,
	Run: func(cmd *cobra.Command, args []string) {
		// re-bind the flags at initCmd run time (they are typically bound to configSetCmd at cobra setup time)
		viper.BindPFlag("APIURL", cmd.Flags().Lookup("api-url"))
		viper.BindPFlag("ClientID", cmd.Flags().Lookup("client-id"))
		viper.BindPFlag("AuthDomain", cmd.Flags().Lookup("auth-domain"))

		var err error

		// Create the config file in case the path hasn't been created yet
		filename := viper.ConfigFileUsed()
		if len(filename) < 1 {
			filename, err = config.WriteConfigFile("config.json", []byte(""))
			if err != nil {
				utils.PrintError("could not write config file to $HOME/.config/snap/config.json")
				os.Exit(1)
			}
		}

		// use viper to write the config to the file
		err = viper.WriteConfig()
		if err != nil {
			utils.PrintErrorMessage(fmt.Sprintf("could not update config file %s", filename), err)
			os.Exit(1)
		} else {
			utils.PrintMessage(fmt.Sprintf("updated config file %s", filename))
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("api-url", "", "", "API URL (defaults to https://dev.snapmaster.io)")
	initCmd.Flags().StringP("client-id", "", "", "Auth0 Client ID (required for any non-default API URL)")
	initCmd.Flags().StringP("auth-domain", "", "", "Auth0 Auth Domain (defaults to snapmaster-dev.auth0.com)")
}
