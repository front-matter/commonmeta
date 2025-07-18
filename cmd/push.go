/*
Copyright © 2024-2025 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/crossrefxml"
	"github.com/front-matter/commonmeta/csl"
	"github.com/front-matter/commonmeta/datacite"
	"github.com/front-matter/commonmeta/inveniordm"
	"github.com/front-matter/commonmeta/jsonfeed"
	"golang.org/x/time/rate"

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
		page, _ := cmd.Flags().GetInt("page")
		from, _ := cmd.Flags().GetString("from")

		client_, _ := cmd.Flags().GetString("client")
		member, _ := cmd.Flags().GetString("member")
		type_, _ := cmd.Flags().GetString("type")
		year, _ := cmd.Flags().GetString("year")
		language, _ := cmd.Flags().GetString("language")
		orcid, _ := cmd.Flags().GetString("orcid")
		affiliation, _ := cmd.Flags().GetString("affiliation")
		ror, _ := cmd.Flags().GetString("ror")
		fromHost, _ := cmd.Flags().GetString("from-host")
		fromToken, _ := cmd.Flags().GetString("from-token")
		community, _ := cmd.Flags().GetString("community")
		subject, _ := cmd.Flags().GetString("subject")
		hasORCID, _ := cmd.Flags().GetBool("has-orcid")
		hasROR, _ := cmd.Flags().GetBool("has-ror-id")
		hasReferences, _ := cmd.Flags().GetBool("has-references")
		hasRelation, _ := cmd.Flags().GetBool("has-relation")
		hasAbstract, _ := cmd.Flags().GetBool("has-abstract")
		hasAward, _ := cmd.Flags().GetBool("has-award")
		hasLicense, _ := cmd.Flags().GetBool("has-license")
		hasArchive, _ := cmd.Flags().GetBool("has-archive")
		isArchived, _ := cmd.Flags().GetBool("is-archived")
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
		password, _ := cmd.Flags().GetString("password")
		development, _ := cmd.Flags().GetBool("development")
		match, _ := cmd.Flags().GetBool("match")

		cmd.SetOut(os.Stdout)
		cmd.SetErr(os.Stderr)

		if input != "" && !strings.HasPrefix(input, "--") {
			_, err = os.Stat(input)
			if err != nil {
				cmd.PrintErrf("File not found: %s", input)
				return
			}
			str = input
		}

		if from == "crossref" {
			data, err = crossref.FetchAll(number, page, member, type_, sample, year, orcid, ror, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense, hasArchive, match)
		} else if from == "datacite" {
			data, err = datacite.FetchAll(number, page, client_, type_, sample, year, language, orcid, ror, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense, match)
		} else if from == "inveniordm" {
			rl := rate.NewLimiter(rate.Every(10*time.Second), 100)
			client := inveniordm.NewClient(rl, fromHost)
			data, err = inveniordm.FetchAll(number, page, fromToken, community, subject, type_, year, language, orcid, affiliation, ror, hasORCID, hasROR, match, client)
		} else if from == "jsonfeed" {
			data, err = jsonfeed.FetchAll(number, page, community, isArchived)
		} else if str != "" && from == "commonmeta" {
			data, err = commonmeta.LoadAll(str)
		} else if str != "" && from == "crossref" {
			data, err = crossref.LoadAll(str, match)
		} else if str != "" && from == "crossrefxml" {
			data, err = crossrefxml.LoadAll(str)
		} else if str != "" && from == "datacite" {
			data, err = datacite.LoadAll(str, match)
		} else if str != "" && from == "inveniordm" {
			data, err = inveniordm.LoadAll(str, match)
		} else if str != "" && from == "jsonfeed" {
			data, err = jsonfeed.LoadAll(str)
		} else if str != "" && from == "csl" {
			data, err = csl.LoadAll(str)
		} else {
			fmt.Println("Please provide a valid input format")
			return
		}

		if err != nil {
			cmd.PrintErr(err)
		}

		var records []commonmeta.APIResponse
		switch to {
		case "crossrefxml":
			account := crossrefxml.Account{
				Depositor:   depositor,
				Email:       email,
				Registrant:  registrant,
				LoginID:     loginID,
				LoginPasswd: loginPasswd,
			}
			records, err = crossrefxml.UpsertAll(data, account, legacyKey)
		case "datacite":
			account := datacite.Account{
				Client:      client_,
				Password:    password,
				Development: development,
			}
			records, err = datacite.UpsertAll(data, account)
		case "inveniordm":
			if host == "" || token == "" {
				fmt.Println("Please provide an inveniordm host and token")
				return
			}
			rl := rate.NewLimiter(rate.Every(10*time.Second), 100) // 100 request every 10 seconds
			client := inveniordm.NewClient(rl, host)
			records, err = inveniordm.UpsertAll(data, fromHost, token, legacyKey, client)
		default:
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
