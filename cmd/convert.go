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
	"github.com/front-matter/commonmeta/crossrefxml"
	"github.com/front-matter/commonmeta/csl"
	"github.com/front-matter/commonmeta/datacite"
	"github.com/front-matter/commonmeta/inveniordm"
	"github.com/front-matter/commonmeta/jsonfeed"
	"github.com/front-matter/commonmeta/schemaorg"
	"github.com/front-matter/commonmeta/utils"

	"github.com/front-matter/commonmeta/crossref"

	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert scholarly metadata from one format to another",
	Long: `Convert scholarly metadata between formats. Currently
supported input formats are Crossref (default) and DataCite DOIs, currently
supported output format and Commonmeta (default). Example usage:

commonmeta 10.5555/12345678`,

	Run: func(cmd *cobra.Command, args []string) {
		var id string  // an identifier, content fetched via API
		var str string // a string, content loaded from a file
		var err error
		var data commonmeta.Data

		// loginID, _ := cmd.Flags().GetString("login_id")
		// loginPassword, _ := cmd.Flags().GetString("login_passwd")
		depositor, _ := cmd.Flags().GetString("depositor")
		email, _ := cmd.Flags().GetString("email")
		registrant, _ := cmd.Flags().GetString("registrant")

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

		if id != "" {
			if from == "" {
				from = utils.FindFromFormatByID(id)
			}
			if from == "crossref" {
				data, err = crossref.Fetch(id)
			} else if from == "crossrefxml" {
				data, err = crossrefxml.Fetch(id)
			} else if from == "datacite" {
				data, err = datacite.Fetch(id)
			} else if from == "inveniordm" {
				data, err = inveniordm.Fetch(id)
			} else if from == "jsonfeed" {
				data, err = jsonfeed.Fetch(id)
			} else if from == "schemaorg" {
				data, err = schemaorg.Fetch(id)
			} else {
				fmt.Println("Please provide a valid input")
				return
			}
			if err != nil {
				fmt.Println("An error occurred:", err)
				return
			}
		} else if str != "" {
			if from == "" {
				from = utils.FindFromFormatByString(str)
			}
			if from == "commonmeta" {
				data, err = commonmeta.Load(str)
			} else if from == "crossref" {
				data, err = crossref.Load(str)
			} else if from == "crossrefxml" {
				data, err = crossrefxml.Load(str)
			} else if from == "datacite" {
				data, err = datacite.Load(str)
			} else if from == "inveniordm" {
				data, err = inveniordm.Load(str)
			} else if from == "csl" {
				data, err = csl.Load(str)
			} else if from == "schemaorg" {
				data, err = schemaorg.Load(str)
			} else {
				cmd.PrintErr("Please provide a valid input")
				return
			}
			if err != nil {
				fmt.Println("An error occurred:", err)
				return
			}
		}

		var output []byte
		to, _ := cmd.Flags().GetString("to")
		if to == "commonmeta" {
			output, err = commonmeta.Write(data)
		} else if to == "csl" {
			output, err = csl.Write(data)
		} else if to == "datacite" {
			output, err = datacite.Write(data)
		} else if to == "schemaorg" {
			output, err = schemaorg.Write(data)
		} else if to == "crossrefxml" {
			account := crossrefxml.Account{
				Depositor:  depositor,
				Email:      email,
				Registrant: registrant,
			}
			output, err = crossrefxml.Write(data, account)
		} else if to == "inveniordm" {
			output, err = inveniordm.Write(data)
		}

		if err != nil {
			cmd.PrintErr(err)
		}

		if to == "crossrefxml" {
			cmd.Printf("%s\n", output)
		} else {
			var out bytes.Buffer
			json.Indent(&out, output, "", "  ")
			cmd.Println(out.String())
		}

		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
}
