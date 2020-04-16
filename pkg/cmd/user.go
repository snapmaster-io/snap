package cmd

import (
	"fmt"
	"os"

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
			fmt.Println("snap: no logged in user.  To login, use the command 'snap login'.")
			os.Exit(1)
		}

		fmt.Printf("snap: current user is %s <%s>\n", name, email)
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
}
