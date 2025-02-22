/*
Copyright © 2025 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/inveniordm"
	"github.com/front-matter/commonmeta/utils"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
)

// deleteCmd represents the decode command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a InvenioRDM record.",
	Long: `Deletes a InvenioRDM record by id.
	
	Example usage:
	
	commonmeta delete fh8y2-aef76`,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		token, _ := cmd.Flags().GetString("token")

		var record commonmeta.APIResponse

		if len(args) == 0 {
			fmt.Println("Please provide an input")
			return
		}
		input := args[0]
		rid, ok := utils.ValidateRID(input)
		if !ok {
			fmt.Println("Please provide a valid input")
			return
		}
		if host == "" || token == "" {
			fmt.Println("Please provide an inveniordm host and token")
			return
		}
		rl := rate.NewLimiter(rate.Every(60*time.Second), 900) // 900 request every 60 seconds
		client := inveniordm.NewClient(rl, host)
		record.ID = rid
		record, err := inveniordm.DeleteDraftRecord(record, client, token)
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
	rootCmd.AddCommand(deleteCmd)
}
