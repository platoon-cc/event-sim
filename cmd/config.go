package cmd

import (
	"github.com/platoon-cc/platoon-cli/internal/client"
	"github.com/platoon-cc/platoon-cli/internal/settings"
	"github.com/platoon-cc/platoon-cli/internal/ui"
	"github.com/spf13/cobra"
)

func init() {
	configCmd := &cobra.Command{
		Use: "config",
	}

	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(&cobra.Command{
		Use: "all",
		RunE: func(cmd *cobra.Command, args []string) error {
			activeTeam, _ := settings.GetActive("team")

			platoon := client.New()
			teams, err := platoon.GetTeamList()
			if err != nil {
				return err
			}

			list := ui.NewList("Select Active Team")
			for _, t := range teams {
				list.AddItem(t.Id, t.Name, t.Id == activeTeam)
			}

			err = list.Run(
				func(i ui.ListItem) {
					settings.ClearCache("project")
					settings.ClearActive("project")
					settings.SetActive("team", i.Key)
				})
			if err != nil {
				return err
			}

			activeProject, _ := settings.GetActive("project")
			projects, err := platoon.GetProjectList()
			if err != nil {
				return err
			}
			list = ui.NewList("Select Active Project")
			for _, t := range projects {
				list.AddItem(t.Id, t.Name, t.Id == activeProject)
			}

			return list.Run(
				func(i ui.ListItem) {
					settings.SetActive("project", i.Key)
				})
		},
	})
}
