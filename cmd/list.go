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
	"strings"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/crossref"
	"github.com/front-matter/commonmeta/crossrefxml"
	"github.com/front-matter/commonmeta/csl"
	"github.com/front-matter/commonmeta/datacite"
	"github.com/front-matter/commonmeta/fileutils"
	"github.com/front-matter/commonmeta/inveniordm"
	"github.com/front-matter/commonmeta/jsonfeed"
	"github.com/front-matter/commonmeta/openalex"
	"github.com/front-matter/commonmeta/ror"
	"github.com/front-matter/commonmeta/schemaorg"
	"golang.org/x/time/rate"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A list of scholarly metadata",
	Long: `A list of scholarly metadata retrieved via file or API. For example:

	commonmeta list --number 10 --member 78 --type journal-article - f crossref,
	commonmeta list --number 10 --client cern.zenodo --type dataset -f datacite,
	commonmeta list --number 10 --from inveniordm --from-host rogue-scholar.org --community front_matter`,
	Run: func(cmd *cobra.Command, args []string) {
		var input string // an identifier, content fetched via API
		var str string   // a string, content loaded from a file
		var err error
		var data []commonmeta.Data
		var orgdata []ror.ROR
		var extension string
		var output []byte

		number, _ := cmd.Flags().GetInt("number")
		page, _ := cmd.Flags().GetInt("page")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")

		client_, _ := cmd.Flags().GetString("client")
		member, _ := cmd.Flags().GetString("member")
		type_, _ := cmd.Flags().GetString("type")
		year, _ := cmd.Flags().GetString("year")
		country, _ := cmd.Flags().GetString("country")
		dateUpdated, _ := cmd.Flags().GetString("date-updated")
		language, _ := cmd.Flags().GetString("language")
		orcid, _ := cmd.Flags().GetString("orcid")
		affiliation, _ := cmd.Flags().GetString("affiliation")
		ror_, _ := cmd.Flags().GetString("ror")
		fromHost, _ := cmd.Flags().GetString("from-host")
		community, _ := cmd.Flags().GetString("community")
		subject, _ := cmd.Flags().GetString("subject")
		dataVersion, _ := cmd.Flags().GetString("data-version")
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
		match, _ := cmd.Flags().GetBool("match")

		depositor, _ := cmd.Flags().GetString("depositor")
		email, _ := cmd.Flags().GetString("email")
		registrant, _ := cmd.Flags().GetString("registrant")

		cmd.SetOut(os.Stdout)
		cmd.SetErr(os.Stderr)

		if len(args) > 0 {
			input = args[0]
		}

		// extract the file extension and check if output file should be compressed
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

		if from == "ror" && (to == "" || to == "commonmeta") {
			to = "ror"
		}

		if from == "commonmeta" {
			data, err = commonmeta.LoadAll(str)
		} else if str != "" && from == "crossref" {
			data, err = crossref.LoadAll(str, match)
		} else if str != "" && from == "crossrefxml" {
			data, err = crossrefxml.LoadAll(str)
		} else if str != "" && from == "datacite" {
			data, err = datacite.LoadAll(str, match)
		} else if str != "" && from == "jsonfeed" {
			data, err = jsonfeed.LoadAll(str)
		} else if str != "" && from == "csl" {
			data, err = csl.LoadAll(str)
		} else if from == "crossref" {
			data, err = crossref.FetchAll(number, page, member, type_, sample, year, orcid, ror_, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense, hasArchive, match)
		} else if from == "datacite" {
			data, err = datacite.FetchAll(number, page, client_, type_, sample, year, language, orcid, ror_, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense, match)
		} else if from == "openalex" {
			r := openalex.NewReader(email)
			data, err = r.FetchAll(number, page, member, type_, sample, "", year, orcid, ror_, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense, hasArchive)
		} else if from == "inveniordm" {
			rl := rate.NewLimiter(rate.Every(10*time.Second), 100)
			client := inveniordm.NewClient(rl, fromHost)
			data, err = inveniordm.FetchAll(number, page, client, community, subject, type_, year, language, orcid, affiliation, ror_, hasORCID, hasROR, match)
		} else if from == "jsonfeed" {
			data, err = jsonfeed.FetchAll(number, page, community, isArchived)
		} else if str != "" && from == "ror" {
			if type_ != "" && !slices.Contains(ror.RORTypes, type_) {
				cmd.PrintErr("Please provide a valid type")
				return
			}
			orgdata, err = ror.LoadAll(str)
		} else if str == "" && dataVersion != "" && from == "ror" {
			// download the ROR data dump with the specified version
			orgdata, err = ror.FetchAll(dataVersion)
			input = ror.Basename(dataVersion)
		} else if str == "" && from == "ror" {
			// if no input is provided, use ROR data dump
			orgdata, err = ror.LoadBuiltin()
			input = ror.Basename("v1.64")
		} else {
			fmt.Println("Please provide a valid input format")
			return
		}
		if err != nil {
			fmt.Println("An error occurred:", err)
			return
		}

		// optionally filter orgdata by type, country, number and page
		if len(orgdata) > 0 && (type_ != "" && !slices.Contains(ror.RORTypes, type_) || country != "" || file != "" || dateUpdated != "" || number != 0) {
			orgdata, err = ror.FilterList(orgdata, type_, country, dateUpdated, file, number, page)
		}
		if err != nil {
			fmt.Println("An error occurred:", err)
			return
		}

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
		} else if data != nil && to == "inveniordm" {
			output, err = inveniordm.WriteAll(data, fromHost)
		} else if len(orgdata) > 0 && to == "ror" {
			output, err = ror.WriteAll(orgdata, extension)
		} else if len(orgdata) > 0 && to == "inveniordm" {
			output, err = ror.WriteAllInvenioRDM(orgdata, extension)
		} else {
			fmt.Println("Please provide a valid output format")
			return
		}

		if to != "crossrefxml" && extension == ".json" {
			var out bytes.Buffer
			json.Indent(&out, output, "", "  ")
			output = out.Bytes()
		}

		if file != "" {
			if input != "" && extension == ".yaml" {
				output = append([]byte("# file generated from "+input+"\n\n"), output...)
			}
			switch compress {
			case "gz":
				err = fileutils.WriteGZFile(file, output)
			case "zip":
				err = fileutils.WriteZIPFile(file, output)
			default:
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
