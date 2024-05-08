/*
Copyright © 2024 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/csl"
	"github.com/xeipuuv/gojsonschema"

	"github.com/front-matter/commonmeta/crossref"
	"github.com/front-matter/commonmeta/crossrefxml"

	"github.com/front-matter/commonmeta/datacite"

	"github.com/front-matter/commonmeta/schemaorg"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A list of works",
	Long: `A list of works. Currently only available for
	the Crossref and DataCite provider. Options include numnber of works, 
	work type, and Crossref member id or DataCite client id. For example:

	commonmeta list --number 10 --member 78 --type journal-article,
	commonmeta list --number 10 --member cern.zenodo --type dataset`,
	Run: func(cmd *cobra.Command, args []string) {
		var input string
		var str string // a string, content loaded from a file
		var err error
		var data []commonmeta.Data

		if len(args) > 0 {
			input = args[0]
		}
		number, _ := cmd.Flags().GetInt("number")
		from, _ := cmd.Flags().GetString("from")

		member, _ := cmd.Flags().GetString("member")
		type_, _ := cmd.Flags().GetString("type")
		hasORCID, _ := cmd.Flags().GetBool("has-orcid")
		hasROR, _ := cmd.Flags().GetBool("has-ror-id")
		hasReferences, _ := cmd.Flags().GetBool("has-references")
		hasRelation, _ := cmd.Flags().GetBool("has-relation")
		hasAbstract, _ := cmd.Flags().GetBool("has-abstract")
		hasAward, _ := cmd.Flags().GetBool("has-award")
		hasLicense, _ := cmd.Flags().GetBool("has-license")
		hasArchive, _ := cmd.Flags().GetBool("has-archive")
		sample := false

		if input != "" {
			_, err = os.Stat(input)
			if err != nil {
				fmt.Printf("File not found: %s", input)
				return
			}
			str = input
		}

		if str != "" && from == "crossref" {
			data, err = crossref.LoadList(str)
		} else if str != "" && from == "crossrefxml" {
			data, err = crossrefxml.LoadList(str)
		} else if str != "" && from == "datacite" {
			data, err = datacite.LoadList(str)
		} else if from == "crossref" {
			data, err = crossref.FetchList(number, member, type_, sample, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense, hasArchive)
		} else if from == "datacite" {
			data, err = datacite.FetchList(number, sample)
		}
		if err != nil {
			fmt.Println(err)
		}

		var output []byte
		var jsErr []gojsonschema.ResultError
		to, _ := cmd.Flags().GetString("to")
		if to == "commonmeta" {
			output, jsErr = commonmeta.WriteAll(data)
		} else if to == "csl" {
			output, jsErr = csl.WriteAll(data)
		} else if to == "datacite" {
			output, jsErr = datacite.WriteAll(data)
		} else if to == "schemaorg" {
			output, jsErr = schemaorg.WriteAll(data)
		}

		if err != nil {
			fmt.Println(err)
		}
		var out bytes.Buffer
		json.Indent(&out, output, "", "  ")
		fmt.Println(out.String())

		if jsErr != nil {
			fmt.Println(jsErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
