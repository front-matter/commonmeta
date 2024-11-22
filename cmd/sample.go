/*
Copyright © 2024 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/crossrefxml"
	"github.com/front-matter/commonmeta/csl"
	"github.com/front-matter/commonmeta/inveniordm"
	"github.com/front-matter/commonmeta/schemaorg"
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

		client_, _ := cmd.Flags().GetString("client")
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

		depositor, _ := cmd.Flags().GetString("depositor")
		email, _ := cmd.Flags().GetString("email")
		registrant, _ := cmd.Flags().GetString("registrant")

		var data []commonmeta.Data
		var err error
		sample := true
		if from == "crossref" {
			data, err = crossref.FetchAll(number, member, type_, sample, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense, hasArchive)
		} else if from == "datacite" {
			data, err = datacite.FetchAll(number, client_, type_, sample)
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
		} else if to == "crossrefxml" {
			account := crossrefxml.Account{
				Depositor:  depositor,
				Email:      email,
				Registrant: registrant,
			}
			output, jsErr = crossrefxml.WriteAll(data, account)
		} else if to == "schemaorg" {
			output, jsErr = schemaorg.WriteAll(data)
		} else if to == "inveniordm" {
			output, jsErr = inveniordm.WriteAll(data)
		}

		if to == "crossrefxml" {
			fmt.Printf("%s\n", output)
		} else {
			var out bytes.Buffer
			json.Indent(&out, output, "", "  ")
			fmt.Println(out.String())
		}

		if jsErr != nil {
			cmd.PrintErr(jsErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(sampleCmd)
}
