package cmd

import (
	"fmt"
	"os"

	"github.com/platoon-cc/platoon-cli/internal/settings"
	"github.com/spf13/cobra"
)

func init() {
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Output directory details",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.UserConfigDir()
			if err != nil {
				return err
			}
			fmt.Printf("Config directory: %s/platoon\n", dir)

			dir, err = os.UserCacheDir()
			if err != nil {
				return err
			}
			fmt.Printf("Database cache directory: %s/platoon\n", dir)

			projectId, err := settings.GetActive("project")
			if err != nil {
				return err
			}
			fmt.Printf("Active project Database: %s/platoon/%s.db\n", dir, projectId)
			return nil
		},
	}

	rootCmd.AddCommand(infoCmd)
}
