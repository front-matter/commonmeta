/*
Copyright © 2024 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/front-matter/commonmeta/commonmeta"

	"github.com/front-matter/commonmeta/crossref"

	"github.com/front-matter/commonmeta/datacite"

	"github.com/front-matter/commonmeta/types"

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

		var data []types.Data
		var err error
		sample := false
		if from == "crossref" {
			data, err = crossref.FetchCrossrefList(number, member, type_, sample, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense, hasArchive)
		} else if from == "datacite" {
			data, err = datacite.FetchDataciteList(number, sample)
		}
		if err != nil {
			fmt.Println(err)
		}
		output, jsErr := commonmeta.WriteCommonmetaList(data)
		var out bytes.Buffer
		json.Indent(&out, output, "=", "\t")
		fmt.Println(out.String())

		if jsErr != nil {
			fmt.Println(jsErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
