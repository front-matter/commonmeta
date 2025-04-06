/*
Copyright © 2025 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/front-matter/commonmeta/fileutils"
	"github.com/front-matter/commonmeta/ror"
	"github.com/front-matter/commonmeta/utils"
	"github.com/spf13/cobra"
)

// obtainCmd represents the obtain command
var obtainCmd = &cobra.Command{
	Use:   "obtain",
	Short: "obtain vocabulary",
	Long: `obtain a vocabulary from a url. Currently only supported for
	ROR (Research Organization Registry) and InvenioRDM awards vocabulary.
	
	Example usage:
	
	commonmeta obtain v1.63-2025-04-03-ror-data_schema_v2.json -f ror -t inveniordm`,
	Run: func(cmd *cobra.Command, args []string) {
		var id string // an identifier, content fetched via API
		var str string
		var err error
		var data []ror.ROR
		var output []byte

		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		file, _ := cmd.Flags().GetString("file")
		compress, _ := cmd.Flags().GetBool("compress")

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

		if str != "" {
			if from == "ror" {
				data, err = ror.LoadAll(str)
			} else {
				cmd.PrintErr("Please provide a valid input")
				return
			}
			if err != nil {
				fmt.Println("An error occurred:", err)
				return
			}
		}

		supportedFormats := []string{"inveniordm"}
		if !slices.Contains(supportedFormats, to) {
			cmd.PrintErr("Please provide a valid output format")
			return
		}

		output, err = ror.WriteAll(data, to)

		if file != "" {
			output = append([]byte("# file generated from "+input+"\n\n"), output...)
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
	rootCmd.AddCommand(obtainCmd)
}
