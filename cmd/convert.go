/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/front-matter/commonmeta-go/datacite"
	"github.com/front-matter/commonmeta-go/doiutils"
	"github.com/front-matter/commonmeta-go/types"

	"github.com/front-matter/commonmeta-go/crossref"

	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert scholarly metadata from one format to another",
	Long: `Convert scholarly metadata between formats. Currently
supported input formats are Crossref and DataCite DOIs, currently
the only supported output format is Commonmeta. Example usage:

commonmeta 10.5555/12345678`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide an input DOI")
			return
		}
		input := args[0]
		from, _ := cmd.Flags().GetString("from")
		var data types.Data
		var err error
		if from == "" {
			var ok bool
			doi, ok := doiutils.ValidateDOI(input)
			if !ok {
				fmt.Println("Please provide a valid DOI from Crossref or Datacite")
				return
			}
			from, ok = doiutils.GetDOIRA(doi)
			if !ok {
				fmt.Println("Please provide a valid DOI from Crossref or Datacite")
				return
			}
			from = strings.ToLower(from)
		}
		if from == "crossref" {
			data, err = crossref.FetchCrossref(input)
		} else if from == "datacite" {
			data, err = datacite.FetchDatacite(input)
		}

		if err != nil {
			fmt.Println(err)
		}
		output, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(output))
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.PersistentFlags().StringP("from", "f", "", "the format to convert from")
	convertCmd.PersistentFlags().StringP("to", "t", "commonmeta", "the format to convert to")
}
