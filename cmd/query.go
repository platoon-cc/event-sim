package cmd

import (
	"fmt"
	"time"

	"github.com/platoon-cc/platoon-cli/internal/model"
	"github.com/platoon-cc/platoon-cli/internal/processor"
	"github.com/platoon-cc/platoon-cli/internal/settings"
	"github.com/spf13/cobra"
)

func common(cmd *cobra.Command, query string, res any) error {
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
	return processor.Query(query, res)
}

func init() {
	queryCmd := &cobra.Command{
		Use: "query",
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
			res := []model.Event{}

			if err := common(cmd, `select * from events;`, &res); err != nil {
				return err
			}

			for _, v := range res {
				t := time.UnixMilli(v.Timestamp).Format("2006/01/02 15:04")
				fmt.Printf("%d %s \t%s \t%s \t%+v\n", v.Id, t, v.Event, v.UserId, v.Payload)
			}

			return nil
		},
	})
	queryCmd.AddCommand(&cobra.Command{
		Use:   "challengeScores",
		Short: "Clodhopper-specific challenges",
		RunE: func(cmd *cobra.Command, args []string) error {
			res := []struct {
				Key   string
				Score float32
			}{}

			if err := common(cmd, `select payload ->> 'name' as Key, avg(payload ->> 'score') as Score from events where event='challengeEnd' group by Key;`, &res); err != nil {
				return err
			}

			for _, v := range res {
				fmt.Printf("%s - %f\n", v.Key, v.Score)
			}

			return nil
		},
	})

	queryCmd.PersistentFlags().BoolP("local", "l", false, "query local database rather than project")
}
