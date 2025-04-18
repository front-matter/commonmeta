/*
Copyright © 2024 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/front-matter/commonmeta/ror"
	"github.com/spf13/cobra"
)

// matchCmd represents the match command
var matchCmd = &cobra.Command{
	Use:   "match",
	Short: "Match a string to an identifier.",
	Long: `Match a string to an identifier. Supports affiliation
  matching for ROR.
	
	Example usage:
	
	commonmeta match "Leibniz Universität Hannover"`,
	Run: func(cmd *cobra.Command, args []string) {
		var input string
		var orgdata ror.ROR
		var output []byte
		var err error

		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")

		if len(args) == 0 {
			fmt.Println("Please provide an input")
			return
		}

		input = args[0]

		if from == "ror" {
			orgdata, err = ror.MatchAffiliation(input)
			if err != nil {
				cmd.Println(err)
				return
			}
		}

		if orgdata.ID == "" {
			cmd.Println("No match found")
			return
		} else if to == "inveniordm" {
			output, err = ror.WriteInvenioRDM(orgdata)
			cmd.Println(string(output))
		} else if to == "ror" {
			output, err = ror.Write(orgdata)
			var out bytes.Buffer
			json.Indent(&out, output, "", "  ")
			cmd.Println(out.String())
		}

		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(matchCmd)
}
