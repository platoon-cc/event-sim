package cmd

import (
	"github.com/platoon-cc/platoon-cli/internal/processor"
	"github.com/platoon-cc/platoon-cli/internal/settings"
	"github.com/spf13/cobra"
)

func common(cmd *cobra.Command, query string) error {
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
	return processor.Query2(query)
}

func init() {
	queryCmd := &cobra.Command{
		Use:   "query",
		Short: "freeform query",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := common(cmd, args[0]); err != nil {
				return err
			}

			return nil
		},
	}

	rootCmd.AddCommand(queryCmd)
	// queryCmd.AddCommand(&cobra.Command{
	// 	Use: "dump",
	// 	RunE: func(cmd *cobra.Command, args []string) error {
	// 		return common(cmd, args, `select *payload ->> 'name' as key, avg(payload ->> 'score') from events where event='challengeEnd' group by key;`)
	// 	},
	// })

	queryCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List out all events in order",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := common(cmd, `select * from events;`); err != nil {
				return err
			}

			return nil
		},
	})

	queryCmd.AddCommand(&cobra.Command{
		Use:   "challengeScores",
		Short: "Clodhopper-specific challenges",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := common(cmd, `select payload ->> 'name' as Key, avg(payload ->> 'score') as Score from events where event='challengeEnd' group by Key;`); err != nil {
				return err
			}
			return nil
		},
	})

	queryCmd.PersistentFlags().BoolP("local", "l", false, "query local database rather than project")
}
