/*
Copyright © 2024-2025 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"fmt"

	"github.com/front-matter/commonmeta/ror"
	"github.com/front-matter/commonmeta/spdx"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a vocabulary",
	Long: `Install a vocabulary. Example usage:
	
	commonmeta install spdx`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var fileName string

		dataVersion, _ := cmd.Flags().GetString("data-version")

		if len(args) == 0 {
			fmt.Println("Please provide an input")
			return
		}
		input := args[0]

		switch input {
		case "spdx":
			_, err = spdx.FetchAll()
			if err != nil {
				cmd.Println(err)
				return
			}
			fileName = spdx.SPDXFilename
		case "ror":
			if dataVersion == "" {
				dataVersion = ror.DefaultVersion
			}
			_, err = ror.FetchAll(dataVersion)
			if err != nil {
				cmd.Println(err)
				return
			}
			fileName = ror.ArchivedFilename
		default:
			cmd.Println("Unsupported vocabulary. Supported vocabularies are: spdx, ror")
			return
		}
		cmd.Printf("Saved %s file: %s\n", input, fileName)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
