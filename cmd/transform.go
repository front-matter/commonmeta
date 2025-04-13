/*
Copyright © 2025 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/front-matter/commonmeta/fileutils"
	"github.com/front-matter/commonmeta/ror"
	"github.com/spf13/cobra"
)

// transformCmd represents the transform command
var transformCmd = &cobra.Command{
	Use:   "transform",
	Short: "transform organization metadata",
	Long: `transform organization metadata. Currently only supported for
	ROR (Research Organization Registry) and InvenioRDM formats.
	
	Example usage:
	
	commonmeta transform v1.63-2025-04-03-ror-data_schema_v2.json -f ror -t inveniordm`,
	Run: func(cmd *cobra.Command, args []string) {
		var input string // an identifier, content fetched via API
		var str string
		var err error
		var data []ror.ROR
		var output []byte

		if len(args) > 0 {
			input = args[0]
		}
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		type_, _ := cmd.Flags().GetString("type")
		country, _ := cmd.Flags().GetString("country")
		file, _ := cmd.Flags().GetString("file")

		cmd.SetOut(os.Stdout)
		cmd.SetErr(os.Stderr)

		// extract the file extension and check if output file should be zipped
		// if the file name is empty, set it to the default value
		file, extension, compress := fileutils.GetExtension(file, ".yaml")

		if input != "" && !strings.HasPrefix(input, "--") {
			_, err = os.Stat(input)
			if err != nil {
				cmd.PrintErrf("File not found: %s", input)
				return
			}
			str = input
		}

		if len(args) > 0 {
			input = args[0]
			_, err = os.Stat(input)
			if err != nil {
				fmt.Printf("File not found: %s", input)
				return
			}
			str = input
		}

		if str != "" && from == "ror" {
			if type_ != "" && !slices.Contains(ror.RORTypes, type_) {
				cmd.PrintErr("Please provide a valid type")
				return
			}
			data, err = ror.LoadAll(str)
		} else {
			// if no input is provided, return the built-in ROR vocabulary
			data, err = ror.LoadBuiltin()
		}
		if err != nil {
			fmt.Println("An error occurred:", err)
			return
		}

		// optionally filter the ROR records by type and/or country
		// file "funders.yaml" is a list of funders for InvenioRDM
		// funders.yaml does not include the acronym field
		// file "affiliations_ror.yaml" is a list of affiliations for InvenioRDM
		// affiliations_ror.yaml does not include the country field

		if type_ != "" && !slices.Contains(ror.RORTypes, type_) || country != "" {
			data, err = ror.FilterRecords(data, type_, country, file)
		}
		if err != nil {
			fmt.Println("An error occurred:", err)
			return
		}

		if to == "ror" {
			output, err = ror.WriteAll(data, extension)
		} else if to == "inveniordm" {
			output, err = ror.WriteInvenioRDM(data, extension)
		} else {
			fmt.Println("Please provide a valid output format")
			return
		}

		if file != "" {
			if input != "" && extension == ".yaml" {
				output = append([]byte("# file generated from "+input+"\n\n"), output...)
			}
			if compress {
				err = fileutils.WriteZIPFile(file, output)
			} else {
				err = fileutils.WriteFile(file, output)
			}
		} else {
			fmt.Printf("%s\n", output)
		}

		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(transformCmd)
}
