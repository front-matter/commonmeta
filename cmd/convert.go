/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/csl"
	"github.com/front-matter/commonmeta/datacite"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/schemaorg"
	"github.com/front-matter/commonmeta/utils"
	"github.com/xeipuuv/gojsonschema"

	"github.com/front-matter/commonmeta/crossref"

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
		var id string  // an identifier, content fetched via API
		var str string // a string, content loaded from a file
		var err error
		var data commonmeta.Data

		if len(args) == 0 {
			fmt.Println("Please provide an input")
			return
		}
		input := args[0]
		id = utils.NormalizeID(input)
		if id == "" {
			_, err = os.Stat(input)
			if err != nil {
				fmt.Printf("File not found: %s", input)
				return
			}
			str = input
		}

		from, _ := cmd.Flags().GetString("from")
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

		if id != "" {
			if from == "crossref" {
				data, err = crossref.Fetch(id)
			} else if from == "datacite" {
				data, err = datacite.Fetch(id)
			} else {
				fmt.Println("Please provide a valid input")
				return
			}
		} else if str != "" {
			if from == "commonmeta" {
				data, err = commonmeta.Load(str)
			} else if from == "crossref" {
				data, err = crossref.Load(str)
			} else if from == "datacite" {
				data, err = datacite.Load(str)
			} else {
				fmt.Println("Please provide a valid input")
				return
			}
		}

		var output []byte
		var jsErr []gojsonschema.ResultError
		to, _ := cmd.Flags().GetString("to")
		if to == "commonmeta" {
			output, jsErr = commonmeta.Write(data)
		} else if to == "csl" {
			output, jsErr = csl.Write(data)
		} else if to == "datacite" {
			output, jsErr = datacite.Write(data)
		} else if to == "schemaorg" {
			output, jsErr = schemaorg.Write(data)
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
	rootCmd.AddCommand(convertCmd)

	convertCmd.PersistentFlags().StringP("from", "f", "", "the format to convert from")
	convertCmd.PersistentFlags().StringP("to", "t", "commonmeta", "the format to convert to")
}
