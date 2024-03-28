package cmd

import (
	"github.com/platoon-cc/platoon-cli/client"
	"github.com/spf13/cobra"
)

var orgCmd = &cobra.Command{
	Use: "org",
}

var listCmd = &cobra.Command{
	Use: "ls",
	RunE: func(cmd *cobra.Command, args []string) error {
		platoon := client.New()
		platoon.OrgList()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(orgCmd)
	orgCmd.AddCommand(listCmd)
}
