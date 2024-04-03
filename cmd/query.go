package cmd

import (
	"github.com/platoon-cc/platoon-cli/internal/processor"
	"github.com/platoon-cc/platoon-cli/internal/settings"
	"github.com/spf13/cobra"
)

func init() {
	queryCmd := &cobra.Command{
		Use: "query",
	}

	rootCmd.AddCommand(queryCmd)
	queryCmd.AddCommand(&cobra.Command{
		Use: "challengeScores",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectId := "local"
			isLocal, err := cmd.Flags().GetBool("local")
			if err != nil {
				return err
			}
			if !isLocal {
				projectId, err = settings.GetActive("project")
				if err != nil {
					return err
				}
			}

			processor, err := processor.New(projectId)
			if err != nil {
				return err
			}
			defer processor.Close()
			return processor.Query()
		},
	})

	queryCmd.PersistentFlags().BoolP("local", "l", false, "query local database rather than project")
}
