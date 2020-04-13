package cmd

import (
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

func init() {
	rootCmd.AddCommand(loginCmd)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// no need for command-line args anymore!  replaced with OAuth2 PKCE flow
	/*
		loginCmd.Flags().StringP("username", "u", "", "Username")
		loginCmd.Flags().StringP("password", "p", "", "Password")
		loginCmd.Flags().BoolP("password-stdin", "", false, "Take the password from stdin")
	*/
}
