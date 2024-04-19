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
	the Crossref API. Options include numnber of samples, DOI prefix
	and work type. For example:

	commonmeta sample --number 10 --prefix 10.7554 --type journal-article`,
	Run: func(cmd *cobra.Command, args []string) {
		number, _ := cmd.Flags().GetInt("number")
		prefix, _ := cmd.Flags().GetString("prefix")
		type_, _ := cmd.Flags().GetString("type")

		data, err := crossref.FetchCrossrefSample(number, prefix, type_)
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

	rootCmd.PersistentFlags().StringP("number", "n", "10", "number of samples")
	rootCmd.PersistentFlags().StringP("prefix", "", "", "DOI prefix")
	rootCmd.PersistentFlags().StringP("type", "", "journal-article", "work type")
}
