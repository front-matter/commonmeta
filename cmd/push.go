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
	"github.com/front-matter/commonmeta/datacite"
	"github.com/front-matter/commonmeta/inveniordm"
	"github.com/front-matter/commonmeta/jsonfeed"

	"github.com/front-matter/commonmeta/crossref"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push scholarly metadata into a service",
	Long: `Convert scholarly metadata between formats and register with
a service. Multiple formats are supported, registration is currently
only supported with InvenioRDM. Example usage:

commonmeta push --sample -f crossref -t inveniordm -h rogue-scholar.org --token mytoken`,

	Run: func(cmd *cobra.Command, args []string) {
		var input string
		var str string // a string, content loaded from a file
		var err error
		var data []commonmeta.Data

		if len(args) > 0 {
			input = args[0]
		}
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
		sample, _ := cmd.Flags().GetBool("sample")

		depositor, _ := cmd.Flags().GetString("depositor")
		email, _ := cmd.Flags().GetString("email")
		registrant, _ := cmd.Flags().GetString("registrant")
		loginID, _ := cmd.Flags().GetString("login_id")
		loginPasswd, _ := cmd.Flags().GetString("login_passwd")
		to, _ := cmd.Flags().GetString("to")
		host, _ := cmd.Flags().GetString("host")
		token, _ := cmd.Flags().GetString("token")
		legacyKey, _ := cmd.Flags().GetString("legacyKey")

		cmd.SetOut(os.Stdout)
		cmd.SetErr(os.Stderr)

		if input != "" {
			_, err = os.Stat(input)
			if err != nil {
				cmd.PrintErrf("File not found: %s", input)
				return
			}
			str = input
		}

		if sample && from == "crossref" {
			data, err = crossref.FetchAll(number, member, type_, sample, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense, hasArchive)
		} else if sample && from == "datacite" {
			data, err = datacite.FetchAll(number, sample)
		} else if str != "" && from == "commonmeta" {
			data, err = commonmeta.LoadAll(str)
		} else if str != "" && from == "crossref" {
			data, err = crossref.LoadAll(str)
		} else if str != "" && from == "crossrefxml" {
			data, err = crossrefxml.LoadAll(str)
		} else if str != "" && from == "datacite" {
			data, err = datacite.LoadAll(str)
		} else if str != "" && from == "jsonfeed" {
			data, err = jsonfeed.LoadAll(str)
		} else {
			fmt.Println("Please provide a valid input format")
			return
		}

		if err != nil {
			cmd.PrintErr(err)
		}

		var records []commonmeta.APIResponse
		if to == "crossrefxml" {
			account := crossrefxml.Account{
				Depositor:   depositor,
				Email:       email,
				Registrant:  registrant,
				LoginID:     loginID,
				LoginPasswd: loginPasswd,
			}
			records, err = crossrefxml.UpsertAll(data, account)
		} else if to == "inveniordm" {
			records, err = inveniordm.UpsertAll(data, host, token, legacyKey)
		} else {
			fmt.Println("Please provide a valid service")
			return
		}

		if err != nil {
			cmd.PrintErr(err)
		}

		var output []byte
		var out bytes.Buffer
		output, err = json.Marshal(records)
		if err != nil {
			cmd.PrintErr(err)
		}

		json.Indent(&out, output, "", "  ")
		cmd.Println(out.String())

		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
