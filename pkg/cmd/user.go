package cmd

import (
	"fmt"
	"os"

	"github.com/snapmaster-io/snap/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Show the active user",
	Long: `Show the active user.

NOTE: snap login must be called before there is an active user.`,
	Run: func(cmd *cobra.Command, args []string) {
		accessToken := viper.GetString("AccessToken")
		name := viper.GetString("Name")
		email := viper.GetString("Email")
		if accessToken == "" {
			utils.PrintError("no logged in user.  To login, use the command 'snap login'.")
			os.Exit(1)
		}

		utils.PrintMessage(fmt.Sprintf("current user is %s <%s>", name, email))
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
}
