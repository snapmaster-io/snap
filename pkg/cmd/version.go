package cmd

import (
	"fmt"

	"github.com/snapmaster-io/snap/pkg/utils"
	"github.com/snapmaster-io/snap/pkg/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the user command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current version",
	Long: `Show the current version.

NOTE: snap login must be called before there is an active user.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintMessage(fmt.Sprintf("version <%s>, git hash <%s>", version.Version, version.GitHash))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
