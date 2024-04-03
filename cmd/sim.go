package cmd

import (
	"fmt"
	"time"

	"github.com/platoon-cc/platoon-cli/internal/client"
	"github.com/platoon-cc/platoon-cli/internal/model"
	"github.com/platoon-cc/platoon-cli/internal/processor"
	"github.com/platoon-cc/platoon-cli/internal/sim"
	"github.com/spf13/cobra"
)

type SimRun struct {
	events [][]model.Event
}

func init() {
	parentCmd := &cobra.Command{
		Use: "sim",
	}

	rootCmd.AddCommand(parentCmd)

	dumpCmd := &cobra.Command{
		Use:   "dump",
		Short: "writes sim results out to terminal",
		RunE: func(cmd *cobra.Command, args []string) error {
			run, err := runSim()
			if err != nil {
				return err
			}

			for _, r := range run.events {
				str, err := sim.Serialise(r)
				if err != nil {
					return err
				}
				fmt.Println(str)
			}

			return nil
		},
	}
	parentCmd.AddCommand(dumpCmd)

	localCmd := &cobra.Command{
		Use:   "local",
		Short: "writes sim results to local database for querying",
		RunE: func(cmd *cobra.Command, args []string) error {
			run, err := runSim()
			if err != nil {
				return err
			}

			processor, err := processor.New("local")
			if err != nil {
				return err
			}
			defer processor.Close()

			for _, r := range run.events {
				peakId, err := processor.GetPeakEventId()
				if err != nil {
					return err
				}
				if err := processor.StoreEvents(r, peakId); err != nil {
					return err
				}
			}

			return nil
		},
	}
	parentCmd.AddCommand(localCmd)

	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "sends sim results to currently configured project on platoon backend",
		RunE: func(cmd *cobra.Command, args []string) error {
			run, err := runSim()
			if err != nil {
				return err
			}

			platoon := client.New()
			for _, r := range run.events {
				err := platoon.PostSimEvents(r)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	parentCmd.AddCommand(sendCmd)
}

func runSim() (SimRun, error) {
	run := SimRun{}

	startTime := time.Now()
	for i := range 10 {
		events := sim.SimulateSessionForUser(i, startTime)
		run.events = append(run.events, events)
	}
	return run, nil
}
