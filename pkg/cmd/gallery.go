package cmd

import (
	"os"

	"github.com/snapmaster-io/snap/pkg/api"
	"github.com/snapmaster-io/snap/pkg/print"
	"github.com/snapmaster-io/snap/pkg/utils"
	"github.com/spf13/cobra"
)

// galleryCmd represents the snaps command
var galleryCmd = &cobra.Command{
	Use:   "gallery [subcommand]",
	Short: "Interact with the gallery",
	Long:  `Interact with the gallery.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

// listGalleryCmd represents the list gallery subcommand
var listGalleryCmd = &cobra.Command{
	Use:   "list",
	Short: "List the snaps in the gallery",
	Long:  `List the snaps in the gallery.`,
	Run: func(cmd *cobra.Command, args []string) {

		// execute the API call
		response, err := api.Get("/gallery")
		if err != nil {
			utils.PrintErrorMessage("could not retrieve data", err)
			os.Exit(1)
		}

		format, err := rootCmd.PersistentFlags().GetString("format")
		if format == "json" {
			print.JSON(response)
			return
		}

		if format == "table" {
			print.SnapsTable(response)
			return
		}

		// unknown format - return the raw response
		print.RawResponse(response)
	},
}

func init() {
	rootCmd.AddCommand(galleryCmd)
	galleryCmd.AddCommand(listGalleryCmd)
	galleryCmd.AddCommand(getSnapCmd)
}
