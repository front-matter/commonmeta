/*
Copyright © 2024 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"fmt"

	"github.com/front-matter/commonmeta/ghost"
	"github.com/spf13/cobra"
)

// encodeCmd represents the encode command
var updateGhostAPICmd = &cobra.Command{
	Use:   "update-ghost-post",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide an input")
			return
		}
		id := args[0]
		apiKey, _ := cmd.Flags().GetString("api-key")
		apiURL, _ := cmd.Flags().GetString("api-url")
		output, err := ghost.UpdateGhostPost(id, apiKey, apiURL)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(updateGhostAPICmd)

	rootCmd.PersistentFlags().StringP("api-key", "", "", "Ghost API key")
	rootCmd.PersistentFlags().StringP("api-url", "", "", "Ghost API URL")
}
