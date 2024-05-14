/*
Copyright © 2024 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"fmt"

	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/utils"
	"github.com/spf13/cobra"
)

// encodeCmd represents the encode command
var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Generate a random DOI string given a prefix",
	Long:  `Generate a random DOI string given a prefix. For example: commonmeta encode 10.5555`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide an input")
			return
		}
		input := args[0]
		prefix, ok := doiutils.ValidatePrefix(input)
		if !ok {
			cmd.PrintErr("Invalid prefix")
			return
		}
		doi := utils.EncodeDOI(prefix)
		cmd.Println(doi)
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)
}
