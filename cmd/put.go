/*
Copyright © 2024 Front Matter <info@front-matter.io>
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

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Put scholarly metadata into a service",
	Long: `Convert scholarly metadata between formats and register with
a service. Multiple formats are supported, registration is currently
only supported with InvenioRDM. Example usage:

commonmeta put 10.5555/12345678 -f crossref -t inveniordm -h rogue-scholar.org --token mytoken`,

	Run: func(cmd *cobra.Command, args []string) {
		var id string  // an identifier, content fetched via API
		var str string // a string, content loaded from a file
		var err error
		var data commonmeta.Data

		to, _ := cmd.Flags().GetString("to")
		host, _ := cmd.Flags().GetString("host")
		token, _ := cmd.Flags().GetString("token")

		cmd.SetOut(os.Stdout)
		cmd.SetErr(os.Stderr)

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
				cmd.PrintErr("Please provide a valid input format")
				return
			}
		}

		if err != nil {
			cmd.PrintErr(err)
		}

		var output []byte
		if to == "inveniordm" {
			var record inveniordm.APIResponse
			record, err = inveniordm.Upsert(record, host, token, data)
			if err != nil {
				cmd.PrintErr(err)
			}
			output, err = json.Marshal(record)
			if err != nil {
				cmd.PrintErr(err)
			}
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
	rootCmd.AddCommand(putCmd)
}
