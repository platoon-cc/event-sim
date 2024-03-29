package cmd

import (
	"fmt"
	"time"

	"github.com/platoon-cc/platoon-cli/internal/client"
	"github.com/platoon-cc/platoon-cli/internal/sim"
	"github.com/spf13/cobra"
)

func init() {
	parentCmd := &cobra.Command{
		Use: "sim",
	}

	rootCmd.AddCommand(parentCmd)

	runCmd := &cobra.Command{
		Use: "run",

		RunE: func(cmd *cobra.Command, args []string) error {
			startTime := time.Now()
			for i := range 10 {
				events := sim.SimulateSessionForUser(i, startTime)
				res, err := sim.Serialise(events)
				if err != nil {
					return err
				}

				fmt.Println(res)

				if send, err := cmd.Flags().GetBool("send"); err != nil {
					return err
				} else if send {
					platoon := client.New()
					err := platoon.PostSimEvents(events)
					if err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
	runCmd.Flags().BoolP("send", "s", false, "whether to send the events to be ingested by the server")
	parentCmd.AddCommand(runCmd)
}
