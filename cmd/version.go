/*
Copyright © 2024 Front Matter <info@front-matter.io>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of commonmeta",
	Long:  `All software has versions. This is commonmeta's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Commonmeta v0.3.7 -- HEAD")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
