package cmd

import (
	"fmt"
	"time"

	"github.com/platoon-cc/platoon-cli/internal/client"
	"github.com/spf13/cobra"
)

func init() {
	eventsCmd := &cobra.Command{
		Use: "events",
	}

	rootCmd.AddCommand(eventsCmd)
	eventsCmd.AddCommand(&cobra.Command{
		Use: "get",
		RunE: func(cmd *cobra.Command, args []string) error {
			platoon := client.New()
			events, err := platoon.GetEvents()
			if err != nil {
				return err
			}
			for _, e := range events {
				t := time.UnixMilli(e.Timestamp).Format("2006/01/02 15:04")
				fmt.Printf("%d %s %s \t%s \t%v\n", e.Id, t, e.UserId, e.Event, e.Params)
			}
			return nil
		},
	})
}
