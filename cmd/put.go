/*
Copyright © 2024-2025 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/crossrefxml"
	"github.com/front-matter/commonmeta/csl"
	"github.com/front-matter/commonmeta/datacite"
	"github.com/front-matter/commonmeta/inveniordm"
	"github.com/front-matter/commonmeta/jsonfeed"
	"github.com/front-matter/commonmeta/schemaorg"
	"github.com/front-matter/commonmeta/utils"
	"golang.org/x/time/rate"

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

		depositor, _ := cmd.Flags().GetString("depositor")
		email, _ := cmd.Flags().GetString("email")
		registrant, _ := cmd.Flags().GetString("registrant")
		loginID, _ := cmd.Flags().GetString("login_id")
		loginPasswd, _ := cmd.Flags().GetString("login_passwd")
		to, _ := cmd.Flags().GetString("to")
		host, _ := cmd.Flags().GetString("host")
		fromHost, _ := cmd.Flags().GetString("from-host")
		token, _ := cmd.Flags().GetString("token")
		legacyKey, _ := cmd.Flags().GetString("legacyKey")
		client_, _ := cmd.Flags().GetString("client")
		password, _ := cmd.Flags().GetString("password")
		development, _ := cmd.Flags().GetBool("development")
		match, _ := cmd.Flags().GetBool("match")

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
			switch from {
			case "crossref":
				data, err = crossref.Fetch(id, match)
			case "crossrefxml":
				data, err = crossrefxml.Fetch(id)
			case "datacite":
				data, err = datacite.Fetch(id, match)
			case "inveniordm":
				rl := rate.NewLimiter(rate.Every(10*time.Second), 100)
				client := inveniordm.NewClient(rl, fromHost)
				data, err = inveniordm.Fetch(id, match, client)
			case "jsonfeed":
				data, err = jsonfeed.Fetch(id)
			case "schemaorg":
				data, err = schemaorg.Fetch(id, match)
			default:
				fmt.Println("Please provide a valid input")
				return
			}
			if err != nil {
				fmt.Println("An error occurred:", err)
				return
			}
		} else if str != "" {
			if from == "" {
				from = utils.FindFromFormatByID(id)
			}
			if from == "commonmeta" {
				data, err = commonmeta.Load(str)
			} else if from == "crossref" {
				data, err = crossref.Load(str, match)
			} else if from == "crossrefxml" {
				data, err = crossrefxml.Load(str)
			} else if from == "datacite" {
				data, err = datacite.Load(str, match)
			} else if from == "csl" {
				data, err = csl.Load(str)
			} else {
				cmd.PrintErr("Please provide a valid input format")
				return
			}
			if err != nil {
				fmt.Println("An error occurred:", err)
				return
			}
		}

		if err != nil {
			cmd.PrintErr(err)
		}

		var record commonmeta.APIResponse
		switch to {
		case "crossrefxml":
			account := crossrefxml.Account{
				Depositor:   depositor,
				Email:       email,
				Registrant:  registrant,
				LoginID:     loginID,
				LoginPasswd: loginPasswd,
			}
			record, err = crossrefxml.Upsert(record, account, legacyKey, data)
		case "datacite":
			account := datacite.Account{
				Client:      client_,
				Password:    password,
				Development: development,
			}
			record, err = datacite.Upsert(record, account, data)
		case "inveniordm":
			if host == "" || token == "" {
				fmt.Println("Please provide an inveniordm host and token")
				return
			}
			rl := rate.NewLimiter(rate.Every(10*time.Second), 100) // 100 request every 10 seconds
			client := inveniordm.NewClient(rl, host)
			record, err = inveniordm.Upsert(record, fromHost, token, legacyKey, data, client)
			if err != nil {
				cmd.PrintErr(err)
			}
		default:
			fmt.Println("Please provide a valid service")
			return
		}

		if err != nil {
			cmd.PrintErr(err)
		}

		var output []byte
		var out bytes.Buffer
		output, err = json.Marshal(record)
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
	rootCmd.AddCommand(putCmd)
}
