/*
Copyright © 2024 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "commonmeta",
	Version: "v0.18.2",
	Short:   "Convert scholarly metadata between formats",
	Long: `Convert scholarly metadata between formats. Currently
supported formats include Crossref, DataCite, Schema.org, CSL, InvenioRDM, 
JSON Feed, and ROR. Example usage:

commonmeta convert 10.5555/12345678
commonmeta convert org 10.5555/12345678`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("root called")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("from", "f", "", "the format to convert from")
	rootCmd.PersistentFlags().StringP("to", "t", "commonmeta", "the format to convert to")

	rootCmd.PersistentFlags().IntP("number", "n", 0, "number of results")
	rootCmd.PersistentFlags().IntP("page", "", 1, "page number")
	rootCmd.PersistentFlags().StringP("file", "", "", "file to write to")
	rootCmd.PersistentFlags().StringP("client", "", "", "DataCite client ID")
	rootCmd.PersistentFlags().StringP("member", "", "", "Crossref member ID")
	rootCmd.PersistentFlags().StringP("type", "", "", "work type")
	rootCmd.PersistentFlags().StringP("year", "", "", "work publication year")
	rootCmd.PersistentFlags().StringP("language", "", "", "work language")
	rootCmd.PersistentFlags().StringP("orcid", "", "", "ORCID ID")
	rootCmd.PersistentFlags().StringP("affiliation", "", "", "affiliation name")
	rootCmd.PersistentFlags().StringP("ror", "", "", "ROR ID")
	rootCmd.PersistentFlags().StringP("from-host", "", "", "from InvenioRDM host")
	rootCmd.PersistentFlags().StringP("community", "", "", "InvenioRDM community slug")
	rootCmd.PersistentFlags().StringP("subject", "", "", "subject area")
	rootCmd.PersistentFlags().StringP("country", "", "", "country code")
	rootCmd.PersistentFlags().BoolP("sample", "", false, "random sample")
	rootCmd.PersistentFlags().BoolP("has-orcid", "", false, "has one or more ORCID IDs")
	rootCmd.PersistentFlags().BoolP("has-ror-id", "", false, "has one or more ROR IDs")
	rootCmd.PersistentFlags().BoolP("has-references", "", false, "has references")
	rootCmd.PersistentFlags().BoolP("has-relation", "", false, "has relation")
	rootCmd.PersistentFlags().BoolP("has-abstract", "", false, "has abstract")
	rootCmd.PersistentFlags().BoolP("has-award", "", false, "has award")
	rootCmd.PersistentFlags().BoolP("has-license", "", false, "has license")
	rootCmd.PersistentFlags().BoolP("has-archive", "", false, "has archive")
	rootCmd.PersistentFlags().BoolP("is-archived", "", true, "is archived")
	rootCmd.PersistentFlags().BoolP("vocabulary", "", false, "with vocabulary")

	// needed for DOI registration
	rootCmd.PersistentFlags().StringP("prefix", "", "", "DOI prefix")
	rootCmd.PersistentFlags().BoolP("development", "", false, "Development mode")

	rootCmd.PersistentFlags().StringP("login_id", "", "", "Crossref account login")
	rootCmd.PersistentFlags().StringP("login_passwd", "", "", "Crossref account password")
	rootCmd.PersistentFlags().StringP("depositor", "", "", "Crossref account depositor")
	rootCmd.PersistentFlags().StringP("email", "", "", "Crossref account email")
	rootCmd.PersistentFlags().StringP("registrant", "", "", "Crossref account registrant")
	rootCmd.PersistentFlags().StringP("host", "", "", "InvenioRDM host")
	rootCmd.PersistentFlags().StringP("token", "", "", "API token")
	rootCmd.PersistentFlags().StringP("password", "", "", "DataCite client password")
	rootCmd.PersistentFlags().StringP("legacyKey", "", "", "Legacy API token")
}
