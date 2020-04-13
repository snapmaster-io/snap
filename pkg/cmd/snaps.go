package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// snapsCmd represents the snaps command
var snapsCmd = &cobra.Command{
	Use:   "snaps [subcommand]",
	Short: "Manage the user's snaps",
	Long:  `Manage the user's snaps.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

// listSnapsCmd represents the get snaps subcommand
var listSnapsCmd = &cobra.Command{
	Use:   "list",
	Short: "List the user's snaps",
	Long:  `List the user's snaps.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("snap: list!")
	},
}

func init() {
	rootCmd.AddCommand(snapsCmd)
	snapsCmd.AddCommand(listSnapsCmd)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")
	// loginCmd.Flags().StringP("username", "u", "", "Username")
}
