package cmd

import (
	"fmt"

	"github.com/snapmaster-io/snap/pkg/auth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to a SnapMaster deployment.",
	Long: `Login to a SnapMaster deployment.

If no server is specified, login to the public SnapMaster service.`,
	Run: func(cmd *cobra.Command, args []string) {
		// hardcode clientId for now
		clientID := viper.GetString("ClientID")
		authDomain := viper.GetString("AuthDomain")
		redirectURL := viper.GetString("RedirectURL")

		auth.AuthorizeUser(clientID, authDomain, redirectURL)
	},
}

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of a SnapMaster deployment.",
	Long:  `Log out of a SnapMaster deployment.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("AccessToken", "")
		viper.Set("Name", "")
		viper.Set("Email", "")
		viper.WriteConfig()

		fmt.Println("snap: no logged in user.")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
}
