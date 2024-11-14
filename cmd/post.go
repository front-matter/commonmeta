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
	"github.com/front-matter/commonmeta/crossrefxml"
	"github.com/front-matter/commonmeta/datacite"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/inveniordm"
	"github.com/front-matter/commonmeta/jsonfeed"
	"github.com/front-matter/commonmeta/utils"

	"github.com/front-matter/commonmeta/crossref"

	"github.com/spf13/cobra"
)

var postCmd = &cobra.Command{
	Use:   "post",
	Short: "Post scholarly metadata with a service",
	Long: `Post scholarly metadata to a service. Currently
the only supported output service is InvenioRDM. Example usage:

commonmeta post 10.5555/12345678 -f crossref -t inveniordm -h rogue-scholar.org --token mytoken`,

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
				cmd.PrintErr("Please provide a valid DOI from Crossref or Datacite")
				return
			}
			from, ok = doiutils.GetDOIRA(doi)
			if !ok {
				cmd.PrintErr("Please provide a valid DOI from Crossref or Datacite")
				return
			}
			from = strings.ToLower(from)
		}

		if id != "" {
			if from == "crossref" {
				data, err = crossref.Fetch(id)
			} else if from == "crossrefxml" {
				data, err = crossrefxml.Fetch(id)
			} else if from == "datacite" {
				data, err = datacite.Fetch(id)
			} else if from == "jsonfeed" {
				data, err = jsonfeed.Fetch(id)
			} else {
				fmt.Println("Please provide a valid input")
				return
			}
		} else if str != "" {
			if from == "commonmeta" {
				data, err = commonmeta.Load(str)
			} else if from == "crossref" {
				data, err = crossref.Load(str)
			} else if from == "crossrefxml" {
				data, err = crossrefxml.Load(str)
			} else if from == "datacite" {
				data, err = datacite.Load(str)
			} else {
				cmd.PrintErr("Please provide a valid input")
				return
			}
		}

		if err != nil {
			cmd.PrintErr(err)
		}

		var output []byte
		to, _ := cmd.Flags().GetString("to")
		host, _ := cmd.Flags().GetString("host")
		token, _ := cmd.Flags().GetString("token")
		if to == "inveniordm" {
			output, err = inveniordm.Post(data, host, token)
		} else {
			fmt.Println("Please provide a valid service")
			return
		}

		var out bytes.Buffer
		json.Indent(&out, output, "", "  ")
		cmd.Println(out.String())

		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(postCmd)
}
