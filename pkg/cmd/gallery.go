package cmd

import (
	"fmt"
	"os"

	"github.com/snapmaster-io/snap/pkg/api"
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
			fmt.Printf("snap: could not retrieve data: %s", err)
			os.Exit(1)
		}

		format, err := rootCmd.PersistentFlags().GetString("format")
		if format == "json" {
			printJSON(response)
			return
		}

		if format == "table" {
			printSnapsTable(response)
			return
		}

		// unknown format - return the raw response
		fmt.Printf("Raw response:\n%s\n", string(response))
	},
}

func init() {
	rootCmd.AddCommand(galleryCmd)
	galleryCmd.AddCommand(listGalleryCmd)
	galleryCmd.AddCommand(getSnapCmd)
}
