/*
Copyright © 2024 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/csl"
	"github.com/xeipuuv/gojsonschema"

	"github.com/front-matter/commonmeta/crossref"

	"github.com/front-matter/commonmeta/datacite"

	"github.com/spf13/cobra"
)

// sampleCmd represents the sample command
var sampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "A random sample of works",
	Long: `A random sample of works. Currently only available for
	the Crossref and DataCite provider. Options include numnber of samples, 
	work type, and Crossref member id or DataCite client id. For example:

	commonmeta sample --number 10 --member 78 --type journal-article,
	commonmeta sample --number 10 --member cern.zenodo --type dataset`,

	Run: func(cmd *cobra.Command, args []string) {
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

		var data []commonmeta.Data
		var err error
		sample := true
		if from == "crossref" {
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
			output, jsErr = commonmeta.WriteList(data)
		} else if to == "csl" {
			output, jsErr = csl.WriteList(data)
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
	rootCmd.AddCommand(sampleCmd)
}
