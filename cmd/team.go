package cmd

import (
	"github.com/platoon-cc/platoon-cli/internal/client"
	"github.com/platoon-cc/platoon-cli/internal/settings"
	"github.com/platoon-cc/platoon-cli/internal/ui"
	"github.com/spf13/cobra"
)

func init() {
	teamCmd := &cobra.Command{
		Use: "team",
	}

	rootCmd.AddCommand(teamCmd)
	teamCmd.AddCommand(&cobra.Command{
		Use: "choose",
		RunE: func(cmd *cobra.Command, args []string) error {
			key := "team"
			activeTeam := settings.GetActive("team")

			platoon := client.New(settings.GetAuthToken())
			teams, err := platoon.GetTeamList()
			if err != nil {
				return err
			}

			list := ui.NewList("Choose Active Team")
			for _, t := range teams {
				list.AddItem(t.Id, t.Name, t.Id == activeTeam)
			}

			return list.Run(
				func(i ui.ListItem) {
					settings.ClearCache("project")
					settings.ClearActive("project")
					settings.SetActive(key, i.Key)
				})
		},
	})
}
