/*
Copyright © 2024 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"bytes"
	"commonmeta/commonmeta"
	"commonmeta/crossref"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// sampleCmd represents the sample command
var sampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "A random sample of works",
	Long: `A random sample of works. Currently only available for
	the Crossref API. Options include numnber of samples, Crossref
	member id and work type. For example:

	commonmeta sample --number 10 --member 78 --type journal-article`,
	Run: func(cmd *cobra.Command, args []string) {
		number, _ := cmd.Flags().GetInt("number")
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

		data, err := crossref.FetchCrossrefSample(number, member, type_, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense, hasArchive)
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
	rootCmd.AddCommand(sampleCmd)
}
