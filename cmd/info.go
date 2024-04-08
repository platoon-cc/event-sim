package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "output directory details",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.UserConfigDir()
			if err != nil {
				return err
			}
			fmt.Printf("Config directory: %s\n", dir)

			dir, err = os.UserCacheDir()
			if err != nil {
				return err
			}
			fmt.Printf("Database cache directory: %s\n", dir)
			return nil
		},
	}

	rootCmd.AddCommand(infoCmd)
}
