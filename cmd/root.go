package cmd

import (
	"fmt"
	"os"

	"github.com/platoon-cc/platoon-cli/internal/settings"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "platoon-cli",
}

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Platoon CLI v0.0.1 -- HEAD")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(settings.Config)
	cobra.OnFinalize(settings.Save)
	rootCmd.AddCommand(versionCmd)
}
