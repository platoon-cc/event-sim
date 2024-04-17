package cmd

import (
	"github.com/platoon-cc/platoon-cli/internal/client"
	"github.com/platoon-cc/platoon-cli/internal/processor"
	"github.com/platoon-cc/platoon-cli/internal/settings"
	"github.com/spf13/cobra"
)

func init() {
	eventsCmd := &cobra.Command{
		Use:   "events",
		Short: "Interact with the stored events on the backend",
	}

	rootCmd.AddCommand(eventsCmd)
	eventsCmd.AddCommand(&cobra.Command{
		Use:   "ingest",
		Short: "Pull all the latest events down into a local database for querying",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectId, err := settings.GetActive("project")
			if err != nil {
				return err
			}
			processor, err := processor.New(projectId)
			if err != nil {
				return err
			}
			defer processor.Close()

			eventId, err := processor.GetPeakEventId()
			if err != nil {
				return err
			}

			platoon, err := client.New()
			if err != nil {
				return err
			}
			events, err := platoon.GetEvents(projectId, eventId)
			if err != nil {
				return err
			}

			return processor.StoreEvents(events, 0)
		},
	})
}
