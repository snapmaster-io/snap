package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to a SnapMaster deployment.",
	Long: `Login to a SnapMaster deployment.

If no server is specified, login to the public SnapMaster service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var username, password string
		var err error

		username, err = cmd.Flags().GetString("username")
		if err != nil {
			return err
		}

		password, err = cmd.Flags().GetString("password")
		if err != nil {
			return err
		}

		fmt.Printf("login called with username %s, password %s\n", username, password)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	loginCmd.Flags().StringP("username", "u", "", "Username")
	loginCmd.Flags().StringP("password", "p", "", "Password")
	loginCmd.Flags().BoolP("password-stdin", "", false, "Take the password from stdin")
}
