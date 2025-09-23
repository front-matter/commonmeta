/*
Copyright © 2025 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/front-matter/commonmeta/inveniordm"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
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

		fromHost, _ := cmd.Flags().GetString("from-host")
		fromToken, _ := cmd.Flags().GetString("from-token")
		host, _ := cmd.Flags().GetString("host")
		token, _ := cmd.Flags().GetString("token")
		action, _ := cmd.Flags().GetString("action")

		if to == "" || to == "commonmeta" {
			to = "inveniordm"
		}
		if action == "" {
			action = "create_subject_communities"
		}

		if to == "inveniordm" {
			if host == "" || token == "" {
				fmt.Println("Please provide an inveniordm host and token")
				return
			}
			rl := rate.NewLimiter(rate.Every(10*time.Second), 100)
			client := inveniordm.NewClient(rl, host)
			oldClient := inveniordm.NewClient(rl, fromHost)
			switch action {
			case "create_subject_communities":
				output, err = inveniordm.CreateSubjectCommunities(token, client)
			case "transfer_blog_communities":
				output, err = inveniordm.TransferCommunities("blog", token, fromToken, oldClient, client)
			case "transfer_topic_communities":
				output, err = inveniordm.TransferCommunities("topic", token, fromToken, oldClient, client)
			}
			if err != nil {
				cmd.PrintErr(err)
				return
			}
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
