/*
Copyright © 2025 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/front-matter/commonmeta/ror"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import vocabulary",
	Long: `Import a vocabulary from a file. Currently only supported for
	ROR (Research Organization Registry) and InvenioRDM awards vocabulary.
	
	Example usage:
	
	commonmeta import v1.63-2025-04-03-ror-data_schema_v2.json -f ror -t inveniordm`,
	Run: func(cmd *cobra.Command, args []string) {
		var input string
		var str string // a string, content loaded from a file
		var err error
		var data []ror.ROR
		if len(args) > 0 {
			input = args[0]
		}
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")

		cmd.SetOut(os.Stdout)
		cmd.SetErr(os.Stderr)

		if input != "" && !strings.HasPrefix(input, "--") {
			_, err = os.Stat(input)
			if err != nil {
				cmd.PrintErrf("File not found: %s", input)
				return
			}
			str = input
		}

		if from == "ror" {
			data, err = ror.LoadAll(str)
		} else {
			fmt.Println("Please provide a valid input format")
			return
		}
		if err != nil {
			fmt.Println("An error occurred:", err)
			return
		}

		var filename string
		if to == "inveniordm" {
			filename, err = ror.WriteAll(data, input)
		}
		fmt.Println("File written:", filename)

		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
