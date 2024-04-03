package cmd

import (
	"github.com/platoon-cc/platoon-cli/internal/client"
	"github.com/platoon-cc/platoon-cli/internal/processor"
	"github.com/platoon-cc/platoon-cli/internal/settings"
	"github.com/spf13/cobra"
)

func init() {
	eventsCmd := &cobra.Command{
		Use: "events",
	}

	rootCmd.AddCommand(eventsCmd)
	eventsCmd.AddCommand(&cobra.Command{
		Use: "ingest",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectId, err := settings.GetActive("project")
			if err != nil {
				return err
			}
			platoon := client.New()
			events, err := platoon.GetEvents(projectId)
			if err != nil {
				return err
			}

			processor, err := processor.New(projectId)
			if err != nil {
				return err
			}
			defer processor.Close()
			return processor.StoreEvents(events, 0)
		},
	})
}
