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
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/csl"
	"github.com/front-matter/commonmeta/fileutils"
	"github.com/front-matter/commonmeta/inveniordm"
	"github.com/front-matter/commonmeta/ror"

	"github.com/front-matter/commonmeta/crossref"
	"github.com/front-matter/commonmeta/crossrefxml"
	"github.com/front-matter/commonmeta/jsonfeed"

	"github.com/front-matter/commonmeta/datacite"

	"github.com/front-matter/commonmeta/schemaorg"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A list of works",
	Long: `A list of works. Currently only available for
	the Crossref and DataCite provider. Options include numnber of works, 
	work type, and Crossref member id or DataCite client id. For example:

	commonmeta list --number 10 --member 78 --type journal-article - f crossref,
	commonmeta list --number 10 --client cern.zenodo --type dataset -f datacite,
	commonmeta list --number 10 --from inveniordm --from-host rogue-scholar.org --community front_matter`,
	Run: func(cmd *cobra.Command, args []string) {
		var input string
		var str string // a string, content loaded from a file
		var err error
		var data []commonmeta.Data
		var extension string

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
		ror_, _ := cmd.Flags().GetString("ror")
		fromHost, _ := cmd.Flags().GetString("from-host")
		community, _ := cmd.Flags().GetString("community")
		subject, _ := cmd.Flags().GetString("subject")
		vocabulary, _ := cmd.Flags().GetBool("vocabulary")
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
		file, _ := cmd.Flags().GetString("file")

		depositor, _ := cmd.Flags().GetString("depositor")
		email, _ := cmd.Flags().GetString("email")
		registrant, _ := cmd.Flags().GetString("registrant")

		cmd.SetOut(os.Stdout)
		cmd.SetErr(os.Stderr)

    // extract the file extension and check if output file should be zipped
		// if the file name is empty, set it to the default value
    file, extension, compress := fileutils.GetExtension(file, ".json")

		if input != "" && !strings.HasPrefix(input, "--") {
			_, err = os.Stat(input)
			if err != nil {
				cmd.PrintErrf("File not found: %s", input)
				return
			}
			str = input
		}

		if from == "commonmeta" {
			data, err = commonmeta.LoadAll(str)
		} else if str != "" && from == "crossref" {
			data, err = crossref.LoadAll(str)
		} else if str != "" && from == "crossrefxml" {
			data, err = crossrefxml.LoadAll(str)
		} else if str != "" && from == "datacite" {
			data, err = datacite.LoadAll(str)
		} else if str != "" && from == "jsonfeed" {
			data, err = jsonfeed.LoadAll(str)
		} else if str != "" && from == "csl" {
			data, err = csl.LoadAll(str)
		} else if from == "crossref" {
			data, err = crossref.FetchAll(number, page, member, type_, sample, year, orcid, ror_, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense, hasArchive)
		} else if from == "datacite" {
			data, err = datacite.FetchAll(number, page, client_, type_, sample, year, language, orcid, ror_, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense)
		} else if from == "inveniordm" {
			data, err = inveniordm.FetchAll(number, page, fromHost, community, subject, type_, year, language, orcid, affiliation, ror_, hasORCID, hasROR)
		} else if from == "jsonfeed" {
			data, err = jsonfeed.FetchAll(number, page, community, isArchived)
		} else {
			fmt.Println("Please provide a valid input format")
			return
		}
		if err != nil {
			fmt.Println("An error occurred:", err)
			return
		}

		if err != nil {
			cmd.PrintErr(err)
		}

		var output []byte
		to, _ := cmd.Flags().GetString("to")
		if to == "commonmeta" {
			output, err = commonmeta.WriteAll(data, extension)
		} else if to == "csl" {
			output, err = csl.WriteAll(data)
		} else if to == "datacite" {
			output, err = datacite.WriteAll(data)
		} else if to == "crossrefxml" {
			account := crossrefxml.Account{
				Depositor:  depositor,
				Email:      email,
				Registrant: registrant,
			}
			output, err = crossrefxml.WriteAll(data, account)
		} else if to == "schemaorg" {
			output, err = schemaorg.WriteAll(data)
		} else if to == "inveniordm" {
			output, err = inveniordm.WriteAll(data)
		}

		if to != "crossrefxml" && extension == ".json" {
			var out bytes.Buffer
			json.Indent(&out, output, "", "  ")
			output = out.Bytes()
		}

		if file != "" {
			if compress {
				err = fileutils.WriteZIPFile(file, output)
			} else {
				err = fileutils.WriteFile(file, output)
			}
		} else {
			fmt.Printf("%s\n", output)
		}

		if to == "inveniordm" && vocabulary {
			file = "affiliations_ror.yaml"
			roroutput, err := ror.ExtractAll(data)
			if err != nil {
				cmd.PrintErr(err)
			}
			today := time.Now().UTC()
			roroutput = append([]byte("# file generated from "+from+" query on "+today.Format("2006-01-02")+".\n\n"), roroutput...)
			err = fileutils.WriteFile(file, roroutput)
			if err != nil {
				cmd.PrintErr(err)
			}
			fmt.Printf("Found ROR IDs written to %s\n", file)
		}

		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
