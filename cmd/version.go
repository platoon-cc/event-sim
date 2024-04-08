package cmd

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed .version
var version string

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Platoon CLI %s\n", "v"+version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
