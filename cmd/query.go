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
	return processor.Query(query)
}

func init() {
	queryCmd := &cobra.Command{
		Use:   "query",
		Short: "Freeform query",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := common(cmd, args[0]); err != nil {
				return err
			}

			return nil
		},
	}
	rootCmd.AddCommand(queryCmd)

	queryCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List out all event in order",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := common(cmd, `select * from event;`); err != nil {
				return err
			}

			return nil
		},
	})

	queryCmd.AddCommand(&cobra.Command{
		Use:   "dau",
		Short: "Daily Active Users",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := common(cmd, `select date(timestamp/1000, 'unixepoch') as date, count(distinct user_id) as user from event where event ='$sessionBegin' group by date;`); err != nil {
				return err
			}
			return nil
		},
	})

	queryCmd.AddCommand(&cobra.Command{
		Use:   "mau",
		Short: "Monthly Active Users",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := common(cmd, `select count(distinct user_id) as user from event where event = '$sessionBegin' and (unixepoch()-timestamp/1000) < (60*60*24*30);`); err != nil {
				return err
			}
			return nil
		},
	})

	queryCmd.AddCommand(&cobra.Command{
		Use:   "challengeScores",
		Short: "Clodhopper-specific challenges",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := common(cmd, `select payload ->> 'name' as Key, avg(payload ->> 'score') as Score from event where event='challengeEnd' group by Key;`); err != nil {
				return err
			}
			return nil
		},
	})

	queryCmd.PersistentFlags().BoolP("local", "l", false, "query local database rather than project")
}

// identifying the pairs of sessions/challenges etc
// select user_id, (timestamp-prev)as duration from (select *,lag(timestamp, 1) over() prev from event where event like 'challenge%') wher
// e event like '%End';

// hourly/daily/weekly active users
// select count(distinct user_id), (timestamp/3600000) as hour, datetime(timestamp/1000, 'unixepoch') from event where event ='$sessionBegin' group by hour;
