/*
Copyright © 2024-2025 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/crossrefxml"
	"github.com/front-matter/commonmeta/csl"
	"github.com/front-matter/commonmeta/datacite"
	"github.com/front-matter/commonmeta/inveniordm"
	"github.com/front-matter/commonmeta/jsonfeed"
	"github.com/front-matter/commonmeta/openalex"
	"github.com/front-matter/commonmeta/ror"
	"github.com/front-matter/commonmeta/schemaorg"
	"github.com/front-matter/commonmeta/utils"
	"golang.org/x/time/rate"

	"github.com/front-matter/commonmeta/crossref"

	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert scholarly metadata from one format to another",
	Long: `Convert scholarly metadata between formats. Example usage:

commonmeta 10.5555/12345678`,

	Run: func(cmd *cobra.Command, args []string) {
		var input, id, identifierType, str string
		var err error
		var data commonmeta.Data
		var orgdata ror.ROR
		var output []byte

		// loginID, _ := cmd.Flags().GetString("login_id")
		// loginPassword, _ := cmd.Flags().GetString("login_passwd")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		fromHost, _ := cmd.Flags().GetString("from-host")
		depositor, _ := cmd.Flags().GetString("depositor")
		email, _ := cmd.Flags().GetString("email")
		registrant, _ := cmd.Flags().GetString("registrant")
		match, _ := cmd.Flags().GetBool("match")

		cmd.SetOut(os.Stdout)
		cmd.SetErr(os.Stderr)

		if len(args) == 0 {
			fmt.Println("Please provide an input")
			return
		}

		if len(args) > 0 {
			input = args[0]
			id, identifierType = utils.ValidateID(input)
			if slices.Contains(commonmeta.WorkTypes, identifierType) {
				if from == "" {
					from = "commonmeta"
				}
			} else if slices.Contains(commonmeta.PersonTypes, identifierType) {
				// TODO
			} else if slices.Contains(commonmeta.OrganizationTypes, identifierType) {
				if from == "" {
					from = "ror"
				}
				if to == "" || to == "commonmeta" {
					to = "ror"
				}
			}
		}

		if id == "" {
			_, err = os.Stat(input)
			if err != nil {
				fmt.Printf("File not found: %s", input)
				return
			}
			str = input
		}

		if id != "" {
			if from == "" {
				from = utils.FindFromFormatByID(id)
			}
			if from == "crossref" {
				data, err = crossref.Fetch(id, match)
			} else if from == "crossrefxml" {
				data, err = crossrefxml.Fetch(id)
			} else if from == "datacite" {
				data, err = datacite.Fetch(id, match)
			} else if from == "inveniordm" {
				rl := rate.NewLimiter(rate.Every(10*time.Second), 100)
				client := inveniordm.NewClient(rl, fromHost)
				data, err = inveniordm.Fetch(id, match, client)
			} else if from == "jsonfeed" {
				data, err = jsonfeed.Fetch(id)
			} else if from == "openalex" {
				r := openalex.NewReader(email)
				data, err = r.Fetch(id)
			} else if from == "schemaorg" {
				data, err = schemaorg.Fetch(id, match)
			} else if slices.Contains(commonmeta.OrganizationTypes, identifierType) && from == "ror" {
				orgdata, err = ror.Search(id)
				if orgdata.ID == "" {
					cmd.Println("No match found")
					return
				}
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
				data, err = crossref.Load(str, match)
			} else if from == "crossrefxml" {
				data, err = crossrefxml.Load(str)
			} else if from == "datacite" {
				data, err = datacite.Load(str, match)
			} else if from == "inveniordm" {
				data, err = inveniordm.Load(str, match)
			} else if from == "csl" {
				data, err = csl.Load(str)
			} else if from == "schemaorg" {
				data, err = schemaorg.Load(str, match)
			} else {
				cmd.PrintErr("Please provide a valid input")
				return
			}
			if err != nil {
				fmt.Println("An error occurred:", err)
				return
			}
		}

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
		} else if data.ID != "" && to == "inveniordm" {
			output, err = inveniordm.Write(data, fromHost)
		} else if orgdata.ID != "" && to == "inveniordm" {
			output, err = ror.WriteInvenioRDM(orgdata)
		} else if orgdata.ID != "" && to == "ror" {
			output, err = ror.Write(orgdata)
		}

		if err != nil {
			cmd.PrintErr(err)
		}

		if to == "crossrefxml" || to == "inveniordm" {
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
