/*
Copyright © 2025 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/front-matter/commonmeta/inveniordm"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup metadata",
	Long: `Setup metadata, currently only InvenioRDM is supported. Example usage:

	commonmeta setup --to inveniordm --host example.org --token mytoken`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var output []byte

		to, _ := cmd.Flags().GetString("to")

		host, _ := cmd.Flags().GetString("host")
		token, _ := cmd.Flags().GetString("token")

		if to == "" || to == "commonmeta" {
			to = "inveniordm"
		}

		if to == "inveniordm" {
			if host == "" || token == "" {
				fmt.Println("Please provide an inveniordm host and token")
				return
			}
			output, err = inveniordm.CreateSubjectCommunities(host, token)
		}

		var out bytes.Buffer
		json.Indent(&out, output, "", "  ")
		cmd.Println(out.String())

		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
