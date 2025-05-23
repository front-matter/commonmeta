/*
Copyright © 2024-2025 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"fmt"

	"github.com/front-matter/commonmeta/utils"
	"github.com/spf13/cobra"
)

// decodeCmd represents the decode command
var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode an identifier.",
	Long: `Decode a DOI, ROR or ORCID identifier. For DOIs only Crockford
	base32-encoding is supported, used by Rogue Scholar and some DataCite
	members.

	Example usage:

	commonmeta decode 10.54900/d3ck1-skq19`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide an input")
			return
		}
		number, err := utils.DecodeID(args[0])
		if err != nil {
			cmd.Println(err)
			return
		}
		cmd.Println(number)
	},
}

func init() {
	rootCmd.AddCommand(decodeCmd)
}
